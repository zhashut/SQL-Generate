package generator

import (
	"sql_generate/consts"
	"sql_generate/core/schema"
	"sql_generate/utils"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/6
 * Time: 17:49
 * Description: 随机值数据生成器
 */

type RandomDataGenerator struct{}

func NewRandomDataGenerator() *RandomDataGenerator {
	return &RandomDataGenerator{}
}

func (r *RandomDataGenerator) DoGenerate(field *schema.Field, rowNum int32) ([]string, error) {
	mockParams := field.MockParams.(string)
	result := make([]string, 0, rowNum)
	for i := 0; i < int(rowNum); i++ {
		randomTypeEnum := consts.GetMockParamsRandomTypeByValue(mockParams)
		randomString := utils.GetRandomValue(randomTypeEnum)
		result = append(result, randomString)
	}
	return result, nil
}
