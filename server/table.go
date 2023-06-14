package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	. "sql_generate/consts"
	"sql_generate/core"
	"sql_generate/core/schema"
	"sql_generate/global"
	"sql_generate/models"
	"sql_generate/respository/db"
	"strings"
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

// AddTableInfo 保存表
func (s *TableService) AddTableInfo(ctx context.Context, t *models.TableInfo) (bool, error) {
	if t == nil {
		return false, fmt.Errorf("参数错误")
	}
	return db.AddTableInfo(ctx, t)
}

// GetMyTableInfoList 分页获取当前用户创建的资源列表
func (s *TableService) GetMyTableInfoList(ctx context.Context, req *models.TableInfoQueryRequest) ([]*models.TableInfo, error) {
	var tableInfo []*models.TableInfo
	db := global.DB.Where("userId = ?", req.UserID).Model(&models.TableInfo{})
	var err error
	db, err = s.GetTableQueryWrapper(db, req)
	if err != nil {
		return nil, err
	}
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&tableInfo); res.Error != nil {
		return nil, res.Error
	}
	return tableInfo, nil
}

// GetTableInfoById 根据 id 获取 Table
func (s *TableService) GetTableInfoById(ctx context.Context, id int64) (*models.TableInfo, error) {
	if id <= 0 {
		return nil, errors.New("请求参数错误")
	}
	resp, err := db.GetTableInfoByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("根据id查询Table错误：%v", err)
	}
	return resp, nil
}

// DeletedTableInfoByID 删除表记录
func (s *TableService) DeletedTableInfoByID(ctx context.Context, id int64) (bool, error) {
	return db.DeletedTableInfoByID(ctx, id)
}

// GetTableInfoList 分页获取列表
func (s *TableService) GetTableInfoList(ctx context.Context, req *models.TableInfoQueryRequest) ([]*models.TableInfo, error) {
	var tableInfo []*models.TableInfo
	db := global.DB.Model(&models.TableInfo{})
	var err error
	db, err = s.GetTableQueryWrapper(db, req)
	if err != nil {
		return nil, err
	}
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&tableInfo); res.Error != nil {
		return nil, res.Error
	}
	return tableInfo, nil
}

// GetTableQueryWrapper 获取查询包装类
func (s *TableService) GetTableQueryWrapper(db *gorm.DB, tableInfoQueryRequest *models.TableInfoQueryRequest) (*gorm.DB, error) {
	if tableInfoQueryRequest == nil {
		return nil, errors.New("请求参数为空")
	}
	tableInfoQuery := &models.TableInfo{}
	copier.Copy(tableInfoQuery, tableInfoQueryRequest)
	sortField := tableInfoQueryRequest.SortField
	sortOrder := tableInfoQueryRequest.SortOrder
	name := tableInfoQuery.Name
	content := tableInfoQuery.Content

	// name、content 需支持模糊搜索
	tableInfoQuery.Name = ""
	tableInfoQuery.Content = ""

	if name != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(name)+"%")
	}
	if content != "" {
		db = db.Where("content LIKE ?", "%"+strings.TrimSpace(content)+"%")
	}
	if sortField != "" {
		order := sortField
		if sortOrder == CommonConstToString[SORT_ORDER_ASC] {
			order += " ASC"
		} else {
			order += " DESC"
		}
		db = db.Order(order)
	}
	return db, nil
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
	if reviewStatus >= 0 && !GetReviewStatus(reviewStatus) {
		return fmt.Errorf("请求参数错误")
	}
	return nil
}
