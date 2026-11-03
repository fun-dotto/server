package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
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

func invalidArgument(reason, field, description string) error {
	return newError(
		connect.CodeInvalidArgument,
		reason,
		fmt.Sprintf("invalid %s: %s", field, description),
		nil,
		&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{Field: field, Description: description}}},
	)
}

func unauthenticated() error {
	return newError(connect.CodeUnauthenticated, "UNAUTHENTICATED", "authentication is required", nil)
}

// toConnectError は service 層のエラーを Connect のエラーに変換する。
func toConnectError(ctx context.Context, operation string, err error) error {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return newError(connect.CodeNotFound, "USER_NOT_FOUND", "user not found; create it with UpdateUser", nil)
	case errors.Is(err, context.DeadlineExceeded):
		return newError(connect.CodeDeadlineExceeded, "DEADLINE_EXCEEDED", "request deadline exceeded", nil)
	case errors.Is(err, context.Canceled):
		return newError(connect.CodeCanceled, "CANCELED", "request canceled", nil)
	default:
		slog.ErrorContext(ctx, "app-api internal error", "operation", operation, "error", err)
		return newError(connect.CodeInternal, "INTERNAL", "internal error", nil)
	}
}

func errorsIsUserNotFound(err error) bool { return errors.Is(err, domain.ErrUserNotFound) }
