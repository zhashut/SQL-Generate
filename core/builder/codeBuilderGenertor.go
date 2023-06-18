package builder

import (
	"bytes"
	"html/template"
	"sql_generate/core/schema"
	"sql_generate/models"
	"sql_generate/utils"
	"strings"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/18
 * Time: 21:33
 * Description: 模板方法设计模式。我们可以创建一个基本代码生成器类，
 *				其通用功能将在一个基本方法中实现，然后通过子类来定制不同的语言生成
 */

type CodeGenerator struct {
	Language        string
	Template        string
	FileTypeHandler func(fieldType string) string
}

func (c *CodeGenerator) BuilderCode(tableSchema *schema.TableSchema) (string, error) {
	entityDTO := models.EntityDTO{
		EntityName: strings.Title(utils.CamelCase(tableSchema.TableName)),
		Comment:    tableSchema.TableComment,
	}

	if entityDTO.Comment == "" {
		entityDTO.Comment = entityDTO.EntityName
	}

	entityDTO.FieldList = c.buildFields(tableSchema.FieldList)
	tmpl, err := template.New(c.Language).Parse(c.Template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, entityDTO); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *CodeGenerator) buildFields(fieldSchemas []*schema.Field) []models.FieldDTO {
	var fields []models.FieldDTO
	for _, field := range fieldSchemas {
		fieldType := c.FileTypeHandler(field.FieldType)
		fieldName := utils.CamelCase(field.FieldName)
		fields = append(fields, models.FieldDTO{
			Comment: field.Comment,
			Type:    fieldType,
			Name:    fieldName,
		})
	}
	return fields
}
