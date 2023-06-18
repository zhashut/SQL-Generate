package db

import (
	"context"
	"fmt"
	. "sql_generate/consts"
	"sql_generate/global"
	"sql_generate/models"
	"strings"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/18
 * Time: 16:17
 * Description: 举报模块的数据处理
 */

type ReportDao struct{}

// SaveReport 添加/更改举报信息
func (dao *ReportDao) SaveReport(ctx context.Context, r *models.Report) (bool, error) {
	if res := global.DB.Save(&r); res.Error != nil {
		return false, res.Error
	}
	return true, nil
}

// GetReportByID 根据 ID 获取举报信息
func (dao *ReportDao) GetReportByID(ctx context.Context, id int64) (*models.Report, error) {
	var report models.Report
	if res := global.DB.Where("id = ?", id).First(&report); res.Error != nil {
		return nil, res.Error
	}
	return &report, nil
}

// DeletedReportByID 删除举报信息
func (dao *ReportDao) DeletedReportByID(ctx context.Context, id int64) (bool, error) {
	if res := global.DB.Delete(&models.Report{ID: id}); res.RowsAffected == 0 {
		return false, fmt.Errorf("举报信息不存在")
	}
	return true, nil
}

// GetReportListPage 分页获取列表
func (dao *ReportDao) GetReportListPage(ctx context.Context, req *models.ReportQueryRequest) ([]*models.Report, error) {
	var reportList []*models.Report
	db := global.DB.Model(&models.Report{})
	sortField := req.SortField
	sortOrder := req.SortOrder
	content := req.Content

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
	if res := db.Scopes(Paginate(int(req.Pages), int(req.PageSize))).Find(&reportList); res.Error != nil {
		return nil, res.Error
	}
	return reportList, nil
}

// GetReportList 获取列表（仅管理员可使用）
func (dao *ReportDao) GetReportList(ctx context.Context, req *models.ReportQueryRequest) ([]*models.Report, error) {
	var reportList []*models.Report
	db := global.DB.Model(&models.Report{})
	if res := db.Find(&reportList); res.Error != nil {
		return nil, res.Error
	}
	return reportList, nil
}
