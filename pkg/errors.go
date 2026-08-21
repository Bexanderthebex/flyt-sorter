package pkg

import "fmt"

type ErrorCode string

const (
	ErrCodeNoStartingPath       ErrorCode = "no_starting_path"
	ErrCodeMultipleStartingPath ErrorCode = "multiple_starting_path"
)

const (
	ErrMsgNoStartingPath       = "No possible starting flight path detected"
	ErrMsgMultipleStartingPath = "Multiple starting flight path detected"
)

type InvalidParameterError struct {
	Message string
	Code    ErrorCode
}

func (e InvalidParameterError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
