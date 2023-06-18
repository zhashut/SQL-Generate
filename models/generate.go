package models

import (
	"sql_generate/core/schema"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/5
 * Time: 15:10
 * Description: 生成器结构体
 */

type Generate struct {
	TableSchema    *schema.TableSchema      `json:"tableSchema"`
	CreateSQL      string                   `json:"createSql"`
	DataList       []map[string]interface{} `json:"dataList"`
	InsertSQL      string                   `json:"insertSql"`
	DataJson       string                   `json:"dataJson"`
	JavaEntityCode string                   `json:"javaEntityCode"`
	GoStructCode   string                   `json:"goStructCode"`
}

type GenerateByAutoRequest struct {
	Content string `json:"content"`
}

// EntityDTO 代码生成器通用实体类
type EntityDTO struct {
	EntityName string
	Comment    string
	FieldList  []FieldDTO
}

type FieldDTO struct {
	Name    string
	Type    string
	Comment string
}
