// Package rpcerror は Connect RPC のエラーを組み立てる共通のヘルパーを提供する。
//
// エラーの規約（.rulesync/rules/api-design.md のエラー）に従い、すべてのエラーに
// google.rpc.ErrorInfo を 1 つ付け、INVALID_ARGUMENT には google.rpc.BadRequest を付ける。
// ErrorInfo.domain は API ごとに異なるため、各 API は自身の Domain を定義して使う。
package rpcerror

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

// Domain は ErrorInfo.domain（例: admin.dotto）。
type Domain string

// New は ErrorInfo（reason・domain・metadata）と details を付けた Connect のエラーを返す。
func (d Domain) New(code connect.Code, reason, message string, metadata map[string]string, details ...protoadapt.MessageV1) *connect.Error {
	err := connect.NewError(code, errors.New(message))
	if detail, derr := connect.NewErrorDetail(&errdetails.ErrorInfo{
		Reason:   reason,
		Domain:   string(d),
		Metadata: metadata,
	}); derr == nil {
		err.AddDetail(detail)
	}
	for _, detail := range details {
		if ed, derr := connect.NewErrorDetail(protoadapt.MessageV2Of(detail)); derr == nil {
			err.AddDetail(ed)
		}
	}
	return err
}

// ContextError は context のキャンセル・タイムアウトを Connect のエラーに変換する。
// err がどちらでもない場合は nil を返す。
func (d Domain) ContextError(err error) error {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return d.New(connect.CodeDeadlineExceeded, "DEADLINE_EXCEEDED", "request deadline exceeded", nil)
	case errors.Is(err, context.Canceled):
		return d.New(connect.CodeCanceled, "CANCELED", "request canceled", nil)
	default:
		return nil
	}
}

// Violations は INVALID_ARGUMENT の BadRequest.field_violations を集める。
type Violations []*errdetails.BadRequest_FieldViolation

// Add は不正なフィールドを追加する。
func (v *Violations) Add(field, description string) {
	*v = append(*v, &errdetails.BadRequest_FieldViolation{Field: field, Description: description})
}

// Err は集めた不正なフィールドを INVALID_ARGUMENT のエラーにして返す。
// 不正なフィールドがない場合は nil を返す。
func (v Violations) Err(d Domain) error {
	if len(v) == 0 {
		return nil
	}
	fields := make([]string, 0, len(v))
	for _, fv := range v {
		fields = append(fields, fv.GetField())
	}
	return d.New(
		connect.CodeInvalidArgument,
		"INVALID_ARGUMENT",
		"request has invalid fields: "+strings.Join(fields, ", "),
		nil,
		&errdetails.BadRequest{FieldViolations: v},
	)
}
