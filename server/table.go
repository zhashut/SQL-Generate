package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sql_generate/consts"
	"sql_generate/core"
	"sql_generate/core/schema"
	"sql_generate/models"
	"sql_generate/respository/db"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/11
 * Time: 19:27
 * Description: No Description
 */

type TableService struct{}

func NewTableService() *TableService {
	return &TableService{}
}

func (s *TableService) AddTableInfo(ctx context.Context, t *models.TableInfo) (bool, error) {
	if t == nil {
		return false, fmt.Errorf("参数错误")
	}
	return db.AddTableInfo(ctx, t)
}

// ValidAndHandleTableInfo 校验并处理 add - 是否为创建校验
func (s *TableService) ValidAndHandleTableInfo(ctx context.Context, tableInfo *models.TableInfo, add bool) error {
	if tableInfo == nil {
		return fmt.Errorf("请求参数错误")
	}
	content := tableInfo.Content
	name := tableInfo.Name
	reviewStatus := tableInfo.ReviewStatus
	// 创建时，所有参数必须非空
	if add && (name == "" || content == "") {
		return fmt.Errorf("请求参数错误")
	}
	if name != "" && len(name) > 30 {
		return fmt.Errorf("名称过长")
	}
	if content != "" {
		if len(content) > 20000 {
			return fmt.Errorf("名称过长")
		}
		// 检验字段内容
		var tableSchema schema.TableSchema
		if err := json.Unmarshal([]byte(content), &tableSchema); err != nil {
			return err
		}
		if err := core.ValidSchema(&tableSchema); err != nil {
			return err
		}
	}
	if reviewStatus >= 0 && !consts.GetReviewStatus(reviewStatus) {
		return fmt.Errorf("请求参数错误")
	}
	return nil
}
