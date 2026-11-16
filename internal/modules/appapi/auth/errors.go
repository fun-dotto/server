package auth

import (
	"connectrpc.com/connect"
	"github.com/fun-dotto/server/internal/shared/rpcerror"
)

const errorDomain rpcerror.Domain = "app.dotto"

// newError は ErrorInfo（metadata なし）を付けた Connect のエラーを返す。
func newError(code connect.Code, reason, message string) *connect.Error {
	return errorDomain.New(code, reason, message, nil)
}
