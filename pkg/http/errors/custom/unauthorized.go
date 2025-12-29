package custom

type Unauthorized struct {
	err error
}

func (*Unauthorized) StatusCode() int {
	return 401
}

func (*Unauthorized) ErrorCode() string {
	return "UNAUTHORIZED"
}

func (e *Unauthorized) Error() string {
	return e.err.Error()
}

func NewUnauthorized(err error) *Unauthorized {
	return &Unauthorized{err: err}
}
