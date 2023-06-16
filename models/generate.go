package models

import (
	"sql_generate/core/schema"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/5
 * Time: 15:10
 * Description: 生成的返回值
 */

type Generate struct {
	TableSchema *schema.TableSchema      `json:"tableSchema"`
	CreateSQL   string                   `json:"createSql"`
	DataList    []map[string]interface{} `json:"dataList"`
	InsertSQL   string                   `json:"insertSql"`
	DataJson    string                   `json:"dataJson"`
}

type GenerateByAutoRequest struct {
	Content string `json:"content"`
}
