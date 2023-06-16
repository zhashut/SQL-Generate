package db

import (
	"context"
	"errors"
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
 * Date: 2023/6/15
 * Time: 15:52
 * Description: 字段模块的数据处理
 */

type FieldDao struct{}

// AddField 添加字段
func (dao *FieldDao) AddField(ctx context.Context, f *models.FieldInfo) (bool, error) {
	if res := global.DB.Where("name = ? and userId = ?", f.Name, f.UserId).First(&models.FieldInfo{}); res.RowsAffected > 0 {
		return false, fmt.Errorf("字段名称已存在")
	}
	global.DB.Save(&f)
	return true, nil
}

// GetFieldByID 根据 ID 获取字段
func (dao *FieldDao) GetFieldByID(ctx context.Context, id int64) (*models.FieldInfo, error) {
	var field models.FieldInfo
	if res := global.DB.Where("id = ?", id).First(&field); res.Error != nil {
		return nil, res.Error
	}
	return &field, nil
}

// DeletedFieldByID 删除字段
func (dao *FieldDao) DeletedFieldByID(ctx context.Context, id int64) (bool, error) {
	if res := global.DB.Delete(&models.FieldInfo{ID: id}); res.RowsAffected == 0 {
		return false, fmt.Errorf("字段不存在")
	}
	return true, nil
}

// GetMyAddFieldListPage 分页获取当前用户创建的资源列表
func (dao *FieldDao) GetMyAddFieldListPage(ctx context.Context, req *models.FieldInfoQueryRequest) ([]*models.FieldInfo, error) {
	var fieldList []*models.FieldInfo
	db := global.DB.Where("userId = ?", req.UserID).Model(&models.FieldInfo{})
	var err error
	db, err = dao.GetFieldQueryWrapper(db, req)
	if err != nil {
		return nil, err
	}
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&fieldList); res.Error != nil {
		return nil, res.Error
	}
	return fieldList, nil
}

// GetMyFieldListPage 分页获取当前用户可选的资源列表
func (dao *FieldDao) GetMyFieldListPage(ctx context.Context, req *models.FieldInfoQueryRequest) ([]*models.FieldInfo, error) {
	var fieldList []*models.FieldInfo
	db := global.DB.Where("userId = ?", req.UserID).Or("reviewStatus = ?", req.ReviewStatus).Model(&models.FieldInfo{})
	var err error
	db, err = dao.GetFieldQueryWrapper(db, req)
	if err != nil {
		return nil, err
	}
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&fieldList); res.Error != nil {
		return nil, res.Error
	}
	return fieldList, nil
}

// GetMyFieldList 获取当前用户可选的全部资源列表（只返回 id 和名称）
func (dao *FieldDao) GetMyFieldList(ctx context.Context, req *models.FieldInfoQueryRequest, isUser bool) ([]*models.FieldInfo, error) {
	var fieldList []*models.FieldInfo
	// isUser 为 false 就是查询审核通过的，为true就是查询本人的
	if !isUser {
		db := global.DB.Where("reviewStatus = ?", ReviewStatusEnumToInt[PASS]).Model(&models.FieldInfo{})
		var err error
		db, err = dao.GetFieldQueryWrapper(db, req)
		if err != nil {
			return nil, err
		}
		if res := db.Select("id", "name").Find(&fieldList); res.Error != nil {
			return nil, res.Error
		}
		return fieldList, nil
	} else {
		db := global.DB.Where("userId = ?", req.UserID).Model(&models.FieldInfo{})
		var err error
		if err != nil {
			return nil, err
		}
		if res := db.Select("id", "name").Find(&fieldList); res.Error != nil {
			return nil, res.Error
		}
		return fieldList, nil
	}
}

// GetFieldListPage 分页获取列表
func (dao *FieldDao) GetFieldListPage(ctx context.Context, req *models.FieldInfoQueryRequest) ([]*models.FieldInfo, error) {
	var fieldList []*models.FieldInfo
	db := global.DB.Model(&models.FieldInfo{})
	var err error
	db, err = dao.GetFieldQueryWrapper(db, req)
	if err != nil {
		return nil, err
	}
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&fieldList); res.Error != nil {
		return nil, res.Error
	}
	return fieldList, nil
}

// dao.GetFieldQueryWrapper 获取查询包装类
func (dao *FieldDao) GetFieldQueryWrapper(db *gorm.DB, fieldQueryRequest *models.FieldInfoQueryRequest) (*gorm.DB, error) {
	if fieldQueryRequest == nil {
		return nil, errors.New("请求参数为空")
	}
	fieldQuery := &models.FieldInfo{}
	_ = copier.Copy(fieldQuery, fieldQueryRequest)
	searchName := fieldQueryRequest.SearchName
	sortField := fieldQueryRequest.SortField
	sortOrder := fieldQueryRequest.SortOrder
	name := fieldQuery.Name
	content := fieldQuery.Content
	fieldName := fieldQuery.FieldName

	// name、content 需支持模糊搜索
	fieldQuery.Name = ""
	fieldQuery.FieldName = ""
	fieldQuery.Content = ""

	if name != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(name)+"%")
	}
	if fieldName != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(fieldName)+"%")
	}
	if content != "" {
		db = db.Where("content LIKE ?", "%"+strings.TrimSpace(content)+"%")
	}
	// 同时按 name、fieldName搜索
	if searchName != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(searchName)+"%").
			Or("fieldName LIKE ?", "%"+strings.TrimSpace(searchName)+"%")
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
