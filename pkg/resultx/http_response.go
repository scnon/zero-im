package resultx

import (
	"context"
	"github.com/pkg/errors"
	zrpcErr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"net/http"
	"zero-im/pkg/xerr"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "",
		Data: data,
	}
}

func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}

func OkHandler(_ context.Context, v interface{}) any {
	return Success(v)
}

func ErrHandler(name string) func(err error) (int, any) {
	return func(err error) (int, any) {
		errCode := xerr.SERVER_COMMON_ERROR
		errMsg := xerr.ErrMsg(errCode)

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*zrpcErr.CodeMsg); ok {
			errCode = e.Code
			errMsg = e.Msg
		} else {
			if gs, ok := status.FromError(causeErr); ok {
				errCode = int(gs.Code())
				errMsg = gs.Message()
			}
		}

		return http.StatusBadRequest, Fail(errCode, errMsg)
	}
}
