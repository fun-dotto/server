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
// テスト終了時に元に戻す。slog.SetDefault は log パッケージの出力先とフラグも書き換え、
// 既定のロガーに戻しても元に戻らないため、log 側も退避して戻す。
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
	wantRes := connect.NewResponse(&emptypb.Empty{})
	wantErr := errors.New("boom")

	var gotCtx context.Context
	next := func(ctx context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
		gotCtx = ctx
		return wantRes, wantErr
	}

	res, err := Timeout(time.Minute)(next)(context.Background(), connect.NewRequest(&emptypb.Empty{}))
	if res != wantRes || err != wantErr {
		t.Errorf("Timeout() = (%v, %v), want (%v, %v)", res, err, wantRes, wantErr)
	}
	if _, ok := gotCtx.Deadline(); !ok {
		t.Error("下流の context に deadline が設定されていない")
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
