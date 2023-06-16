package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sql_generate/respository/db"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/14
 * Time: 11:55
 * Description: 获取词条链表
 */

type DictService struct {
	DB *db.DictDao
}

func (s DictService) GetWordList(ctx context.Context, id int64) ([]string, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	dict, err := s.DB.GetDictByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get dict: %v", err)
	}

	var wordList []string
	if err := json.Unmarshal([]byte(dict.Content), &wordList); err != nil {
		return nil, fmt.Errorf("cannot unmarshal content: %v", err)
	}

	return wordList, nil
}
