package generator

import (
	"context"
	"math/rand"
	"sql_generate/client"
	"sql_generate/core/schema"
	"time"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/14
 * Time: 10:10
 * Description: 词库数据生成器
 */

type DictDataGenerator struct {
	DictServiceInterface DictServiceInterface
}

type DictServiceInterface interface {
	GetWordList(context.Context, int64) ([]string, error)
}

func NewDictDataGenerator() *DictDataGenerator {
	return &DictDataGenerator{
		DictServiceInterface: client.DictService{},
	}
}

func (r *DictDataGenerator) DoGenerate(field schema.Field, rowNum int32) ([]string, error) {
	mockParams := field.MockParams.(float64)
	wordList, err := r.DictServiceInterface.GetWordList(context.Background(), int64(mockParams))
	if err != nil {
		return nil, err
	}
	list := make([]string, 0, rowNum)
	rand.Seed(time.Now().Unix())
	for i := 0; i < int(rowNum); i++ {
		randomStr := wordList[rand.Intn(len(wordList))]
		list = append(list, randomStr)
	}
	return list, nil
}
