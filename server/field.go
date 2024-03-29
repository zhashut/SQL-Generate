package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	. "sql_generate/consts"
	"sql_generate/core"
	"sql_generate/core/builder"
	"sql_generate/core/schema"
	"sql_generate/models"
	"sql_generate/respository/cache"
	"sql_generate/respository/db"
	"strconv"
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
	Cache            *cache.Cache
	GenerateResolver GenerateResolver
	BuilderResolver  BuilderResolver
}

func NewFieldService() *FieldService {
	return &FieldService{
		GenerateResolver: core.NewGeneratorFace(),
		BuilderResolver:  builder.NewSQLBuilder(),
	}
}

// AddField 添加词条
func (s *FieldService) AddField(ctx context.Context, fieldAddReq *models.FieldInfoAddRequest, uid int64) (int64, error) {
	if fieldAddReq == nil {
		return 0, fmt.Errorf("field cannot be nil")
	}
	field := &models.FieldInfo{}
	_ = copier.Copy(field, fieldAddReq)
	// 检验
	if err := s.ValidAndHandleField(ctx, field, true); err != nil {
		return 0, err
	}
	field.UserId = uid
	result, err := s.DB.AddField(ctx, field)
	if !result || err != nil {
		return 0, fmt.Errorf("cannot add field: %v", err)
	}
	if err := s.Cache.DeleteKV(ctx, CACHE_FIELD_KEY+"*"); err != nil {
		zap.S().Errorf("failed to delete: %s", CACHE_FIELD_KEY)
		return field.ID, nil
	}
	return field.ID, nil
}

// GetFieldByID 根据id获取词条
func (s *FieldService) GetFieldByID(ctx context.Context, id int64) (*models.FieldInfo, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	var field *models.FieldInfo
	cacheKey := CACHE_DICT_KEY + strconv.FormatInt(id, 10)
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			field, err = s.DB.GetFieldByID(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("cannot get field: %v", err)
			}
			marshal, _ := json.Marshal(&field)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return field, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &field); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal response: %v", err)
	}
	return field, nil
}

// DeleteField 删除词条
func (s *FieldService) DeleteField(ctx context.Context, req *models.OnlyIDRequest, user *models.User) (bool, error) {
	if req == nil || req.ID <= 0 {
		return false, fmt.Errorf("incorrect request parameters: %v", req.ID)
	}
	// 判断是否存在
	field, err := s.DB.GetFieldByID(ctx, req.ID)
	if err != nil {
		return false, fmt.Errorf("cannot get field: %v", err)
	}
	// 仅本人和管理员可以删除
	if field.UserId != user.ID && user.UserRole != ADMIN {
		return false, fmt.Errorf("not access delete field")
	}
	b, err := s.DB.DeletedFieldByID(ctx, field.ID)
	if err != nil {
		return false, fmt.Errorf("cannot delete field: %v", err)
	}
	if err := s.Cache.DeleteKV(ctx, CACHE_FIELD_KEY+"*"); err != nil {
		zap.S().Errorf("failed to delete: %s", CACHE_FIELD_KEY)
		return b, nil
	}
	return b, nil
}

// GetMyAddFieldListPage 分页获取当前用户创建的资源列表
func (s *FieldService) GetMyAddFieldListPage(ctx context.Context, req *models.FieldInfoQueryRequest, user *models.User) ([]*models.FieldInfo, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	var fields []*models.FieldInfo
	cacheKey := CACHE_FIELD_KEY + ADD_LIST_PAGE + strconv.FormatInt(user.ID, 10)
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			req.UserID = user.ID
			fields, err = s.DB.GetMyAddFieldListPage(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("cannot get MyAddFieldListPage: %v", err)
			}
			marshal, _ := json.Marshal(&fields)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return fields, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &fields); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal MyAddFieldListPage: %v", err)
	}
	return fields, nil
}

// GetMyFieldListPage 分页获取当前用户可选的资源列表
func (s *FieldService) GetMyFieldListPage(ctx context.Context, req *models.FieldInfoQueryRequest, user *models.User) ([]*models.FieldInfo, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	var fields []*models.FieldInfo
	cacheKey := CACHE_FIELD_KEY + MY_LIST_PAGE + strconv.FormatInt(user.ID, 10)
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			req.UserID = user.ID
			req.ReviewStatus = ReviewStatusEnumToInt[PASS]
			fields, err = s.DB.GetMyFieldListPage(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("cannot get MyFieldListPage: %v", err)
			}
			marshal, _ := json.Marshal(&fields)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return fields, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &fields); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal MyFieldListPage: %v", err)
	}
	return fields, nil
}

// GetMyFieldList 获取当前用户可选的全部资源列表（只返回 id 和名称）
func (s *FieldService) GetMyFieldList(ctx context.Context, req *models.FieldInfoQueryRequest, user *models.User) ([]*models.FieldInfo, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	var fields []*models.FieldInfo
	cacheKey := CACHE_FIELD_KEY + MY_LIST + strconv.FormatInt(user.ID, 10)
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			fieldList := make([]*models.FieldInfo, 0)
			// 先查询审核通过的
			req.ReviewStatus = ReviewStatusEnumToInt[PASS]
			passFieldList, err := s.DB.GetMyFieldList(ctx, req, false)
			if err != nil {
				return nil, fmt.Errorf("cannot get pass field list: %v", err)
			}
			fieldList = append(fieldList, passFieldList...)
			// 查询当前登录用户的词条
			req.ReviewStatus = ReviewStatusEnumToInt[REVIEWING]
			req.UserID = user.ID
			myFieldList, err := s.DB.GetMyFieldList(ctx, req, true)
			if err != nil {
				return nil, fmt.Errorf("cannot get my field list: %v", err)
			}
			fieldList = append(fieldList, myFieldList...)
			// 根据id去重
			fields = fieldDeduplicate(fieldList)
			marshal, _ := json.Marshal(&fields)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return fields, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &fields); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal MyFieldList: %v", err)
	}
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
	var fields []*models.FieldInfo
	cacheKey := CACHE_FIELD_KEY + LIST_PAGE
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			fields, err = s.DB.GetFieldListPage(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("cannot get FieldListPage: %v", err)
			}
			marshal, _ := json.Marshal(&fields)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return fields, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &fields); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal FieldListPage: %v", err)
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
		if err := s.GenerateResolver.ValidField(&schemaField); err != nil {
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
