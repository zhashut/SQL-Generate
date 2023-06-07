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
 * Description: No Description
 */

type Response struct {
	Code    int32       `json:"code"`
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

func ResponseFailed(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    40001,
		Message: err.Error(),
	})
	return
}
