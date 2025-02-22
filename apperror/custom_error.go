package apperror

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func NewCustomError(code int, message string, err error) *CustomError {
	if err != nil {
		message = message + ": " + err.Error()
	}
	return &CustomError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

func (ce *CustomError) Error() string {
	return ce.Message
}

func (ce *CustomError) Unwrap() error {
	return ce.Cause
}
