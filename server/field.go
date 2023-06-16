package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	. "sql_generate/consts"
	"sql_generate/core"
	"sql_generate/core/builder"
	"sql_generate/core/schema"
	"sql_generate/global"
	"sql_generate/models"
	"sql_generate/respository/db"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/15
 * Time: 16:09
 * Description: 字段业务逻辑代码
 */

type FieldService struct {
	DB               *db.FieldDao
	UserResolver     UserResolver
	GenerateResolver GenerateResolver
	BuilderResolver  BuilderResolver
}

func NewFieldService() *FieldService {
	return &FieldService{
		UserResolver:     NewUserService(),
		GenerateResolver: core.NewGeneratorFace(),
		BuilderResolver:  builder.NewSQLBuilder(),
	}
}

// AddField 添加词条
func (s *FieldService) AddField(ctx context.Context, fieldAddReq *models.FieldInfoAddRequest) (int64, error) {
	if fieldAddReq == nil {
		return 0, fmt.Errorf("field cannot be nil")
	}
	field := &models.FieldInfo{}
	_ = copier.Copy(field, fieldAddReq)
	// 检验
	if err := s.ValidAndHandleField(ctx, field, true); err != nil {
		return 0, err
	}
	// 获取当前登录用户ID
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return 0, fmt.Errorf("cannot get login user: %v", err)
	}
	field.UserId = user.ID
	result, err := s.DB.AddField(ctx, field)
	if !result || err != nil {
		return 0, fmt.Errorf("cannot add field: %v", err)
	}
	return field.ID, nil
}

// GetFieldByID 根据id获取词条
func (s *FieldService) GetFieldByID(ctx context.Context, id int64) (*models.FieldInfo, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	field, err := s.DB.GetFieldByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get field: %v", err)
	}
	return field, nil
}

// DeleteField 删除词条
func (s *FieldService) DeleteField(ctx context.Context, req *models.OnlyIDRequest) (bool, error) {
	if req == nil || req.ID <= 0 {
		return false, fmt.Errorf("incorrect request parameters: %v", req.ID)
	}
	// 获取当前登录用户
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return false, fmt.Errorf("cannot get login user: %v", err)
	}
	// 判断是否存在
	field, err := s.DB.GetFieldByID(ctx, req.ID)
	if err != nil {
		return false, fmt.Errorf("cannot get field: %v", err)
	}
	// 仅本人和管理员可以删除
	admin, _ := s.UserResolver.IsAdmin(ctx, global.Session)
	if field.UserId != user.ID && !admin {
		return false, fmt.Errorf("not access delete field")
	}
	b, err := s.DB.DeletedFieldByID(ctx, field.ID)
	if err != nil {
		return false, fmt.Errorf("cannot delete field: %v", err)
	}
	return b, nil
}

// GetMyAddFieldListPage 分页获取当前用户创建的资源列表
func (s *FieldService) GetMyAddFieldListPage(ctx context.Context, req *models.FieldInfoQueryRequest) ([]*models.FieldInfo, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	// 获取当前登录用户
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return nil, fmt.Errorf("cannot get login user: %v", err)
	}
	req.UserID = user.ID
	fields, err := s.DB.GetMyAddFieldListPage(ctx, req)
	return fields, nil
}

// GetMyFieldListPage 分页获取当前用户可选的资源列表
func (s *FieldService) GetMyFieldListPage(ctx context.Context, req *models.FieldInfoQueryRequest) ([]*models.FieldInfo, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	// 获取当前登录用户
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return nil, fmt.Errorf("cannot get login user: %v", err)
	}
	req.UserID = user.ID
	req.ReviewStatus = ReviewStatusEnumToInt[PASS]
	fields, err := s.DB.GetMyFieldListPage(ctx, req)
	return fields, nil
}

// GetMyFieldList 获取当前用户可选的全部资源列表（只返回 id 和名称）
func (s *FieldService) GetMyFieldList(ctx context.Context, req *models.FieldInfoQueryRequest) ([]*models.FieldInfo, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	fieldList := make([]*models.FieldInfo, 0)
	// 先查询审核通过的
	req.ReviewStatus = ReviewStatusEnumToInt[PASS]
	passFieldList, err := s.DB.GetMyFieldList(ctx, req, false)
	if err != nil {
		return nil, fmt.Errorf("cannot get pass field list: %v", err)
	}
	fieldList = append(fieldList, passFieldList...)
	//查询本人的词条
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return nil, err
	}
	req.ReviewStatus = ReviewStatusEnumToInt[REVIEWING]
	req.UserID = user.ID
	myFieldList, err := s.DB.GetMyFieldList(ctx, req, true)
	if err != nil {
		return nil, fmt.Errorf("cannot get my field list: %v", err)
	}
	fieldList = append(fieldList, myFieldList...)
	// 根据id去重
	fields := fieldDeduplicate(fieldList)
	return fields, nil
}

// 根据id去重
func fieldDeduplicate(fieldList []*models.FieldInfo) []*models.FieldInfo {
	fieldMap := make(map[int64]*models.FieldInfo)

	for _, field := range fieldList {
		if _, ok := fieldMap[field.ID]; !ok {
			fieldMap[field.ID] = field
		}
	}

	deduplicated := make([]*models.FieldInfo, 0, len(fieldMap))
	for _, field := range fieldMap {
		deduplicated = append(deduplicated, field)
	}

	return deduplicated
}

// GetFieldListPage 分页获取列表
func (s *FieldService) GetFieldListPage(ctx context.Context, req *models.FieldInfoQueryRequest) ([]*models.FieldInfo, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	fields, err := s.DB.GetFieldListPage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot get FieldListPage: %v", err)
	}
	return fields, nil
}

// GenerateCreateSQL 生成创建表的 SQL
func (s *FieldService) GenerateCreateSQL(ctx context.Context, id int64) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("incorrect request parameters: %v", id)
	}
	field, err := s.DB.GetFieldByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("cannot get field: %v", err)
	}
	var schemaField schema.Field
	if err = json.Unmarshal([]byte(field.Content), &schemaField); err != nil {
		return "", fmt.Errorf("cannot unmarshal conetent: %v", err)
	}
	sql, err := s.BuilderResolver.BuildCreateFieldSQL(&schemaField)
	if err != nil {
		return "", fmt.Errorf("failed buildCreateFieldSQL: %v", err)
	}
	return sql, nil
}

// ValidAndHandleField 检验字段
func (s *FieldService) ValidAndHandleField(ctx context.Context, field *models.FieldInfo, add bool) error {
	if field == nil {
		return fmt.Errorf("请求参数错误")
	}
	content := field.Content
	name := field.Name
	reviewStatus := field.ReviewStatus
	// 创建时，所有参数必须非空
	if add && (name == "" || content == "") {
		return fmt.Errorf("请求参数错误")
	}
	if name != "" && len(name) > 30 {
		return fmt.Errorf("名称过长")
	}
	if content != "" {
		if len(content) > 20000 {
			return fmt.Errorf("内容过长")
		}
		// 检验字段内容
		var schemaField schema.Field
		if err := json.Unmarshal([]byte(content), &schemaField); err != nil {
			return fmt.Errorf("cannot unmarshal conetent: %v", err)
		}
		if err := s.GenerateResolver.ValidField(schemaField); err != nil {
			return fmt.Errorf("failed ValidField field: %v", err)
		}
		// 填充 fieldName
		field.FieldName = schemaField.FieldName
	}
	if reviewStatus >= 0 && !GetReviewStatus(reviewStatus) {
		return fmt.Errorf("请求参数错误")
	}
	return nil
}
