package custom

type CommonError struct {
	err        error  // error what will be logged in our service
	errorCode  string // error code for user
	statusCode int    // http status code
	message    string // message for user
	service    string // service name
}

func NewCommonError(statusCode int, errorCode string, err error, message string, service string) *CommonError {
	return &CommonError{
		err:        err,
		errorCode:  errorCode,
		statusCode: statusCode,
		message:    message,
		service:    service,
	}
}

func (e *CommonError) LogError() error {
	return e.err
}

func (e *CommonError) ErrorCode() string {
	return e.errorCode
}

func (e *CommonError) StatusCode() int {
	return e.statusCode
}

func (e *CommonError) Error() string {
	return e.message
}

func (e *CommonError) Service() string {
	return e.service
}
