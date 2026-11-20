package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/fun-dotto/server/internal/shared/rpcerror"
	"gorm.io/gorm"
)

const errorDomain rpcerror.Domain = "app.dotto"

func invalidArgument(field, description string) error {
	var v rpcerror.Violations
	v.Add(field, description)
	return v.Err(errorDomain)
}

func unauthenticated() error {
	return errorDomain.New(connect.CodeUnauthenticated, "UNAUTHENTICATED", "authentication is required", nil)
}

// permissionDenied はログイン中のユーザーのものではないリソースを操作しようとした場合のエラー。
// 存在の有無を漏らさないよう、存在しない場合も同じエラーを返す（AIP-211）。
func permissionDenied(resource string, metadata map[string]string) error {
	return errorDomain.New(connect.CodePermissionDenied, strings.ToUpper(resource)+"_PERMISSION_DENIED",
		fmt.Sprintf("permission denied on %s; it does not belong to the authenticated user", strings.ReplaceAll(resource, "_", " ")), metadata)
}

func notFound(resource string, metadata map[string]string) error {
	return errorDomain.New(connect.CodeNotFound, strings.ToUpper(resource)+"_NOT_FOUND",
		strings.ReplaceAll(resource, "_", " ")+" not found", metadata)
}

func userNotFound() error {
	return errorDomain.New(connect.CodeFailedPrecondition, "USER_NOT_FOUND",
		"user not found; create the user with UpdateUser first", nil)
}

// unavailable は依存先（microCMS など）に到達できない、または未設定の場合のエラー。リトライ可能。
// 依存先の HTTP クライアントのタイムアウトも context.DeadlineExceeded を満たすため、
// err ではなくリクエストの ctx で、リクエスト側の期限切れ・キャンセルと区別する。
func unavailable(ctx context.Context, operation string, err error) error {
	if ctx.Err() != nil {
		return internalError(ctx, operation, ctx.Err())
	}
	slog.WarnContext(ctx, "app-api dependency unavailable", "operation", operation, "error", err)
	return errorDomain.New(connect.CodeUnavailable, "UNAVAILABLE", "service temporarily unavailable", nil)
}

// internalError は DB などのエラーを Connect のエラーに変換する。
func internalError(ctx context.Context, operation string, err error) error {
	if cerr := errorDomain.ContextError(err); cerr != nil {
		return cerr
	}
	slog.ErrorContext(ctx, "app-api internal error", "operation", operation, "error", err)
	return errorDomain.New(connect.CodeInternal, "INTERNAL", "internal error", nil)
}

func isNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }
