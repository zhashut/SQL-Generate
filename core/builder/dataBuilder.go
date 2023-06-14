package builder

import (
	"sql_generate/consts"
	"sql_generate/core/generator"
	"sql_generate/core/schema"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/6
 * Time: 16:13
 * Description: 数据生成器
 */

func GenerateData(tableSchema *schema.TableSchema, rowNum int32) ([]map[string]interface{}, error) {
	fieldList := tableSchema.FieldList
	// 初始化结果数据
	resultList := make([]map[string]interface{}, rowNum)
	for i := range resultList {
		resultList[i] = make(map[string]interface{})
	}
	// 依次生成每一列
	for _, field := range fieldList {
		mockTypeEnum := consts.GetMockTypeEnumByValue(field.MockType)
		dataGenerator := generator.GetGenerator(mockTypeEnum)
		mockDataSlice, err := dataGenerator.DoGenerate(field, rowNum)
		if err != nil {
			return nil, err
		}
		fieldName := field.FieldName
		if len(mockDataSlice) > 0 {
			for i := 0; i < int(rowNum); i++ {
				resultList[i][fieldName] = mockDataSlice[i]
			}
		}
	}
	return resultList, nil
}
