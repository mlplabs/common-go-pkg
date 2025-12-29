package errors

import (
	"encoding/json"
	"errors"
	"github.com/mlplabs/common-go-pkg/pkg/http/errors/custom"
	"log"
	"net/http"
)

type CommonError interface {
	StatusCode() int
	ErrorCode() string
	Error() string
}

type CommonService interface {
	Service() string
}

type CommonErrorWithLog interface {
	LogError() error
}

type ResponseError struct {
	Error Response `json:"error"`
}

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Service string      `json:"service,omitempty"` // Need for understanding to get name service with error
}

func parseError(err error) (CommonError, string, error) {
	var logMsg error
	var e CommonError
	ok := errors.As(err, &e)
	if !ok {
		logMsg = err
		if ok := errors.As(err, &e); !ok {
			e = &custom.ServerError{}
		}
	}

	if logMsg == nil {
		if errorLog, ok := e.(CommonErrorWithLog); ok {
			logMsg = errorLog.LogError()
		}
	}

	var serviceName string
	if service, ok := e.(CommonService); ok {
		serviceName = service.Service()
	}
	return e, serviceName, logMsg
}

func SetError(w http.ResponseWriter, r *http.Request, err error) {
	clientServiceName := "this service"
	if r != nil {
		clientServiceName = r.Header.Get("X-Service-Name")
	}
	commonErr, serviceName, errLog := parseError(err)
	if errLog != nil {
		log.Printf("ERROR: client: %s, err: %v", clientServiceName, errLog)
	}
	w.WriteHeader(commonErr.StatusCode())
	data := ResponseError{
		Error: Response{
			Code:    commonErr.ErrorCode(),
			Message: commonErr.Error(),
			Data:    nil,
			Service: serviceName,
		},
	}
	body, err := json.Marshal(data)
	if err != nil {
		SetError(w, nil, err)
	}
	w.Write(body)
}
