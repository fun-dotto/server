package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
	"gorm.io/gorm"
)

const errorDomain = "app.dotto"

func newError(code connect.Code, reason, message string, metadata map[string]string, details ...protoadapt.MessageV1) *connect.Error {
	err := connect.NewError(code, errors.New(message))
	if d, derr := connect.NewErrorDetail(&errdetails.ErrorInfo{
		Reason:   reason,
		Domain:   errorDomain,
		Metadata: metadata,
	}); derr == nil {
		err.AddDetail(d)
	}
	for _, detail := range details {
		if d, derr := connect.NewErrorDetail(protoadapt.MessageV2Of(detail)); derr == nil {
			err.AddDetail(d)
		}
	}
	return err
}

// violations は INVALID_ARGUMENT の BadRequest.field_violations を集める。
type violations []*errdetails.BadRequest_FieldViolation

func (v *violations) add(field, description string) {
	*v = append(*v, &errdetails.BadRequest_FieldViolation{Field: field, Description: description})
}

func (v violations) err() error {
	if len(v) == 0 {
		return nil
	}
	fields := make([]string, 0, len(v))
	for _, fv := range v {
		fields = append(fields, fv.GetField())
	}
	return newError(
		connect.CodeInvalidArgument,
		"INVALID_ARGUMENT",
		"request has invalid fields: "+strings.Join(fields, ", "),
		nil,
		&errdetails.BadRequest{FieldViolations: v},
	)
}

func invalidArgument(field, description string) error {
	var v violations
	v.add(field, description)
	return v.err()
}

func unauthenticated() error {
	return newError(connect.CodeUnauthenticated, "UNAUTHENTICATED", "authentication is required", nil)
}

// permissionDenied はログイン中のユーザーのものではないリソースを操作しようとした場合のエラー。
// 存在の有無を漏らさないよう、存在しない場合も同じエラーを返す（AIP-211）。
func permissionDenied(resource string, metadata map[string]string) error {
	return newError(connect.CodePermissionDenied, strings.ToUpper(resource)+"_PERMISSION_DENIED",
		fmt.Sprintf("permission denied on %s; it does not belong to the authenticated user", strings.ReplaceAll(resource, "_", " ")), metadata)
}

func notFound(resource string, metadata map[string]string) error {
	return newError(connect.CodeNotFound, strings.ToUpper(resource)+"_NOT_FOUND",
		strings.ReplaceAll(resource, "_", " ")+" not found", metadata)
}

func userNotFound() error {
	return newError(connect.CodeFailedPrecondition, "USER_NOT_FOUND",
		"user not found; create the user with UpdateUser first", nil)
}

// internalError は DB などのエラーを Connect のエラーに変換する。
func internalError(ctx context.Context, operation string, err error) error {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return newError(connect.CodeDeadlineExceeded, "DEADLINE_EXCEEDED", "request deadline exceeded", nil)
	case errors.Is(err, context.Canceled):
		return newError(connect.CodeCanceled, "CANCELED", "request canceled", nil)
	default:
		slog.ErrorContext(ctx, "app-api internal error", "operation", operation, "error", err)
		return newError(connect.CodeInternal, "INTERNAL", "internal error", nil)
	}
}

func isNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }

func postgresErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
)
