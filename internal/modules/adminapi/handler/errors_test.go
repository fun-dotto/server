package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

func TestToConnectError(t *testing.T) {
	keys := []keyField{{name: "id", value: "s1"}}
	metadata := map[string]string{"subjectId": "s1"}
	errorInfo := func(reason string, metadata map[string]string) *errdetails.ErrorInfo {
		return &errdetails.ErrorInfo{Reason: reason, Domain: "admin.dotto", Metadata: metadata}
	}
	pgError := func(code string) error {
		return fmt.Errorf("create subject: %w", &pgconn.PgError{Code: code})
	}
	invalid := []proto.Message{errorInfo("SUBJECT_INVALID", metadata), &errdetails.BadRequest{}}

	tests := []struct {
		name string
		err  error
		code connect.Code
		want []proto.Message
	}{
		{name: "レコードがない", err: gorm.ErrRecordNotFound, code: connect.CodeNotFound, want: []proto.Message{errorInfo("SUBJECT_NOT_FOUND", metadata)}},
		{name: "一意制約違反", err: pgError(pgerrcode.UniqueViolation), code: connect.CodeAlreadyExists, want: []proto.Message{errorInfo("SUBJECT_ALREADY_EXISTS", metadata)}},
		{name: "外部キー制約違反", err: pgError(pgerrcode.ForeignKeyViolation), code: connect.CodeFailedPrecondition, want: []proto.Message{errorInfo("SUBJECT_REFERENCE_VIOLATION", metadata)}},
		{name: "CHECK 制約違反", err: pgError(pgerrcode.CheckViolation), code: connect.CodeInvalidArgument, want: invalid},
		{name: "不正な UUID", err: pgError(pgerrcode.InvalidTextRepresentation), code: connect.CodeInvalidArgument, want: invalid},
		{name: "日時の形式が不正", err: pgError(pgerrcode.InvalidDatetimeFormat), code: connect.CodeInvalidArgument, want: invalid},
		{name: "日時のフィールドが範囲外", err: pgError(pgerrcode.DatetimeFieldOverflow), code: connect.CodeInvalidArgument, want: invalid},
		{name: "タイムアウト", err: fmt.Errorf("query: %w", context.DeadlineExceeded), code: connect.CodeDeadlineExceeded, want: []proto.Message{errorInfo("DEADLINE_EXCEEDED", nil)}},
		{name: "それ以外は INTERNAL", err: errors.New("boom"), code: connect.CodeInternal, want: []proto.Message{errorInfo("INTERNAL", nil)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cerr, ok := errors.AsType[*connect.Error](toConnectError(context.Background(), "subject", keys, tt.err))
			if !ok {
				t.Fatalf("toConnectError() did not return *connect.Error")
			}
			if got := cerr.Code(); got != tt.code {
				t.Errorf("code = %v, want %v", got, tt.code)
			}
			details := cerr.Details()
			if len(details) != len(tt.want) {
				t.Fatalf("len(details) = %d, want %d", len(details), len(tt.want))
			}
			for i, d := range details {
				got, err := d.Value()
				if err != nil {
					t.Fatalf("details[%d].Value() error = %v", i, err)
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("details[%d] = %v, want %v", i, got, tt.want[i])
				}
			}
		})
	}
}
