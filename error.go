package tencent

import "fmt"

type ErrResponse struct {
	Err *Error `json:"error"`
}

func (e *ErrResponse) Error() string {
	return e.Err.Error()
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}
