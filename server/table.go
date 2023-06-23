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
	"sql_generate/global"
	"sql_generate/models"
	"sql_generate/respository/cache"
	"sql_generate/respository/db"
	"strconv"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/11
 * Time: 19:27
 * Description: No Description
 */

type TableService struct {
	DB               *db.TableDao
	Cache            *cache.Cache
	UserResolver     UserResolver
	GenerateResolver GenerateResolver
	BuilderResolver  BuilderResolver
}

func NewTableService() *TableService {
	return &TableService{
		UserResolver:     NewUserService(),
		GenerateResolver: core.NewGeneratorFace(),
		BuilderResolver:  builder.NewSQLBuilder(),
	}
}

// AddTableInfo 添加表
func (s *TableService) AddTableInfo(ctx context.Context, tableAddReq *models.TableInfoAddRequest) (int64, error) {
	table := &models.TableInfo{}
	_ = copier.Copy(table, tableAddReq)
	// 检验
	if err := s.ValidAndHandleTableInfo(ctx, table, true); err != nil {
		return 0, err
	}
	// 获取当前登录用户ID
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return 0, fmt.Errorf("cannot get login user: %v", err)
	}
	table.UserId = user.ID
	result, err := s.DB.AddTableInfo(ctx, table)
	if !result || err != nil {
		return 0, fmt.Errorf("cannot add table: %v", err)
	}
	if err := s.Cache.DeleteKV(ctx, CACHE_TABLE_KEY+"*"); err != nil {
		zap.S().Errorf("failed to delete: %s", CACHE_TABLE_KEY)
		return table.ID, nil
	}
	return table.ID, nil
}

// GetTableInfoByID 根据id获取表
func (s *TableService) GetTableInfoByID(ctx context.Context, id int64) (*models.TableInfo, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	var table *models.TableInfo
	cacheKey := CACHE_TABLE_KEY + strconv.FormatInt(id, 10)
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			table, err = s.DB.GetTableInfoByID(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("cannot get table: %v", err)
			}
			marshal, _ := json.Marshal(&table)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return table, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &table); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal response: %v", err)
	}
	return table, nil
}

// DeleteTableInfo 删除表
func (s *TableService) DeleteTableInfo(ctx context.Context, req *models.OnlyIDRequest) (bool, error) {
	if req == nil || req.ID <= 0 {
		return false, fmt.Errorf("incorrect request parameters: %v", req.ID)
	}
	// 获取当前登录用户
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return false, fmt.Errorf("cannot get login user: %v", err)
	}
	// 判断是否存在
	table, err := s.DB.GetTableInfoByID(ctx, req.ID)
	if err != nil {
		return false, fmt.Errorf("cannot get table: %v", err)
	}
	// 仅本人和管理员可以删除
	admin, _ := s.UserResolver.IsAdmin(ctx, global.Session)
	if table.UserId != user.ID && !admin {
		return false, fmt.Errorf("not access delete table")
	}
	b, err := s.DB.DeletedTableInfoByID(ctx, table.ID)
	if err != nil {
		return false, fmt.Errorf("cannot delete table: %v", err)
	}
	if err := s.Cache.DeleteKV(ctx, CACHE_TABLE_KEY+"*"); err != nil {
		zap.S().Errorf("failed to delete: %s", CACHE_TABLE_KEY)
		return b, nil
	}
	return b, nil
}

// GetMyAddTableInfoListPage 分页获取当前用户创建的资源列表
func (s *TableService) GetMyAddTableInfoListPage(ctx context.Context, req *models.TableInfoQueryRequest) ([]*models.TableInfo, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	// 获取当前登录用户
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return nil, fmt.Errorf("cannot get login user: %v", err)
	}
	var tables []*models.TableInfo
	cacheKey := CACHE_TABLE_KEY + ADD_LIST_PAGE + strconv.FormatInt(user.ID, 10)
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			req.UserID = user.ID
			tables, err = s.DB.GetMyAddTableInfoListPage(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("cannot get TableInfoListPage: %v", err)
			}
			marshal, _ := json.Marshal(&tables)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return tables, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &tables); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal TableInfoListPage: %v", err)
	}
	return tables, nil
}

// GetMyTableInfoListPage 分页获取当前用户可选的资源列表
func (s *TableService) GetMyTableInfoListPage(ctx context.Context, req *models.TableInfoQueryRequest) ([]*models.TableInfo, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	// 获取当前登录用户
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return nil, fmt.Errorf("cannot get login user: %v", err)
	}
	var tables []*models.TableInfo
	cacheKey := CACHE_TABLE_KEY + MY_LIST_PAGE + strconv.FormatInt(user.ID, 10)
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			req.UserID = user.ID
			req.ReviewStatus = ReviewStatusEnumToInt[PASS]
			tables, err = s.DB.GetMyTableInfoListPage(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("cannot get MyTableInfoListPage: %v", err)
			}
			marshal, _ := json.Marshal(&tables)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return tables, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &tables); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal MyTableInfoListPage: %v", err)
	}
	return tables, nil
}

// GetMyTableInfoList 获取当前用户可选的全部资源列表（只返回 id 和名称）
func (s *TableService) GetMyTableInfoList(ctx context.Context, req *models.TableInfoQueryRequest) ([]*models.TableInfo, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return nil, err
	}
	var tables []*models.TableInfo
	cacheKey := CACHE_TABLE_KEY + MY_LIST + strconv.FormatInt(user.ID, 10)
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			tableList := make([]*models.TableInfo, 0)
			// 先查询审核通过的
			req.ReviewStatus = ReviewStatusEnumToInt[PASS]
			passTableInfoList, err := s.DB.GetMyTableInfoList(ctx, req, false)
			if err != nil {
				return nil, fmt.Errorf("cannot get pass table list: %v", err)
			}
			tableList = append(tableList, passTableInfoList...)
			//查询本人的表
			req.ReviewStatus = ReviewStatusEnumToInt[REVIEWING]
			req.UserID = user.ID
			myTableInfoList, err := s.DB.GetMyTableInfoList(ctx, req, true)
			if err != nil {
				return nil, fmt.Errorf("cannot get my table list: %v", err)
			}
			tableList = append(tableList, myTableInfoList...)
			// 根据id去重
			tables = tableDeduplicate(tableList)

			marshal, _ := json.Marshal(&tables)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return tables, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &tables); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal MyTableInfoList: %v", err)
	}
	return tables, nil
}

// 根据id去重
func tableDeduplicate(tableList []*models.TableInfo) []*models.TableInfo {
	tableMap := make(map[int64]*models.TableInfo)

	for _, table := range tableList {
		if _, ok := tableMap[table.ID]; !ok {
			tableMap[table.ID] = table
		}
	}

	deduplicated := make([]*models.TableInfo, 0, len(tableMap))
	for _, table := range tableMap {
		deduplicated = append(deduplicated, table)
	}

	return deduplicated
}

// GetTableInfoListPage 分页获取列表
func (s *TableService) GetTableInfoListPage(ctx context.Context, req *models.TableInfoQueryRequest) ([]*models.TableInfo, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	var tables []*models.TableInfo
	cacheKey := CACHE_TABLE_KEY + LIST_PAGE
	value, err := s.Cache.GetKV(ctx, cacheKey)
	if err != nil {
		if err == redis.Nil {
			tables, err = s.DB.GetTableInfoListPage(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("cannot get TableInfoListPage: %v", err)
			}
			marshal, _ := json.Marshal(&tables)
			if err = s.Cache.SetKV(ctx, cacheKey, marshal, 0); err != nil {
				return nil, fmt.Errorf("cannot set KV: %v", err)
			}
			return tables, nil
		}
		return nil, fmt.Errorf("unknow failed error: %v", err)
	}
	if err := json.Unmarshal([]byte(value), &tables); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal TableInfoListPage: %v", err)
	}
	return tables, nil
}

// GenerateCreateSQL 生成创建表的 SQL
func (s *TableService) GenerateCreateSQL(ctx context.Context, id int64) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("incorrect request parameters: %v", id)
	}
	table, err := s.DB.GetTableInfoByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("cannot get table: %v", err)
	}
	var tableSchema schema.TableSchema
	if err = json.Unmarshal([]byte(table.Content), &tableSchema); err != nil {
		return "", fmt.Errorf("cannot unmarshal conetent: %v", err)
	}
	sql, err := s.BuilderResolver.BuildCreateTableSql(&tableSchema)
	if err != nil {
		return "", fmt.Errorf("failed buildCreateTableInfoSQL: %v", err)
	}
	return sql, nil
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
		if err := s.GenerateResolver.ValidSchema(&tableSchema); err != nil {
			return err
		}
	}
	if reviewStatus >= 0 && !GetReviewStatus(reviewStatus) {
		return fmt.Errorf("请求参数错误")
	}
	return nil
}
