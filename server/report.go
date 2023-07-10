package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	. "sql_generate/consts"
	"sql_generate/models"
	"sql_generate/respository/cache"
	"sql_generate/respository/db"
	"strconv"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/18
 * Time: 16:26
 * Description: 举报业务逻辑代码
 */

type ReportService struct {
	DB      *db.ReportDao
	Cache   *cache.Cache
	DictDao *db.DictDao
}

func NewReportService() *ReportService {
	return &ReportService{}
}

// AddReport 添加举报信息
func (s *ReportService) AddReport(ctx context.Context, req *models.ReportAddRequest, uid int64) (int64, error) {
	if req == nil {
		return 0, fmt.Errorf("report cannot be nil")
	}
	report := &models.Report{}
	_ = copier.Copy(report, req)
	// 检验
	if err := s.ValidReport(ctx, report, true); err != nil {
		return 0, err
	}
	// TODO 根据 reportedID 获取Dict，这里估计是不好控制动态获取当前页面，就直接用 词库页了
	dict, err := s.DictDao.GetDictByID(ctx, report.ReportedID)
	if err != nil {
		return 0, fmt.Errorf("举报对象不存在： %v", err)
	}

	report.ReportedUserId = dict.UserId
	report.UserId = uid
	report.Status = ReportStatusEnumToInt[DEFAULT]
	result, err := s.DB.SaveReport(ctx, report)
	if !result || err != nil {
		return 0, fmt.Errorf("cannot add report: %v", err)
	}
	if err := s.Cache.DeleteKV(ctx, CACHE_REPORT_KEY+"*"); err != nil {
		zap.S().Errorf("failed to delete: %s", CACHE_REPORT_KEY)
		return report.ID, nil
	}
	return report.ID, nil
}

// GetReportByID 根据id获取举报信息
func (s *ReportService) GetReportByID(ctx context.Context, id int64) (*models.Report, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	var report *models.Report
	cacheKey := CACHE_REPORT_KEY + strconv.FormatInt(id, 10)
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			report, err = s.DB.GetReportByID(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("cannot get report: %v", err)
			}
			marshal, _ := json.Marshal(&report)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return report, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &report); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal response: %v", err)
	}
	return report, nil
}

// DeleteReport 删除举报信息
func (s *ReportService) DeleteReport(ctx context.Context, req *models.OnlyIDRequest, user *models.User) (bool, error) {
	if req == nil || req.ID <= 0 {
		return false, fmt.Errorf("incorrect request parameters: %v", req.ID)
	}
	// 获取当前登录用户
	// 判断是否存在
	report, err := s.DB.GetReportByID(ctx, req.ID)
	if err != nil {
		return false, fmt.Errorf("cannot get report: %v", err)
	}
	// 仅本人和管理员可以删除
	if report.UserId != user.ID && user.UserRole != ADMIN {
		return false, fmt.Errorf("not access delete report")
	}
	b, err := s.DB.DeletedReportByID(ctx, report.ID)
	if err != nil {
		return false, fmt.Errorf("cannot delete report: %v", err)
	}
	if err := s.Cache.DeleteKV(ctx, CACHE_REPORT_KEY+"*"); err != nil {
		zap.S().Errorf("failed to delete: %s", CACHE_REPORT_KEY)
		return b, nil
	}
	return b, nil
}

// UpdateReport 更新（仅管理员）
func (s *ReportService) UpdateReport(ctx context.Context, req *models.ReportUpdateRequest, user *models.User) (bool, error) {
	if user.UserRole != ADMIN {
		return false, fmt.Errorf("权限不足")
	}

	if req == nil || req.ID <= 0 {
		return false, fmt.Errorf("report cannot be nil")
	}
	report := &models.Report{}
	_ = copier.Copy(report, req)
	// 检验
	if err := s.ValidReport(ctx, report, false); err != nil {
		return false, err
	}
	// 判断是否存在
	oldReport, err := s.DB.GetReportByID(ctx, req.ID)
	if err != nil {
		return false, fmt.Errorf("举报信息不存在: %v", oldReport)
	}

	result, err := s.DB.SaveReport(ctx, report)
	if !result || err != nil {
		return false, fmt.Errorf("cannot update report: %v", err)
	}
	if err := s.Cache.DeleteKV(ctx, CACHE_REPORT_KEY+"*"); err != nil {
		zap.S().Errorf("failed to delete: %s", CACHE_REPORT_KEY)
		return true, nil
	}
	return true, nil
}

// GetReportListPage 分页获取列表
func (s *ReportService) GetReportListPage(ctx context.Context, req *models.ReportQueryRequest) ([]*models.Report, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	var reports []*models.Report
	cacheKey := CACHE_REPORT_KEY + LIST_PAGE
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			reports, err = s.DB.GetReportListPage(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("cannot get ReportListPage: %v", err)
			}
			marshal, _ := json.Marshal(&reports)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return reports, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &reports); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal response: %v", err)
	}
	return reports, nil
}

// GetReportList 获取列表（仅管理员可使用）
func (s *ReportService) GetReportList(ctx context.Context, req *models.ReportQueryRequest, user *models.User) ([]*models.Report, error) {
	if user.UserRole != ADMIN {
		return nil, fmt.Errorf("权限不足")
	}
	var reports []*models.Report
	cacheKey := CACHE_REPORT_KEY + LIST
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			reports, err = s.DB.GetReportList(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("cannot get ReportList: %v", err)
			}
			marshal, _ := json.Marshal(&reports)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return reports, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &reports); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal response: %v", err)
	}
	return reports, nil
}

// ValidReport 检验
func (s *ReportService) ValidReport(ctx context.Context, report *models.Report, add bool) error {
	if report == nil {
		return fmt.Errorf("incorrect request parameters: %v", report)
	}
	content := report.Content
	reportedID := report.ReportedID
	status := report.Status
	if content != "" && len(content) > 1024 {
		return fmt.Errorf("内容过长")
	}
	// 创建时必须指定
	if add {
		if reportedID <= 0 {
			return fmt.Errorf("无效的 reportedID: %d", reportedID)
		}
	} else {
		if status >= 0 && !GetReportStatus(status) {
			return fmt.Errorf("请求参数错误")
		}
	}

	return nil
}
