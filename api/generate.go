package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sql_generate/core"
	"sql_generate/core/schema"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 9:32
 * Description: 数据生成器api
 */

func GenerateSQL(c *gin.Context) {
	var tableSchema schema.TableSchema
	if err := c.ShouldBind(&tableSchema); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
	}
	generator := core.NewGeneratorFace()
	all, err := generator.GenerateAll(&tableSchema)
	zap.S().Infof("GenerateAll.result: %v", all)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
	}

	ResponseSuccess(c, all)
}
