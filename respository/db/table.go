package db

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	. "sql_generate/consts"
	"sql_generate/global"
	"sql_generate/models"
	"strings"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/11
 * Time: 19:17
 * Description: No Description
 */

type TableDao struct{}

// AddTableInfo 添加表
func (dao *TableDao) AddTableInfo(ctx context.Context, d *models.TableInfo) (bool, error) {
	if res := global.DB.Where("name = ? and userId = ?", d.Name, d.UserId).First(&models.TableInfo{}); res.RowsAffected > 0 {
		return false, fmt.Errorf("表名称已存在")
	}
	global.DB.Save(&d)
	return true, nil
}

// GetTableInfoByID 根据 ID 获取表
func (dao *TableDao) GetTableInfoByID(ctx context.Context, id int64) (*models.TableInfo, error) {
	var table models.TableInfo
	if res := global.DB.Where("id = ?", id).First(&table); res.Error != nil {
		return nil, res.Error
	}
	return &table, nil
}

// DeletedTableInfoByID 删除表
func (dao *TableDao) DeletedTableInfoByID(ctx context.Context, id int64) (bool, error) {
	if res := global.DB.Delete(&models.TableInfo{ID: id}); res.RowsAffected == 0 {
		return false, fmt.Errorf("表不存在")
	}
	return true, nil
}

// GetMyAddTableInfoListPage 分页获取当前用户创建的资源列表
func (dao *TableDao) GetMyAddTableInfoListPage(ctx context.Context, req *models.TableInfoQueryRequest) ([]*models.TableInfo, error) {
	var tableList []*models.TableInfo
	db := global.DB.Where("userId = ?", req.UserID).Model(&models.TableInfo{})
	var err error
	db, err = dao.GetTableQueryWrapper(db, req)
	if err != nil {
		return nil, err
	}
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&tableList); res.Error != nil {
		return nil, res.Error
	}
	return tableList, nil
}

// GetMyTableInfoListPage 分页获取当前用户可选的资源列表
func (dao *TableDao) GetMyTableInfoListPage(ctx context.Context, req *models.TableInfoQueryRequest) ([]*models.TableInfo, error) {
	var tableList []*models.TableInfo
	db := global.DB.Where("userId = ?", req.UserID).Or("reviewStatus = ?", req.ReviewStatus).Model(&models.TableInfo{})
	var err error
	db, err = dao.GetTableQueryWrapper(db, req)
	if err != nil {
		return nil, err
	}
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&tableList); res.Error != nil {
		return nil, res.Error
	}
	return tableList, nil
}

// GetMyTableInfoList 获取当前用户可选的全部资源列表（只返回 id 和名称）
func (dao *TableDao) GetMyTableInfoList(ctx context.Context, req *models.TableInfoQueryRequest, isUser bool) ([]*models.TableInfo, error) {
	var tableList []*models.TableInfo
	// isUser 为 false 就是查询审核通过的，为true就是查询本人的
	if !isUser {
		db := global.DB.Where("reviewStatus = ?", ReviewStatusEnumToInt[PASS]).Model(&models.TableInfo{})
		var err error
		db, err = dao.GetTableQueryWrapper(db, req)
		if err != nil {
			return nil, err
		}
		if res := db.Select("id", "name").Find(&tableList); res.Error != nil {
			return nil, res.Error
		}
		return tableList, nil
	} else {
		db := global.DB.Where("userId = ?", req.UserID).Model(&models.TableInfo{})
		var err error
		if err != nil {
			return nil, err
		}
		if res := db.Select("id", "name").Find(&tableList); res.Error != nil {
			return nil, res.Error
		}
		return tableList, nil
	}
}

// GetTableInfoListPage 分页获取列表
func (dao *TableDao) GetTableInfoListPage(ctx context.Context, req *models.TableInfoQueryRequest) ([]*models.TableInfo, error) {
	var tableList []*models.TableInfo
	db := global.DB.Model(&models.TableInfo{})
	var err error
	db, err = dao.GetTableQueryWrapper(db, req)
	if err != nil {
		return nil, err
	}
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&tableList); res.Error != nil {
		return nil, res.Error
	}
	return tableList, nil
}

// GetTableQueryWrapper 获取查询包装类
func (dao *TableDao) GetTableQueryWrapper(db *gorm.DB, tableInfoQueryRequest *models.TableInfoQueryRequest) (*gorm.DB, error) {
	if tableInfoQueryRequest == nil {
		return nil, fmt.Errorf("请求参数为空")
	}
	tableInfoQuery := &models.TableInfo{}
	_ = copier.Copy(tableInfoQuery, tableInfoQueryRequest)
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
