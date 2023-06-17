package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sql_generate/core"
	"sql_generate/core/builder"
	"sql_generate/core/schema"
	"sql_generate/models"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 9:32
 * Description: 数据生成器api
 */

var (
	generator          = core.NewGeneratorFace()
	tableSchemaBuilder = builder.NewTableSchemaBuilder()
)

func GenerateSQL(c *gin.Context) {
	var tableSchema schema.TableSchema
	if err := c.ShouldBind(&tableSchema); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
	}
	all, err := generator.GenerateAll(&tableSchema)
	zap.S().Infof("GenerateAll.result: %v", all)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
	}

	ResponseSuccess(c, all)
}

// GetSchemaByAuto 智能导入
func GetSchemaByAuto(c *gin.Context) {
	var req models.GenerateByAutoRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	auto, err := tableSchemaBuilder.BuildFromAuto(req.Content)
	if err != nil {
		ResponseFailed(c, ErrorPERATION)
		return
	}
	ResponseSuccess(c, auto)
}

// GetSchemaByExcel Excel导入
func GetSchemaByExcel(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	defer file.Close()

	tableSchema, err := tableSchemaBuilder.BuildFromExcel(file)

	ResponseSuccess(c, tableSchema)
}
