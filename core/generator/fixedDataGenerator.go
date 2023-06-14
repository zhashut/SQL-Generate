package generator

import "sql_generate/core/schema"

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 19:59
 * Description: 固定值数据生成器
 */

type FixedDataGenerator struct{}

func NewFixedDataGenerator() *FixedDataGenerator {
	return &FixedDataGenerator{}
}

func (r *FixedDataGenerator) DoGenerate(field schema.Field, rowNum int32) ([]string, error) {
	mockParams := field.MockParams.(string)
	if mockParams == "" {
		mockParams = "6"
	}
	result := make([]string, 0, rowNum)
	for i := 0; i < int(rowNum); i++ {
		result = append(result, mockParams)
	}
	return result, nil
}
