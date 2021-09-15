package errors

type ErrorCode int

const (
	GeneralPromptError      ErrorCode = 10
	InterruptError          ErrorCode = 20
	InstructionMissingError ErrorCode = 30
	SchemaError             ErrorCode = 40
)

type PromptError struct {
	code    ErrorCode
	message string
	err     error
}

func NewPromptError(err error) *PromptError {
	return &PromptError{
		code:    GeneralPromptError,
		message: "an error occurred",
		err:     err,
	}
}

func NewSchemaError(err error) *PromptError {
	return &PromptError{
		code:    SchemaError,
		message: "invalid schema",
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
