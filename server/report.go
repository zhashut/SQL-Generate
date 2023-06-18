package server

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	. "sql_generate/consts"
	"sql_generate/global"
	"sql_generate/models"
	"sql_generate/respository/db"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/18
 * Time: 16:26
 * Description: 举报业务逻辑代码
 */

type ReportService struct {
	DB           *db.ReportDao
	DictDao      *db.DictDao
	UserResolver UserResolver
}

func NewReportService() *ReportService {
	return &ReportService{
		UserResolver: NewUserService(),
	}
}

// AddReport 添加举报信息
func (s *ReportService) AddReport(ctx context.Context, req *models.ReportAddRequest) (int64, error) {
	if req == nil {
		return 0, fmt.Errorf("report cannot be nil")
	}
	report := &models.Report{}
	_ = copier.Copy(report, req)
	// 检验
	if err := s.ValidReport(ctx, report, true); err != nil {
		return 0, err
	}
	// 获取当前登录用户ID
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return 0, fmt.Errorf("cannot get login user: %v", err)
	}
	// TODO 根据 reportedID 获取Dict，这里估计是不好控制动态获取当前页面，就直接用 词库页了
	dict, err := s.DictDao.GetDictByID(ctx, report.ReportedID)
	if err != nil {
		return 0, fmt.Errorf("举报对象不存在： %v", err)
	}

	report.ReportedUserId = dict.UserId
	report.UserId = user.ID
	report.Status = ReportStatusEnumToInt[DEFAULT]
	result, err := s.DB.SaveReport(ctx, report)
	if !result || err != nil {
		return 0, fmt.Errorf("cannot add report: %v", err)
	}
	return report.ID, nil
}

// GetReportByID 根据id获取举报信息
func (s *ReportService) GetReportByID(ctx context.Context, id int64) (*models.Report, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	report, err := s.DB.GetReportByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get report: %v", err)
	}
	return report, nil
}

// DeleteReport 删除举报信息
func (s *ReportService) DeleteReport(ctx context.Context, req *models.OnlyIDRequest) (bool, error) {
	if req == nil || req.ID <= 0 {
		return false, fmt.Errorf("incorrect request parameters: %v", req.ID)
	}
	// 获取当前登录用户
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return false, fmt.Errorf("cannot get login user: %v", err)
	}
	// 判断是否存在
	report, err := s.DB.GetReportByID(ctx, req.ID)
	if err != nil {
		return false, fmt.Errorf("cannot get report: %v", err)
	}
	// 仅本人和管理员可以删除
	admin, _ := s.UserResolver.IsAdmin(ctx, global.Session)
	if report.UserId != user.ID && !admin {
		return false, fmt.Errorf("not access delete report")
	}
	b, err := s.DB.DeletedReportByID(ctx, report.ID)
	if err != nil {
		return false, fmt.Errorf("cannot delete report: %v", err)
	}
	return b, nil
}

// UpdateReport 更新（仅管理员）
func (s *ReportService) UpdateReport(ctx context.Context, req *models.ReportUpdateRequest) (bool, error) {
	// 获取当前登录用户ID
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil || user.UserRole != "admin" {
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
	return true, nil
}

// GetReportListPage 分页获取列表
func (s *ReportService) GetReportListPage(ctx context.Context, req *models.ReportQueryRequest) ([]*models.Report, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	fields, err := s.DB.GetReportListPage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot get ReportListPage: %v", err)
	}
	return fields, nil
}

// GetReportList 获取列表（仅管理员可使用）
func (s *ReportService) GetReportList(ctx context.Context, req *models.ReportQueryRequest) ([]*models.Report, error) {
	// 获取当前登录用户ID
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil || user.UserRole != "admin" {
		return nil, fmt.Errorf("权限不足")
	}
	fields, err := s.DB.GetReportList(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot get ReportListPage: %v", err)
	}
	return fields, nil
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
