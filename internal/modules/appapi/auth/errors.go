package auth

import (
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

const errorDomain = "app.dotto"

func newError(code connect.Code, reason, message string) *connect.Error {
	err := connect.NewError(code, errors.New(message))
	if d, derr := connect.NewErrorDetail(&errdetails.ErrorInfo{Reason: reason, Domain: errorDomain}); derr == nil {
		err.AddDetail(d)
	}
	return err
}
