package rpcerror_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"

	"github.com/fun-dotto/server/internal/shared/rpcerror"
)

const testDomain rpcerror.Domain = "test.dotto"

// assertConnectError は err が code・message の Connect のエラーで、details が want と一致することを検証する。
func assertConnectError(t *testing.T, err error, code connect.Code, message string, want ...proto.Message) {
	t.Helper()
	cerr, ok := errors.AsType[*connect.Error](err)
	if !ok {
		t.Fatalf("err = %v (%T), want *connect.Error", err, err)
	}
	if got := cerr.Code(); got != code {
		t.Errorf("code = %v, want %v", got, code)
	}
	if got := cerr.Message(); got != message {
		t.Errorf("message = %q, want %q", got, message)
	}
	details := cerr.Details()
	if len(details) != len(want) {
		t.Fatalf("len(details) = %d, want %d", len(details), len(want))
	}
	for i, d := range details {
		got, derr := d.Value()
		if derr != nil {
			t.Fatalf("details[%d].Value() error = %v", i, derr)
		}
		if !proto.Equal(got, want[i]) {
			t.Errorf("details[%d] = %v, want %v", i, got, want[i])
		}
	}
}

func TestDomainNew(t *testing.T) {
	tests := []struct {
		name     string
		code     connect.Code
		reason   string
		message  string
		metadata map[string]string
		details  []protoadapt.MessageV1
		want     []proto.Message
	}{
		{
			name:    "metadata と details なし",
			code:    connect.CodeInternal,
			reason:  "INTERNAL",
			message: "internal error",
			want: []proto.Message{
				&errdetails.ErrorInfo{Reason: "INTERNAL", Domain: "test.dotto"},
			},
		},
		{
			name:     "metadata あり",
			code:     connect.CodeNotFound,
			reason:   "SUBJECT_NOT_FOUND",
			message:  "subject not found",
			metadata: map[string]string{"subjectId": "s1"},
			want: []proto.Message{
				&errdetails.ErrorInfo{Reason: "SUBJECT_NOT_FOUND", Domain: "test.dotto", Metadata: map[string]string{"subjectId": "s1"}},
			},
		},
		{
			name:    "details は ErrorInfo の後に順に付く",
			code:    connect.CodeInvalidArgument,
			reason:  "SUBJECT_INVALID",
			message: "invalid subject",
			details: []protoadapt.MessageV1{
				&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{Field: "id", Description: "is required"}}},
				&errdetails.LocalizedMessage{Locale: "ja-JP", Message: "科目が不正です"},
			},
			want: []proto.Message{
				&errdetails.ErrorInfo{Reason: "SUBJECT_INVALID", Domain: "test.dotto"},
				&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{{Field: "id", Description: "is required"}}},
				&errdetails.LocalizedMessage{Locale: "ja-JP", Message: "科目が不正です"},
			},
		},
		{
			name:    "空の BadRequest も付く",
			code:    connect.CodeInvalidArgument,
			reason:  "ZONE_INVALID",
			message: "zone has a value that violates a constraint",
			details: []protoadapt.MessageV1{&errdetails.BadRequest{}},
			want: []proto.Message{
				&errdetails.ErrorInfo{Reason: "ZONE_INVALID", Domain: "test.dotto"},
				&errdetails.BadRequest{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testDomain.New(tt.code, tt.reason, tt.message, tt.metadata, tt.details...)
			assertConnectError(t, err, tt.code, tt.message, tt.want...)
		})
	}
}

func TestViolations(t *testing.T) {
	type violation struct{ field, description string }
	tests := []struct {
		name       string
		violations []violation
		wantNil    bool
		wantMsg    string
		want       []*errdetails.BadRequest_FieldViolation
	}{
		{
			name:    "不正なフィールドがなければ nil",
			wantNil: true,
		},
		{
			name:       "1 件",
			violations: []violation{{"subject_id", "must be a valid ID"}},
			wantMsg:    "request has invalid fields: subject_id",
			want:       []*errdetails.BadRequest_FieldViolation{{Field: "subject_id", Description: "must be a valid ID"}},
		},
		{
			name: "複数件は追加した順",
			violations: []violation{
				{"from_date", "must be a valid date"},
				{"until_date", "must be a valid date"},
				{"faculties.0.faculty_id", "is required"},
			},
			wantMsg: "request has invalid fields: from_date, until_date, faculties.0.faculty_id",
			want: []*errdetails.BadRequest_FieldViolation{
				{Field: "from_date", Description: "must be a valid date"},
				{Field: "until_date", Description: "must be a valid date"},
				{Field: "faculties.0.faculty_id", Description: "is required"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v rpcerror.Violations
			for _, fv := range tt.violations {
				v.Add(fv.field, fv.description)
			}
			err := v.Err(testDomain)
			if tt.wantNil {
				if err != nil {
					t.Fatalf("Err() = %v, want nil", err)
				}
				return
			}
			assertConnectError(t, err, connect.CodeInvalidArgument, tt.wantMsg,
				&errdetails.ErrorInfo{Reason: "INVALID_ARGUMENT", Domain: "test.dotto"},
				&errdetails.BadRequest{FieldViolations: tt.want},
			)
		})
	}
}

func TestDomainContextError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantNil bool
		code    connect.Code
		message string
		reason  string
	}{
		{name: "DeadlineExceeded", err: context.DeadlineExceeded, code: connect.CodeDeadlineExceeded, message: "request deadline exceeded", reason: "DEADLINE_EXCEEDED"},
		{name: "ラップされた Canceled", err: fmt.Errorf("query: %w", context.Canceled), code: connect.CodeCanceled, message: "request canceled", reason: "CANCELED"},
		{name: "context のエラーでなければ nil", err: errors.New("boom"), wantNil: true},
		{name: "nil は nil", err: nil, wantNil: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testDomain.ContextError(tt.err)
			if tt.wantNil {
				if err != nil {
					t.Fatalf("ContextError() = %v, want nil", err)
				}
				return
			}
			assertConnectError(t, err, tt.code, tt.message,
				&errdetails.ErrorInfo{Reason: tt.reason, Domain: "test.dotto"},
			)
		})
	}
}
