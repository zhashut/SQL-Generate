package generator

import "sql_generate/core/schema"

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/6
 * Time: 17:49
 * Description: No Description
 */

type RandomDataGenerator struct{}

func NewRandomDataGenerator() *RandomDataGenerator {
	return &RandomDataGenerator{}
}

// DoGenerate TODO
func (r *RandomDataGenerator) DoGenerate(field schema.Field, rowNum int32) []string {

}
