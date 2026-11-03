package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
)

type fakeIDToken struct{}

func (fakeIDToken) VerifyIDToken(_ context.Context, token string) (User, error) {
	if token != "valid-id" {
		return User{}, errors.New("invalid")
	}
	return User{ID: "user1", Email: "user1@example.com"}, nil
}

var appCheck = AppCheckTokenVerifierFunc(func(token string) error {
	if token != "valid-app" {
		return errors.New("invalid")
	}
	return nil
})

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

type announcementHandler struct {
	appv1connect.UnimplementedAnnouncementServiceHandler
}

func (announcementHandler) ListAnnouncements(ctx context.Context, _ *connect.Request[appv1.ListAnnouncementsRequest]) (*connect.Response[appv1.ListAnnouncementsResponse], error) {
	return connect.NewResponse(&appv1.ListAnnouncementsResponse{}), nil
}

func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	opts := connect.WithInterceptors(NewInterceptor(appCheck, fakeCFAccess{}, fakeService{}, fakeIDToken{}))
	mux := http.NewServeMux()
	mux.Handle(appv1connect.NewUserServiceHandler(userHandler{}, opts))
	mux.Handle(appv1connect.NewAnnouncementServiceHandler(announcementHandler{}, opts))
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestInterceptor(t *testing.T) {
	srv := newServer(t)
	users := appv1connect.NewUserServiceClient(http.DefaultClient, srv.URL)
	announcements := appv1connect.NewAnnouncementServiceClient(http.DefaultClient, srv.URL)

	tests := []struct {
		name     string
		call     func(req func(http.Header)) error
		appCheck string
		cfAccess string
		service  string
		idToken  string
		wantCode connect.Code
	}{
		{"ユーザー不要の RPC は App Check のみで通る", listAnnouncements(announcements), "Bearer valid-app", "", "", "", 0},
		{"App Check がない場合は UNAUTHENTICATED", listAnnouncements(announcements), "", "", "", "", connect.CodeUnauthenticated},
		{"App Check が不正な場合は UNAUTHENTICATED", listAnnouncements(announcements), "invalid", "", "", "", connect.CodeUnauthenticated},
		{"ユーザーが必要な RPC は ID トークンで通る", getUser(users), "valid-app", "", "", "Bearer valid-id", 0},
		{"ID トークンがない場合は UNAUTHENTICATED", getUser(users), "valid-app", "", "", "", connect.CodeUnauthenticated},
		{"ID トークンが不正な場合は UNAUTHENTICATED", getUser(users), "valid-app", "", "", "Bearer invalid", connect.CodeUnauthenticated},
		{"Cloudflare Access の JWT があれば App Check なしで通る", listAnnouncements(announcements), "", "valid-cf", "", "", 0},
		{"Cloudflare Access の JWT が不正で App Check がない場合は UNAUTHENTICATED", listAnnouncements(announcements), "", "invalid", "", "", connect.CodeUnauthenticated},
		{"Cloudflare Access の JWT が不正でも App Check があれば通る", listAnnouncements(announcements), "valid-app", "invalid", "", "", 0},
		{"Cloudflare Access の JWT があってもユーザーが必要な RPC は ID トークンが必要", getUser(users), "", "valid-cf", "", "", connect.CodeUnauthenticated},
		{"Cloudflare Access の JWT と ID トークンで通る", getUser(users), "", "valid-cf", "", "Bearer valid-id", 0},
		{"許可されたサービスの ID トークンがあれば App Check なしで通る", listAnnouncements(announcements), "", "", "Bearer valid-svc", "", 0},
		{"サービスの ID トークンが不正で App Check がない場合は UNAUTHENTICATED", listAnnouncements(announcements), "", "", "Bearer invalid", "", connect.CodeUnauthenticated},
		{"サービスの ID トークンに Bearer がない場合は App Check が必要", listAnnouncements(announcements), "", "", "valid-svc", "", connect.CodeUnauthenticated},
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
					h.Set(cfAccessJWTHeader, tt.cfAccess)
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

func listAnnouncements(c appv1connect.AnnouncementServiceClient) func(func(http.Header)) error {
	return func(set func(http.Header)) error {
		req := connect.NewRequest(&appv1.ListAnnouncementsRequest{})
		set(req.Header())
		_, err := c.ListAnnouncements(context.Background(), req)
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
