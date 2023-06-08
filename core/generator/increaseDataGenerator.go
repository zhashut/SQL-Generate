package generator

import (
	"sql_generate/core/schema"
	"strconv"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 20:30
 * Description: 递增值数据生成器
 */

type IncreaseDataGenerator struct{}

func NewIncreaseDataGenerator() *IncreaseDataGenerator {
	return &IncreaseDataGenerator{}
}

func (r *IncreaseDataGenerator) DoGenerate(field schema.Field, rowNum int32) []string {
	mockParams := field.MockParams
	result := make([]string, 0, rowNum)
	if mockParams == "" {
		mockParams = "1"
	}
	initValue, _ := strconv.Atoi(mockParams)
	for i := 0; i < int(rowNum); i++ {
		result = append(result, strconv.Itoa(initValue+i))
	}
	return result
}
