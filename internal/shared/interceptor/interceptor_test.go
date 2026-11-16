package interceptor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/fun-dotto/server/internal/shared/logging"
	"google.golang.org/protobuf/types/known/emptypb"
)

// captureLog は slog の既定ロガーの出力を JSON 形式でバッファに書き出すよう差し替え、
// テスト終了時に元に戻す。
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	prevLogger := slog.Default()
	prevWriter, prevFlags := log.Writer(), log.Flags()
	t.Cleanup(func() {
		slog.SetDefault(prevLogger)
		log.SetOutput(prevWriter)
		log.SetFlags(prevFlags)
	})
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))
	return &buf
}

func TestTimeout(t *testing.T) {
	const d = time.Minute
	req := connect.NewRequest(&emptypb.Empty{})
	wantRes := connect.NewResponse(&emptypb.Empty{})
	wantErr := errors.New("boom")

	var gotCtx context.Context
	next := func(ctx context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
		gotCtx = ctx
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatal("下流の context に deadline が設定されていない")
		}
		if limit := time.Now().Add(d); deadline.After(limit) {
			t.Errorf("deadline = %v, want <= %v", deadline, limit)
		}
		if err := ctx.Err(); err != nil {
			t.Errorf("下流の処理中に context が終了している: %v", err)
		}
		return wantRes, wantErr
	}

	res, err := Timeout(d)(next)(context.Background(), req)
	if res != wantRes {
		t.Errorf("response = %v, want %v", res, wantRes)
	}
	if err != wantErr {
		t.Errorf("error = %v, want %v", err, wantErr)
	}
	if gotCtx == nil {
		t.Fatal("下流の処理が呼ばれていない")
	}
	if err := gotCtx.Err(); !errors.Is(err, context.Canceled) {
		t.Errorf("return 後の context の Err() = %v, want %v", err, context.Canceled)
	}
}

func TestTimeoutExceeded(t *testing.T) {
	next := func(ctx context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}

	_, err := Timeout(time.Millisecond)(next)(context.Background(), connect.NewRequest(&emptypb.Empty{}))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("error = %v, want %v", err, context.DeadlineExceeded)
	}
}

func TestRecordError(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "エラーがなければそのまま返す",
			err:  nil,
		},
		{
			name: "エラーがあればそのまま返す",
			err:  connect.NewError(connect.CodeInternal, errors.New("boom")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := captureLog(t)
			wantRes := connect.NewResponse(&emptypb.Empty{})
			next := func(context.Context, connect.AnyRequest) (connect.AnyResponse, error) {
				return wantRes, tt.err
			}

			// logging.Middleware を通っていない context でも panic せずに動くこと。
			res, err := RecordError()(next)(context.Background(), connect.NewRequest(&emptypb.Empty{}))
			if res != wantRes {
				t.Errorf("response = %v, want %v", res, wantRes)
			}
			if err != tt.err {
				t.Errorf("error = %v, want %v", err, tt.err)
			}

			// Middleware を通っていない場合、logging.RecordError はエラーを直接ログに出す。
			got := buf.String()
			if tt.err == nil && got != "" {
				t.Errorf("エラーがないのにログが出力された: %s", got)
			}
			if tt.err != nil && !strings.Contains(got, tt.err.Error()) {
				t.Errorf("ログ %q にエラー %q が含まれていない", got, tt.err.Error())
			}
		})
	}
}

func TestRecordErrorWithMiddleware(t *testing.T) {
	const procedure = "/test.v1.TestService/Call"
	wantErr := connect.NewError(connect.CodeInternal, errors.New("boom"))
	handler := connect.NewUnaryHandler(
		procedure,
		func(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
			return nil, wantErr
		},
		connect.WithInterceptors(RecordError()),
	)

	buf := captureLog(t)
	req := httptest.NewRequest(http.MethodPost, procedure, strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	logging.Middleware(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}

	var entry struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("アクセスログを JSON として読めない: %v (log: %s)", err, buf.String())
	}
	if entry.Error != wantErr.Error() {
		t.Errorf("アクセスログの error = %q, want %q", entry.Error, wantErr.Error())
	}
}
