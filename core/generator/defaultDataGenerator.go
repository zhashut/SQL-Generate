package generator

import (
	"sql_generate/core/schema"
	"strconv"
	"time"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 20:01
 * Description: 默认数据生成器
 */

const (
	CURRENT_TIMESTAMP = "CURRENT_TIMESTAMP"
)

type DefaultDataGenerator struct{}

func NewDefaultDataGenerator() *DefaultDataGenerator {
	return &DefaultDataGenerator{}
}

func (r *DefaultDataGenerator) DoGenerate(field *schema.Field, rowNum int32) ([]string, error) {
	mockParams := field.MockParams.(string)
	result := make([]string, 0, rowNum)
	// 主键采用递增策略
	if field.PrimaryKey {
		if mockParams == "" {
			mockParams = "1"
		}
		initValue, _ := strconv.Atoi(mockParams)
		for i := 0; i < int(rowNum); i++ {
			result = append(result, strconv.Itoa(initValue+i))
		}
		return result, nil
	}
	// 使用默认值
	defaultValue := field.DefaultValue
	// 特殊逻辑，日期瑶伪造数据
	if defaultValue == CURRENT_TIMESTAMP {
		defaultValue = time.Now().Format("2006-01-02 15:04:05")
	}
	if defaultValue != "" {
		for i := 0; i < int(rowNum); i++ {
			result = append(result, defaultValue)
		}
	}
	return result, nil
}
