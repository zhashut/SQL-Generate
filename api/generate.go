package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"net/url"
	"sql_generate/core"
	"sql_generate/core/builder"
	"sql_generate/core/schema"
	"sql_generate/models"
	"strings"
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

// DownloadDataExcel 下载模拟数据 Excel
func DownloadDataExcel(c *gin.Context) {
	var generateVO models.Generate
	err := json.NewDecoder(c.Request.Body).Decode(&generateVO)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, "Invalid request body")
		return
	}

	tableSchema := generateVO.TableSchema
	tableName := tableSchema.TableName

	// Create a new Excel workbook
	f := excelize.NewFile()

	// Set table headers
	for index, field := range tableSchema.FieldList {
		err = f.SetCellValue("Sheet1", fmt.Sprintf("%s1", columnIndexToExcelColumn(index)), field.FieldName)
		if err != nil {
			ResponseErrorWithMsg(c, ErrorPERATION, "Error setting header value")
			return
		}
	}

	// Set data rows
	for rowIndex, data := range generateVO.DataList {
		for colIndex, field := range tableSchema.FieldList {
			err = f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", columnIndexToExcelColumn(colIndex), rowIndex+2), data[field.FieldName])
			if err != nil {
				ResponseErrorWithMsg(c, ErrorPERATION, "Error setting cell value")
				return
			}
		}
	}

	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	charset := "utf-8"

	fileName := tableName + "表数据"
	disposition := "attachment;filename*=utf-8''" + url.QueryEscape(strings.ReplaceAll(fileName, " ", "%20")) + ".xlsx"

	c.Header("Content-Disposition", disposition)
	c.Header("Content-Type", contentType+";charset="+charset)
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")

	err = f.Write(c.Writer)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorInvalidParams, "Error writing Excel file")
	}
}
