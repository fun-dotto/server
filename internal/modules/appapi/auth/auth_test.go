package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"connectrpc.com/connect"
	"firebase.google.com/go/v4/appcheck"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/shared/cfaccess"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type fakeIDToken struct{}

func (fakeIDToken) VerifyIDToken(_ context.Context, token string) (User, error) {
	switch token {
	case "valid-id":
		return User{ID: "user1", Email: "user1@example.com"}, nil
	case "revoked-id":
		return User{}, ErrIDTokenRevoked
	case "disabled-id":
		return User{}, ErrUserDisabled
	}
	return User{}, errors.New("invalid")
}

const allowedAppID = "1:123:ios:abc"

// fakeAppCheck は App Check の backend を模す。"valid-app" は毎回新しいトークンとして扱い、
// "once-" で始まるトークンは 2 回目以降を使用済みとして扱う。
type fakeAppCheck struct {
	mu       sync.Mutex
	consumed map[string]bool
}

func (f *fakeAppCheck) VerifyOneTimeToken(_ context.Context, token string) (*appcheck.DecodedAppCheckToken, error) {
	appID := allowedAppID
	switch {
	case token == "valid-app":
	case token == "other-app":
		appID = "1:123:web:def"
	case len(token) > 5 && token[:5] == "once-":
	default:
		return nil, errors.New("invalid")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	consumed := f.consumed[token]
	if token != "valid-app" && token != "other-app" {
		f.consumed[token] = true
	}
	return &appcheck.DecodedAppCheckToken{AppID: appID, AlreadyConsumed: &consumed}, nil
}

func newAppCheck(t *testing.T) *AppCheck {
	t.Helper()
	a, err := NewAppCheck(&fakeAppCheck{consumed: map[string]bool{}}, " "+allowedAppID+" , 1:123:android:xyz")
	if err != nil {
		t.Fatal(err)
	}
	return a
}

type fakeCFAccess struct{}

func (fakeCFAccess) Verify(_ context.Context, token string) error {
	if token != "valid-cf" {
		return errors.New("invalid")
	}
	return nil
}

type fakeService struct{}

func (fakeService) Verify(_ context.Context, token string) error {
	if token != "valid-svc" {
		return errors.New("invalid")
	}
	return nil
}

type userHandler struct {
	appv1connect.UnimplementedUserServiceHandler
}

func (userHandler) GetUser(ctx context.Context, _ *connect.Request[appv1.GetUserRequest]) (*connect.Response[appv1.GetUserResponse], error) {
	if u, ok := UserFromContext(ctx); !ok || u.ID != "user1" {
		return nil, connect.NewError(connect.CodeInternal, errors.New("user not in context"))
	}
	return connect.NewResponse(&appv1.GetUserResponse{User: &appv1.User{}}), nil
}

type articleHandler struct {
	appv1connect.UnimplementedArticleServiceHandler
}

func (articleHandler) ListArticles(ctx context.Context, _ *connect.Request[appv1.ListArticlesRequest]) (*connect.Response[appv1.ListArticlesResponse], error) {
	return connect.NewResponse(&appv1.ListArticlesResponse{}), nil
}

func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	opts := connect.WithInterceptors(NewInterceptor(newAppCheck(t), fakeCFAccess{}, fakeService{}, fakeIDToken{}))
	mux := http.NewServeMux()
	mux.Handle(appv1connect.NewUserServiceHandler(userHandler{}, opts))
	mux.Handle(appv1connect.NewArticleServiceHandler(articleHandler{}, opts))
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestInterceptor(t *testing.T) {
	srv := newServer(t)
	users := appv1connect.NewUserServiceClient(http.DefaultClient, srv.URL)
	articles := appv1connect.NewArticleServiceClient(http.DefaultClient, srv.URL)

	// "once-used" を一度使用済みにしておく。
	if err := listArticles(articles)(func(h http.Header) { h.Set(appCheckHeader, "once-used") }); err != nil {
		t.Fatalf("初回の App Check トークンは通るはず: %v", err)
	}

	tests := []struct {
		name     string
		call     func(req func(http.Header)) error
		appCheck string
		cfAccess string
		service  string
		idToken  string
		wantCode connect.Code
	}{
		{"ユーザー不要の RPC は App Check のみで通る", listArticles(articles), "Bearer valid-app", "", "", "", 0},
		{"App Check がない場合は UNAUTHENTICATED", listArticles(articles), "", "", "", "", connect.CodeUnauthenticated},
		{"App Check が不正な場合は UNAUTHENTICATED", listArticles(articles), "invalid", "", "", "", connect.CodeUnauthenticated},
		{"使用済みの App Check トークンは UNAUTHENTICATED", listArticles(articles), "once-used", "", "", "", connect.CodeUnauthenticated},
		{"許可されていないアプリの App Check トークンは UNAUTHENTICATED", listArticles(articles), "other-app", "", "", "", connect.CodeUnauthenticated},
		{"ユーザーが必要な RPC は ID トークンで通る", getUser(users), "valid-app", "", "", "Bearer valid-id", 0},
		{"ID トークンがない場合は UNAUTHENTICATED", getUser(users), "valid-app", "", "", "", connect.CodeUnauthenticated},
		{"ID トークンが不正な場合は UNAUTHENTICATED", getUser(users), "valid-app", "", "", "Bearer invalid", connect.CodeUnauthenticated},
		{"失効した ID トークンは UNAUTHENTICATED", getUser(users), "valid-app", "", "", "Bearer revoked-id", connect.CodeUnauthenticated},
		{"無効化されたユーザーの ID トークンは UNAUTHENTICATED", getUser(users), "valid-app", "", "", "Bearer disabled-id", connect.CodeUnauthenticated},
		{"Cloudflare Access の JWT があれば App Check なしで通る", listArticles(articles), "", "valid-cf", "", "", 0},
		{"Cloudflare Access の JWT が不正で App Check がない場合は UNAUTHENTICATED", listArticles(articles), "", "invalid", "", "", connect.CodeUnauthenticated},
		{"Cloudflare Access の JWT が不正でも App Check があれば通る", listArticles(articles), "valid-app", "invalid", "", "", 0},
		{"Cloudflare Access の JWT があってもユーザーが必要な RPC は ID トークンが必要", getUser(users), "", "valid-cf", "", "", connect.CodeUnauthenticated},
		{"Cloudflare Access の JWT と ID トークンで通る", getUser(users), "", "valid-cf", "", "Bearer valid-id", 0},
		{"許可されたサービスの ID トークンがあれば App Check なしで通る", listArticles(articles), "", "", "Bearer valid-svc", "", 0},
		{"サービスの ID トークンが不正で App Check がない場合は UNAUTHENTICATED", listArticles(articles), "", "", "Bearer invalid", "", connect.CodeUnauthenticated},
		{"サービスの ID トークンに Bearer がない場合は App Check が必要", listArticles(articles), "", "", "valid-svc", "", connect.CodeUnauthenticated},
		{"サービスの ID トークンがあってもユーザーが必要な RPC は ID トークンが必要", getUser(users), "", "", "Bearer valid-svc", "", connect.CodeUnauthenticated},
		{"サービスの ID トークンとユーザーの ID トークンで通る", getUser(users), "", "", "Bearer valid-svc", "Bearer valid-id", 0},
		{"Bearer がない場合は UNAUTHENTICATED", getUser(users), "valid-app", "", "", "valid-id", connect.CodeUnauthenticated},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call(func(h http.Header) {
				if tt.appCheck != "" {
					h.Set(appCheckHeader, tt.appCheck)
				}
				if tt.cfAccess != "" {
					h.Set(cfaccess.JWTHeader, tt.cfAccess)
				}
				if tt.service != "" {
					h.Set(serviceAuthHeader, tt.service)
				}
				if tt.idToken != "" {
					h.Set(authorizationHeader, tt.idToken)
				}
			})
			if tt.wantCode == 0 {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if got := connect.CodeOf(err); got != tt.wantCode {
				t.Fatalf("code = %v, want %v (err: %v)", got, tt.wantCode, err)
			}
		})
	}
}

func listArticles(c appv1connect.ArticleServiceClient) func(func(http.Header)) error {
	return func(set func(http.Header)) error {
		req := connect.NewRequest(&appv1.ListArticlesRequest{})
		set(req.Header())
		_, err := c.ListArticles(context.Background(), req)
		return err
	}
}

func getUser(c appv1connect.UserServiceClient) func(func(http.Header)) error {
	return func(set func(http.Header)) error {
		req := connect.NewRequest(&appv1.GetUserRequest{})
		set(req.Header())
		_, err := c.GetUser(context.Background(), req)
		return err
	}
}

func TestInterceptorWithoutCFAccess(t *testing.T) {
	opts := connect.WithInterceptors(NewInterceptor(newAppCheck(t), nil, fakeService{}, fakeIDToken{}))
	mux := http.NewServeMux()
	mux.Handle(appv1connect.NewArticleServiceHandler(articleHandler{}, opts))
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	articles := appv1connect.NewArticleServiceClient(http.DefaultClient, srv.URL)

	call := listArticles(articles)
	if err := call(func(h http.Header) { h.Set(cfaccess.JWTHeader, "valid-cf") }); connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Fatalf("Cloudflare Access を使わない場合は JWT があっても App Check が必要: err = %v", err)
	}
	if err := call(func(h http.Header) { h.Set(appCheckHeader, "valid-app") }); err != nil {
		t.Fatalf("App Check で通るはず: %v", err)
	}
	if err := call(func(h http.Header) { h.Set(serviceAuthHeader, "Bearer valid-svc") }); err != nil {
		t.Fatalf("サービスの ID トークンで通るはず: %v", err)
	}
}

func TestAppCheckReplay(t *testing.T) {
	srv := newServer(t)
	articles := appv1connect.NewArticleServiceClient(http.DefaultClient, srv.URL)
	call := listArticles(articles)
	set := func(h http.Header) { h.Set(appCheckHeader, "once-token") }

	if err := call(set); err != nil {
		t.Fatalf("初回は通るはず: %v", err)
	}
	err := call(set)
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Fatalf("2 回目は UNAUTHENTICATED のはず: err = %v", err)
	}
	if got := reasonOf(t, err); got != "APP_CHECK_TOKEN_ALREADY_CONSUMED" {
		t.Fatalf("reason = %q", got)
	}
}

func TestErrorReasons(t *testing.T) {
	srv := newServer(t)
	users := appv1connect.NewUserServiceClient(http.DefaultClient, srv.URL)
	articles := appv1connect.NewArticleServiceClient(http.DefaultClient, srv.URL)

	tests := []struct {
		name string
		call func(func(http.Header)) error
		set  func(http.Header)
		want string
	}{
		{"許可されていないアプリ", listArticles(articles), func(h http.Header) { h.Set(appCheckHeader, "other-app") }, "APP_CHECK_APP_NOT_ALLOWED"},
		{"失効した ID トークン", getUser(users), func(h http.Header) {
			h.Set(appCheckHeader, "valid-app")
			h.Set(authorizationHeader, "Bearer revoked-id")
		}, "ID_TOKEN_REVOKED"},
		{"無効化されたユーザー", getUser(users), func(h http.Header) {
			h.Set(appCheckHeader, "valid-app")
			h.Set(authorizationHeader, "Bearer disabled-id")
		}, "USER_DISABLED"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reasonOf(t, tt.call(tt.set)); got != tt.want {
				t.Fatalf("reason = %q, want %q", got, tt.want)
			}
		})
	}
}

func reasonOf(t *testing.T, err error) string {
	t.Helper()
	var cerr *connect.Error
	if !errors.As(err, &cerr) {
		t.Fatalf("not a connect error: %v", err)
	}
	for _, d := range cerr.Details() {
		v, derr := d.Value()
		if derr != nil {
			continue
		}
		if info, ok := v.(*errdetails.ErrorInfo); ok {
			return info.GetReason()
		}
	}
	return ""
}

func TestNewAppCheck(t *testing.T) {
	if _, err := NewAppCheck(&fakeAppCheck{}, " , "); err == nil {
		t.Fatal("許可するアプリ ID が空の場合はエラーのはず")
	}
	if _, err := NewAppCheck(nil, allowedAppID); err == nil {
		t.Fatal("Verifier が nil の場合はエラーのはず")
	}
	a, err := NewAppCheck(&fakeAppCheck{}, " a , ,b ")
	if err != nil {
		t.Fatal(err)
	}
	if len(a.allowedAppIDs) != 2 || a.allowedAppIDs[0] != "a" || a.allowedAppIDs[1] != "b" {
		t.Fatalf("allowedAppIDs = %v", a.allowedAppIDs)
	}
}
