package tester

import "fmt"

type ResponseError struct {
	Message string
	Err     error
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("%s. Error: %v", r.Message, r.Err)
}

type PostDataFileError struct {
	FileName string
	Err      error
}

func (r *PostDataFileError) Error() string {
	return fmt.Sprintf("[%s], %v", r.FileName, r.Err)
}

type HttpCodeError struct {
	Err error
}

func (r *HttpCodeError) Error() string {
	return fmt.Sprintf("Error: %v", r.Err)
}
