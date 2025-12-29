package custom

type ErrorNoRows struct {
	err error
}

func (*ErrorNoRows) StatusCode() int {
	return 404
}

func (*ErrorNoRows) ErrorCode() string {
	return "OBJECT_DOES_NOT_EXIST"
}

func (e *ErrorNoRows) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return "Объект не существует"
}
func NewErrorNoRows(err error) *ErrorNoRows {
	return &ErrorNoRows{err: err}
}
