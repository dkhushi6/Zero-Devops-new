package helper

import (
	"Zero_Devops/server/domain"
	"net/http"
	rtdebug "runtime/debug"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func BuildErrorResponse(message string,err error,reqId string , opts ...DebugOption) domain.ErrorResponse{
	resp := domain.ErrorResponse{
		Success: false,
		Error: domain.ErrorBody{
			Code:    getStatusCode(err),
			Message: message,
		},
		RequestId: reqId,
	}

	if viper.GetString("APP_ENV") != "production"{
		debug := &domain.DebugInfo{
			RawError: err.Error(),
			Stack: string(rtdebug.Stack()),
		}
		for _, opt := range opts {
			opt(debug)
		}
		resp.Error.Debug = debug
	}
	
	return resp
}

type DebugOption func(* domain.DebugInfo)

func WithReason(r string) DebugOption {
    return func(d * domain.DebugInfo) { d.Reason = r }
}

func WithQuery(q string) DebugOption {
    return func(d * domain.DebugInfo) { d.Query = q }
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}