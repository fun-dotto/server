package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/fun-dotto/server/internal/modules/adminapi/repository"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

const errorDomain = "admin.dotto"

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

// toConnectError は repository のエラーを Connect のエラーに変換する。
func toConnectError(ctx context.Context, r resourceName, keys []keyField, err error) error {
	metadata := make(map[string]string, len(keys))
	for _, k := range keys {
		metadata[r.metadataKey(k.name)] = k.String()
	}
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return newError(connect.CodeNotFound, r.upper()+"_NOT_FOUND", r.human()+" not found", metadata)
	case errors.Is(err, repository.ErrAlreadyExists):
		return newError(connect.CodeAlreadyExists, r.upper()+"_ALREADY_EXISTS",
			r.human()+" already exists; use a different identifier or update the existing one", metadata)
	case errors.Is(err, repository.ErrForeignKeyViolation):
		return newError(connect.CodeFailedPrecondition, r.upper()+"_REFERENCE_VIOLATION",
			r.human()+" references a resource that does not exist, or is referenced by another resource", metadata)
	case errors.Is(err, repository.ErrCheckViolation):
		return newError(connect.CodeInvalidArgument, r.upper()+"_INVALID",
			r.human()+" has a value that violates a constraint; check the request fields", metadata,
			&errdetails.BadRequest{})
	case errors.Is(err, context.DeadlineExceeded):
		return newError(connect.CodeDeadlineExceeded, "DEADLINE_EXCEEDED", "request deadline exceeded", nil)
	case errors.Is(err, context.Canceled):
		return newError(connect.CodeCanceled, "CANCELED", "request canceled", nil)
	default:
		slog.ErrorContext(ctx, "admin-api internal error", "resource", string(r), "error", err)
		return newError(connect.CodeInternal, "INTERNAL", "internal error", nil)
	}
}

// keyField はリソースを特定する識別子のフィールド（proto のフィールド名 = DB の列名）。
type keyField struct {
	name  string
	value any
}

func (k keyField) String() string {
	if s, ok := k.value.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprint(k.value)
}

func keyOf(keys []keyField) repository.Key {
	m := make(repository.Key, len(keys))
	for _, k := range keys {
		m[k.name] = k.value
	}
	return m
}

// resourceName は lower_snake_case のリソース名（例: calendar_date）。
type resourceName string

func (r resourceName) upper() string { return strings.ToUpper(string(r)) }

func (r resourceName) human() string { return strings.ReplaceAll(string(r), "_", " ") }

// metadataKey は ErrorInfo.metadata のキーを lowerCamelCase で返す。
// `id` はリソース名を前置する（例: subjectID）。
func (r resourceName) metadataKey(field string) string {
	if field == "id" {
		field = string(r) + "_id"
	}
	parts := strings.Split(field, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
