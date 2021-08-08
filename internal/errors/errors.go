package errors

type ErrorCode int

const (
	GeneralError            ErrorCode = 1
	InterruptError          ErrorCode = 2
	InstructionMissingError ErrorCode = 3
)

type PromptError struct {
	code    ErrorCode
	message string
	err     error
}

func New(err error) *PromptError {
	return &PromptError{
		code:    GeneralError,
		message: "an error occurred",
		err:     err,
	}
}

func NewInstructionMissingError(err error) *PromptError {
	return &PromptError{
		code:    InstructionMissingError,
		message: "instruction file is missing",
		err:     err,
	}
}

func NewInterruptError(err error) *PromptError {
	return &PromptError{
		code:    InterruptError,
		message: "keyboard interrupt identified",
		err:     err,
	}
}

func (e *PromptError) Message() string {
	return e.message
}

func (e *PromptError) Code() ErrorCode {
	return e.code
}

func (e *PromptError) GoError() error {
	return e.err
}
