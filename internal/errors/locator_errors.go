package errors

const (
	GeneralLocatorError ErrorCode = 1
	AlreadyInitialized  ErrorCode = 2
)

type LocatorError struct {
	code    ErrorCode
	message string
	err     error
}

func NewLocatorError(err error) *LocatorError {
	return &LocatorError{
		code:    GeneralLocatorError,
		message: "an error occurred",
		err:     err,
	}
}

func NewAlreadyInitializedError(err error) *LocatorError {
	return &LocatorError{
		code:    AlreadyInitialized,
		message: "locator was already initialized",
		err:     err,
	}
}

func (e *LocatorError) Message() string {
	return e.message
}

func (e *LocatorError) Code() ErrorCode {
	return e.code
}

func (e *LocatorError) GoError() error {
	return e.err
}
