package errors

type Validation struct {
	Error
}

func (v Validation) GetError(message string, status int) Error {
	return Error{
		Code:       "VALIDATION_ERROR",
		Message:    message,
		Details:    GetDetails(),
		StatusCode: status,
	}
}
