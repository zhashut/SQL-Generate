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
 * Description: 数据生成器工厂
 * 工厂 + 单例模式，降低开销
 */

type DataGenerator interface {
	DoGenerate(field schema.Field, rowNum int32) []string
}

var mockTypeDataGeneratorMap map[MockTypeEnum]DataGenerator

func init() {
	/**
	 * 模拟类型 => 生成器映射
	 */
	mockTypeDataGeneratorMap = make(map[MockTypeEnum]DataGenerator)
	mockTypeDataGeneratorMap[RANDOM] = NewRandomDataGenerator()
	mockTypeDataGeneratorMap[FIXED] = NewFixedDataGenerator()
	mockTypeDataGeneratorMap[NONE] = NewDefaultDataGenerator()
	mockTypeDataGeneratorMap[INCREASE] = NewIncreaseDataGenerator()
}

// GetGenerator 获取实例
func GetGenerator(mockTypeEnum MockTypeEnum) DataGenerator {
	return mockTypeDataGeneratorMap[mockTypeEnum]
}
