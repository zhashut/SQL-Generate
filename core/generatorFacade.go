package core

import (
	"fmt"
	"sql_generate/core/builder"
	"sql_generate/core/models"
	"sql_generate/core/schema"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/5
 * Time: 14:46
 * Description: No Description
 */

type GeneratorFace struct {
}

func (g *GeneratorFace) GenerateAll(tableSchema *schema.TableSchema) (*models.Generate, error) {
	// 校验
	if err := validSchema(tableSchema); err != nil {
		return nil, err
	}
	sqlBuilder := builder.NewSQLBuilder()
	createSQL, err := sqlBuilder.BuildCreateTableSql(tableSchema)
	if err != nil {
		return nil, err
	}
	mockNum := tableSchema.MockNum
	// 生成模拟数据
	dataList := builder.GenerateData(tableSchema, mockNum)
	// 生成插入 SQL

}

func validSchema(tableSchema *schema.TableSchema) error {
	if tableSchema == nil {
		return fmt.Errorf("failed valid schema")
	}
	tableName := tableSchema.TableName
	if tableName == "" {
		return fmt.Errorf("表名不能为空")
	}
	mockNum := tableSchema.MockNum
	// 默认生成 20 条
	if mockNum == 0 {
		mockNum = 20
	}
	if mockNum > 100 || mockNum < 10 {
		return fmt.Errorf("生成条数设置错误")
	}
	fieldList := tableSchema.FieldList
	if len(fieldList) == 0 {
		return fmt.Errorf("字段列表不能为空")
	}
	for _, field := range fieldList {
		err := validField(field)
		if err != nil {
			return err
		}
	}
	return nil
}

func validField(field schema.Field) error {
	fieldName := field.FieldName
	fieldType := field.FieldType
	if fieldName == "" {
		return fmt.Errorf("字段名不能为空")
	}
	if fieldType == "" {
		return fmt.Errorf("字段类型不能为空")
	}
	return nil
}
