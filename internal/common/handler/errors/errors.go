package errors

import (
	"fmt"

	"github.com/Hypocrite/gorder/common/consts"
	"github.com/pkg/errors"
)

type Error struct {
	code int
	msg  string
	err  error
}

func New(code int) error {
	return &Error{
		code: code,
	}
}

func NewWithError(code int, err error) error {
	if err == nil {
		return New(code)
	}
	return &Error{
		code: code,
		err:  err,
	}
}

func NewWithMsg(code int, format string, args ...any) error {
	return &Error{
		code: code,
		msg:  fmt.Sprintf(format, args...),
	}
}

func (e *Error) Error() string {
	var msg string
	if e.msg != "" {
		msg = e.msg
	}
	msg = consts.ErrMsg[e.code]
	return msg + " -> " + e.err.Error()
}

func Errcode(err error) int {
	if err == nil {
		return consts.ErrcodeSuccess
	}
	targetErr := &Error{}
	if errors.As(err, targetErr) {
		return targetErr.code
	}
	return -1
}

func Output(err error) (int, string) {
	if err == nil {
		return consts.ErrcodeSuccess, consts.ErrMsg[consts.ErrcodeSuccess]
	}

	errcode := Errcode(err)
	if errcode == -1 {
		return consts.ErrcodeUnknownError, err.Error()
	}
	return errcode, err.Error()
}
