package custom

type ServerError struct {
	err error
}

func (e *ServerError) LogError() error {
	return e.err
}

func (e *ServerError) ErrorCode() string {
	return "SERVER_UNEXPECTED"
}

func (*ServerError) StatusCode() int {
	return 500
}

func (e *ServerError) Error() string {
	return "Ошибка сервиса"
}

func NewServerError(err error) *ServerError {
	return &ServerError{err: err}
}
