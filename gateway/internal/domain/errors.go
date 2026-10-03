package domain

type FieldError struct {
	Field string
	Error string
}

type ValidationError struct {
	Message string
	Errors  []FieldError
}

func (e *ValidationError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}
