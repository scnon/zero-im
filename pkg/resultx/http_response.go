package resultx

import (
	"context"
	"net/http"
	"zero-im/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zrpcErr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
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

func ErrHandler(name string) func(c context.Context, err error) (int, any) {
	return func(c context.Context, err error) (int, any) {
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

		// 日志记录
		logx.WithContext(c).Errorf("【%s】 err %v", name, err)

		return http.StatusBadRequest, Fail(errCode, errMsg)
	}
}
