package vo

type ErrorWithHttpStatus interface {
	error
	HttpStatus() int
}

type errorWithHttpStatus struct {
	error      string
	httpStatus int
}

func NewErrorWithHttpStatus(error string, httpStatus int) ErrorWithHttpStatus {
	return errorWithHttpStatus{error: error, httpStatus: httpStatus}
}

func (e errorWithHttpStatus) Error() string {
	return e.error
}

func (e errorWithHttpStatus) HttpStatus() int {
	return e.httpStatus
}
