package builder

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sql_generate/core/schema"
	"sql_generate/global"
	"sql_generate/models"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/16
 * Time: 18:45
 * Description: 表概要生成器
 */

type TableSchemaBuilder struct {
	SQLDialect SQLDialect
}

func NewTableSchemaBuilder() *TableSchemaBuilder {
	return &TableSchemaBuilder{
		SQLDialect: &Dialect{},
	}
}

// BuildFromAuto 智能导入
func (b *TableSchemaBuilder) BuildFromAuto(content string) (*schema.TableSchema, error) {
	if len(content) == 0 || content == "" {
		return nil, fmt.Errorf("content is empty")
	}
	// 切分单词
	pattern := regexp.MustCompile("[,，]")
	words := pattern.Split(content, -1)
	if len(words) == 0 || len(words) > 20 {
		return nil, fmt.Errorf("content no compliance: %v", words)
	}
	// 根据单词去词库里匹配列信息，未匹配到的使用默认值
	var fieldInfoList []*models.FieldInfo
	if res := global.DB.Where("name in ?", words).Or("fieldName in ?", words).Find(&fieldInfoList); res.Error != nil {
		return nil, res.Error
	}
	// 名称 => 字段信息
	// 名称 => 字段信息
	nameFieldInfoMap := make(map[string][]*models.FieldInfo)
	fieldNameFieldInfoMap := make(map[string][]*models.FieldInfo)
	for _, fieldInfo := range fieldInfoList {
		nameFieldInfoMap[fieldInfo.Name] = append(nameFieldInfoMap[fieldInfo.Name], fieldInfo)
		fieldNameFieldInfoMap[fieldInfo.FieldName] = append(fieldNameFieldInfoMap[fieldInfo.FieldName], fieldInfo)
	}
	tableSchema := &schema.TableSchema{}
	tableSchema.TableName = "my_table"
	tableSchema.TableComment = "自动生成的表"
	fieldList := make([]schema.Field, 0)
	for _, word := range words {
		var field schema.Field
		var infoList []*models.FieldInfo
		infoList = append(infoList, nameFieldInfoMap[word]...)
		infoList = append(infoList, fieldNameFieldInfoMap[word]...)

		if len(infoList) > 0 {
			if err := json.Unmarshal([]byte(infoList[0].Content), &field); err != nil {
				return nil, fmt.Errorf("cannot unmarshal content: %v", err)
			}
		} else {
			field = getDefaultField(word)
		}

		fieldList = append(fieldList, field)
	}
	tableSchema.FieldList = fieldList
	return tableSchema, nil
}

// 获取默认字段
func getDefaultField(word string) schema.Field {
	field := schema.Field{
		FieldName:     word,
		FieldType:     "text",
		DefaultValue:  "",
		NotNull:       false,
		Comment:       word,
		PrimaryKey:    false,
		AutoIncrement: false,
		MockType:      "",
		MockParams:    "",
		OnUpdate:      "",
	}
	return field
}
