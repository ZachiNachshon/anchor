package errors

const (
	GeneralExtractorError ErrorCode = 1
	EmptyActionsError     ErrorCode = 2
	EmptyWorkflowsError   ErrorCode = 3
)

type ExtractorError struct {
	code    ErrorCode
	message string
	err     error
}

func NewExtractorError(err error) *ExtractorError {
	return &ExtractorError{
		code:    GeneralExtractorError,
		message: "an error occurred",
		err:     err,
	}
}

func NewEmptyActionsError(err error) *ExtractorError {
	return &ExtractorError{
		code:    EmptyActionsError,
		message: "empty actions",
		err:     err,
	}
}

func NewEmptyWorkflowsError(err error) *ExtractorError {
	return &ExtractorError{
		code:    EmptyWorkflowsError,
		message: "empty workflows",
		err:     err,
	}
}

func (e *ExtractorError) Message() string {
	return e.message
}

func (e *ExtractorError) Code() ErrorCode {
	return e.code
}

func (e *ExtractorError) GoError() error {
	return e.err
}
