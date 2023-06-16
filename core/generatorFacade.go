package core

import (
	"fmt"
	"sql_generate/core/builder"
	"sql_generate/core/schema"
	"sql_generate/models"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/5
 * Time: 14:46
 * Description: No Description
 */

type GeneratorFace struct {
	SQLBuilder  *builder.SQLBuilder
	DataBuilder *builder.DataBuilder
	JsonBuilder *builder.JsonBuilder
}

func NewGeneratorFace() *GeneratorFace {
	return &GeneratorFace{
		SQLBuilder:  builder.NewSQLBuilder(),
		DataBuilder: builder.NewDataBuilder(),
		JsonBuilder: builder.NewJsonBuilder(),
	}
}

// GenerateAll 生成所有内容
func (g *GeneratorFace) GenerateAll(tableSchema *schema.TableSchema) (*models.Generate, error) {
	// 校验
	if err := g.ValidSchema(tableSchema); err != nil {
		return nil, err
	}
	// 生成创建 SQL
	createSQL, err := g.SQLBuilder.BuildCreateTableSql(tableSchema)
	if err != nil {
		return nil, err
	}
	mockNum := tableSchema.MockNum
	// 生成模拟数据
	dataList, err := g.DataBuilder.GenerateData(tableSchema, mockNum)
	if err != nil {
		return nil, err
	}
	// 生成插入 SQL
	insertSQL, err := g.SQLBuilder.BuildInsertSQL(tableSchema, dataList)
	if err != nil {
		return nil, err
	}
	// 生成数据 Json
	dataJson, err := g.JsonBuilder.BuildJSON(dataList)
	if err != nil {
		return nil, fmt.Errorf("dataJson 生成失败")
	}
	// 封装返回
	return &models.Generate{
		TableSchema: tableSchema,
		CreateSQL:   createSQL,
		DataList:    dataList,
		InsertSQL:   insertSQL,
		DataJson:    dataJson,
	}, nil
}

// ValidSchema 校验 schema
func (g *GeneratorFace) ValidSchema(tableSchema *schema.TableSchema) error {
	if tableSchema == nil {
		return fmt.Errorf("failed valid schema")
	}
	tableName := tableSchema.TableName
	if tableName == "" {
		return fmt.Errorf("表名不能为空")
	}
	mockNum := tableSchema.MockNum
	// 默认生成 20 条
	if tableSchema.MockNum == 0 {
		tableSchema.MockNum = 20
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
		err := g.ValidField(field)
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidField 检验字段
func (g *GeneratorFace) ValidField(field schema.Field) error {
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
