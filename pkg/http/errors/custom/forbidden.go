package custom

type Forbidden struct {
	err error
}

func (*Forbidden) StatusCode() int {
	return 403
}

func (*Forbidden) ErrorCode() string {
	return "FORBIDDEN"
}

func (e *Forbidden) Error() string {
	return e.err.Error()
}

func NewForbidden(err error) *Forbidden {
	return &Forbidden{err: err}
}
