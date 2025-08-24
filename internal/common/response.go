package common

import (
	"net/http"

	"github.com/Hypocrite/gorder/common/tracing"
	"github.com/gin-gonic/gin"
)

type BaseResponse struct{}
type response struct {
	ErrCode int    `json:"err_code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
}

func (base *BaseResponse) Response(c *gin.Context, err error, data interface{}) {
	if err != nil {
		base.error(c, err)
	} else {
		base.success(c, data)
	}
}

func (base *BaseResponse) success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, response{
		ErrCode: 0,
		Message: "success",
		Data:    data,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
}

func (base *BaseResponse) error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, response{
		ErrCode: 2,
		Message: err.Error(),
		Data:    nil,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
}
