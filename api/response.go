package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 9:36
 * Description: 全局异常
 */

type ErrorCode struct {
	HttpStatus int
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

// 定义错误码
var (
	ErrorInvalidParams = ErrorCode{http.StatusForbidden, 40000, "请求参数错误"}
	ErrorNotLogin      = ErrorCode{http.StatusUnauthorized, 40100, "未登录"}
	ErrorNoAuth        = ErrorCode{http.StatusUnauthorized, 40101, "无权限"}
	ErrorNotFound      = ErrorCode{http.StatusNotFound, 40400, "请求数据不存在"}
	ErrorFORBIDDEN     = ErrorCode{http.StatusForbidden, 40300, "禁止访问"}
	ErrorSYSTEM        = ErrorCode{http.StatusInternalServerError, 50000, "系统内部异常"}
	ErrorPERATION      = ErrorCode{http.StatusBadRequest, 50001, "操作失败"}
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Data:    data,
		Message: "ok",
	})
}

func ResponseFailed(c *gin.Context, err ErrorCode) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    err.Code,
		Message: err.Message,
	})
}

func ResponseErrorWithMsg(c *gin.Context, errorCode ErrorCode, message string) {
	c.JSON(errorCode.HttpStatus, &Response{
		Message: message,
		Code:    errorCode.Code,
		Data:    nil,
	})
}
