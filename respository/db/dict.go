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
 * Date: 2023/6/14
 * Time: 8:38
 * Description: No Description
 */

// AddDict 添加词条
func AddDict(ctx context.Context, d *models.Dict) (bool, error) {
	if res := global.DB.Where("name = ? and userId = ?", d.Name, d.UserId).First(&models.Dict{}); res.RowsAffected > 0 {
		return false, fmt.Errorf("词条名称已存在")
	}
	global.DB.Save(&d)
	return true, nil
}

// GetDictByID 根据 ID 获取词条
func GetDictByID(ctx context.Context, id int64) (*models.Dict, error) {
	var dict models.Dict
	if res := global.DB.Where("id = ?", id).First(&dict); res.Error != nil {
		return nil, res.Error
	}
	return &dict, nil
}

// DeletedDictByID 删除词条
func DeletedDictByID(ctx context.Context, id int64) (bool, error) {
	if res := global.DB.Delete(&models.Dict{ID: id}); res.RowsAffected == 0 {
		return false, fmt.Errorf("词条不存在")
	}
	return true, nil
}

// GetMyAddDictListPage 分页获取当前用户创建的资源列表
func GetMyAddDictListPage(ctx context.Context, req *models.DictQueryRequest) ([]*models.Dict, error) {
	var dictList []*models.Dict
	db := global.DB.Where("userId = ?", req.UserID).Model(&models.Dict{})
	var err error
	db, err = GetDictQueryWrapper(db, req)
	if err != nil {
		return nil, err
	}
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&dictList); res.Error != nil {
		return nil, res.Error
	}
	return dictList, nil
}

// GetMyDictListPage 分页获取当前用户可选的资源列表
func GetMyDictListPage(ctx context.Context, req *models.DictQueryRequest) ([]*models.Dict, error) {
	var dictList []*models.Dict
	db := global.DB.Where("userId = ?", req.UserID).Or("reviewStatus = ?", req.ReviewStatus).Model(&models.Dict{})
	var err error
	db, err = GetDictQueryWrapper(db, req)
	if err != nil {
		return nil, err
	}
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&dictList); res.Error != nil {
		return nil, res.Error
	}
	return dictList, nil
}

// GetMyDictList 获取当前用户可选的全部资源列表（只返回 id 和名称）
func GetMyDictList(ctx context.Context, req *models.DictQueryRequest, isUser bool) ([]*models.Dict, error) {
	var dictList []*models.Dict
	// isUser 为 false 就是查询审核通过的，为true就是查询本人的
	if !isUser {
		db := global.DB.Where("reviewStatus = ?", ReviewStatusEnumToInt[PASS]).Model(&models.Dict{})
		var err error
		db, err = GetDictQueryWrapper(db, req)
		if err != nil {
			return nil, err
		}
		if res := db.Select("id", "name").Find(&dictList); res.Error != nil {
			return nil, res.Error
		}
		return dictList, nil
	} else {
		db := global.DB.Where("userId = ?", req.UserID).Model(&models.Dict{})
		var err error
		if err != nil {
			return nil, err
		}
		if res := db.Select("id", "name").Find(&dictList); res.Error != nil {
			return nil, res.Error
		}
		return dictList, nil
	}
}

// GetDictListPage 分页获取列表
func GetDictListPage(ctx context.Context, req *models.DictQueryRequest) ([]*models.Dict, error) {
	var dictList []*models.Dict
	db := global.DB.Model(&models.Dict{})
	var err error
	db, err = GetDictQueryWrapper(db, req)
	if err != nil {
		return nil, err
	}
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&dictList); res.Error != nil {
		return nil, res.Error
	}
	return dictList, nil
}

// GetDictQueryWrapper 获取查询包装类
func GetDictQueryWrapper(db *gorm.DB, dictQueryRequest *models.DictQueryRequest) (*gorm.DB, error) {
	if dictQueryRequest == nil {
		return nil, errors.New("请求参数为空")
	}
	dictQuery := &models.Dict{}
	_ = copier.Copy(dictQuery, dictQueryRequest)
	sortField := dictQueryRequest.SortField
	sortOrder := dictQueryRequest.SortOrder
	name := dictQuery.Name
	content := dictQuery.Content

	// name、content 需支持模糊搜索
	dictQuery.Name = ""
	dictQuery.Content = ""

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
