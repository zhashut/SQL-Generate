package builder

import (
	. "sql_generate/consts"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/19
 * Time: 1:25
 * Description: 基于模板模式，生成各类代码实体
 */

// CodeBuilder code generator
type CodeBuilder struct {
	CodeGenerator *CodeGenerator
}

func NewCodeGenerator(language, template string, index int) *CodeBuilder {
	return &CodeBuilder{
		CodeGenerator: &CodeGenerator{
			Language: language,
			Template: template,
			FileTypeHandler: func(fieldType string) string {
				return FieldTypeEnumToString[GetFieldTypeEnumByValue(fieldType)][index]
			},
		},
	}
}
