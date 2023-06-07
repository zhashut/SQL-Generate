package generator

import (
	. "sql_generate/consts"
	"sql_generate/core/schema"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/6
 * Time: 17:45
 * Description: No Description
 */

type DataGenerator interface {
	DoGenerate(field schema.Field, rowNum int32) []string
}

var mockTypeDataGeneratorMap map[MockTypeEnum]DataGenerator

// TODO 这个地方可能会报错，因为这里面的配置有依赖其他配置
func init() {
	mockTypeDataGeneratorMap = make(map[MockTypeEnum]DataGenerator)
	mockTypeDataGeneratorMap[RANDOM] = NewRandomDataGenerator()
}

// GetGenerator 获取实例
func GetGenerator(mockTypeEnum MockTypeEnum) DataGenerator {
	return mockTypeDataGeneratorMap[mockTypeEnum]
}
