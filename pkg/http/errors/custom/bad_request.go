package custom

type BadRequest struct {
	err error
}

func (*BadRequest) StatusCode() int {
	return 400
}

func (*BadRequest) ErrorCode() string {
	return "BAD_REQUEST"
}

func (e *BadRequest) Error() string {
	return e.err.Error()
}

func NewBadRequest(err error) *BadRequest {
	return &BadRequest{err: err}
}
