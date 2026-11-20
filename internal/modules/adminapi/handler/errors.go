package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	shareddb "github.com/fun-dotto/server/internal/shared/db"
	"github.com/fun-dotto/server/internal/shared/rpcerror"
	"github.com/jackc/pgerrcode"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"gorm.io/gorm"
)

const errorDomain rpcerror.Domain = "admin.dotto"

func invalidArgument(reason, field, description string) error {
	return errorDomain.New(
		connect.CodeInvalidArgument,
		reason,
		fmt.Sprintf("invalid %s: %s", field, description),
		nil,
		&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{Field: field, Description: description}}},
	)
}

// toConnectError は DB のエラーを Connect のエラーに変換する。
func toConnectError(ctx context.Context, r resourceName, keys []keyField, err error) error {
	metadata := make(map[string]string, len(keys))
	for _, k := range keys {
		metadata[r.metadataKey(k.name)] = k.String()
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errorDomain.New(connect.CodeNotFound, r.upper()+"_NOT_FOUND", r.human()+" not found", metadata)
	}
	switch shareddb.PostgresErrorCode(err) {
	case pgerrcode.UniqueViolation:
		return errorDomain.New(connect.CodeAlreadyExists, r.upper()+"_ALREADY_EXISTS",
			r.human()+" already exists; use a different identifier or update the existing one", metadata)
	case pgerrcode.ForeignKeyViolation:
		return errorDomain.New(connect.CodeFailedPrecondition, r.upper()+"_REFERENCE_VIOLATION",
			r.human()+" references a resource that does not exist, or is referenced by another resource", metadata)
	// CHECK 制約違反に加え、不正な UUID・日付など DB が値を解釈できなかった場合も含む
	case pgerrcode.CheckViolation, pgerrcode.InvalidTextRepresentation,
		pgerrcode.InvalidDatetimeFormat, pgerrcode.DatetimeFieldOverflow:
		return errorDomain.New(connect.CodeInvalidArgument, r.upper()+"_INVALID",
			r.human()+" has a value that violates a constraint; check the request fields", metadata,
			&errdetails.BadRequest{})
	}
	if cerr := errorDomain.ContextError(err); cerr != nil {
		return cerr
	}
	slog.ErrorContext(ctx, "admin-api internal error", "resource", string(r), "error", err)
	return errorDomain.New(connect.CodeInternal, "INTERNAL", "internal error", nil)
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

func keyOf(keys []keyField) map[string]any {
	m := make(map[string]any, len(keys))
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
