// Package interceptor は Connect の Unary Interceptor のうち、各 API で共通して使うものを提供する。
package interceptor

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/fun-dotto/server/internal/shared/logging"
)

// Timeout は 1 リクエストあたりの処理時間の上限を設ける。
func Timeout(d time.Duration) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			ctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			return next(ctx, req)
		}
	}
}

// RecordError は RPC のエラーをアクセスログに記録する。
// アクセスログに含めるには logging.Middleware を通ったリクエストである必要がある。
func RecordError() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			res, err := next(ctx, req)
			if err != nil {
				logging.RecordError(ctx, err)
			}
			return res, err
		}
	}
}
