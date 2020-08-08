package errors

type (
	ErrorCode = string
)

const (
	CodeUnexpect         ErrorCode = "UNEXPECT"
	CodeObjectNotFound   ErrorCode = "OBJECT_NOT_FOUND"
	CodeInvalidArguments ErrorCode = "INVALID_ARGUMENTS"
)
