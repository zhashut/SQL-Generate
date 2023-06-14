package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/jinzhu/copier"
	"regexp"
	. "sql_generate/consts"
	"sql_generate/core"
	"sql_generate/core/schema"
	"sql_generate/global"
	"sql_generate/models"
	"sql_generate/respository/db"
	"strconv"
	"strings"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/14
 * Time: 7:53
 * Description: 词库业务逻辑代码
 */

type DictService struct {
	UserResolver     UserResolver
	GenerateResolver GenerateResolver
}

type UserResolver interface {
	GetLoginUser(ctx context.Context, session sessions.Session) (*models.User, error)
	IsAdmin(ctx context.Context, session sessions.Session) (bool, error)
}

type GenerateResolver interface {
	GenerateAll(tableSchema *schema.TableSchema) (*models.Generate, error)
}

func NewDictService() *DictService {
	return &DictService{
		UserResolver:     NewUserService(),
		GenerateResolver: core.NewGeneratorFace(),
	}
}

// AddDict 添加词条
func (s *DictService) AddDict(ctx context.Context, dictAddReq *models.DictAddRequest) (int64, error) {
	dict := &models.Dict{}
	_ = copier.Copy(dict, dictAddReq)
	// 检验
	if err := s.ValidAndHandleDict(ctx, dict, true); err != nil {
		return 0, err
	}
	// 获取当前登录用户ID
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return 0, fmt.Errorf("cannot get login user: %v", err)
	}
	dict.UserId = user.ID
	result, err := db.AddDict(ctx, dict)
	if !result || err != nil {
		return 0, fmt.Errorf("cannot add dict: %v", err)
	}
	return dict.ID, nil
}

// GetDictByID 根据id获取词条
func (s *DictService) GetDictByID(ctx context.Context, id int64) (*models.Dict, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	dict, err := db.GetDictByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get dict: %v", err)
	}
	return dict, nil
}

// DeleteDict 删除词条
func (s *DictService) DeleteDict(ctx context.Context, req *models.OnlyIDRequest) (bool, error) {
	if req == nil || req.ID <= 0 {
		return false, fmt.Errorf("incorrect request parameters: %v", req.ID)
	}
	// 获取当前登录用户
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return false, fmt.Errorf("cannot get login user: %v", err)
	}
	// 判断是否存在
	dict, err := db.GetDictByID(ctx, req.ID)
	if err != nil {
		return false, fmt.Errorf("cannot get dict: %v", err)
	}
	// 仅本人和管理员可以删除
	admin, _ := s.UserResolver.IsAdmin(ctx, global.Session)
	if dict.UserId != user.ID && !admin {
		return false, fmt.Errorf("not access delete dict")
	}
	b, err := db.DeletedDictByID(ctx, dict.ID)
	if err != nil {
		return false, fmt.Errorf("cannot delete dict: %v", err)
	}
	return b, nil
}

// GetMyAddDictListPage 分页获取当前用户创建的资源列表
func (s *DictService) GetMyAddDictListPage(ctx context.Context, req *models.DictQueryRequest) ([]*models.Dict, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	// 获取当前登录用户
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return nil, fmt.Errorf("cannot get login user: %v", err)
	}
	req.UserID = user.ID
	dicts, err := db.GetMyAddDictListPage(ctx, req)
	return dicts, nil
}

// GetMyDictListPage 分页获取当前用户可选的资源列表
func (s *DictService) GetMyDictListPage(ctx context.Context, req *models.DictQueryRequest) ([]*models.Dict, error) {
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
	dicts, err := db.GetMyDictListPage(ctx, req)
	return dicts, nil
}

// GetMyDictList 获取当前用户可选的全部资源列表（只返回 id 和名称）
func (s *DictService) GetMyDictList(ctx context.Context, req *models.DictQueryRequest) ([]*models.Dict, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	dictList := make([]*models.Dict, 0)
	// 先查询审核通过的
	req.ReviewStatus = ReviewStatusEnumToInt[PASS]
	passDictList, err := db.GetMyDictList(ctx, req, false)
	if err != nil {
		return nil, fmt.Errorf("cannot get pass dict list: %v", err)
	}
	dictList = append(dictList, passDictList...)
	//查询本人的词条
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return nil, err
	}
	req.ReviewStatus = ReviewStatusEnumToInt[REVIEWING]
	req.UserID = user.ID
	myDictList, err := db.GetMyDictList(ctx, req, true)
	if err != nil {
		return nil, fmt.Errorf("cannot get my dict list: %v", err)
	}
	dictList = append(dictList, myDictList...)
	// 根据id去重
	dicts := deduplicate(dictList)
	return dicts, nil
}

// 根据id去重
func deduplicate(dictList []*models.Dict) []*models.Dict {
	dictMap := make(map[int64]*models.Dict)

	for _, dict := range dictList {
		if _, ok := dictMap[dict.ID]; !ok {
			dictMap[dict.ID] = dict
		}
	}

	deduplicated := make([]*models.Dict, 0, len(dictMap))
	for _, dict := range dictMap {
		deduplicated = append(deduplicated, dict)
	}

	return deduplicated
}

// GetDictListPage 分页获取列表
func (s *DictService) GetDictListPage(ctx context.Context, req *models.DictQueryRequest) ([]*models.Dict, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request parameters: %v", req)
	}
	// 获取当前登录用户
	user, err := s.UserResolver.GetLoginUser(ctx, global.Session)
	if err != nil {
		return nil, fmt.Errorf("cannot get login user: %v", err)
	}
	req.UserID = user.ID
	dicts, err := db.GetDictListPage(ctx, req)
	return dicts, nil
}

// GenerateCreateSQL 生成创建表的 SQL
func (s *DictService) GenerateCreateSQL(ctx context.Context, id int64) (*models.Generate, error) {
	if id <= 0 {
		return nil, fmt.Errorf("incorrect request parameters: %v", id)
	}
	dict, err := db.GetDictByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get dict: %v", err)
	}
	// 根据词库生成 Schema
	tableSchema := &schema.TableSchema{}
	name := dict.Name
	tableSchema.TableName = "dict"
	tableSchema.TableComment = name
	fieldList := make([]schema.Field, 0)
	var idField schema.Field
	idField.FieldName = "id"
	idField.FieldType = "bigint"
	idField.NotNull = true
	idField.Comment = "id"
	idField.PrimaryKey = true
	idField.AutoIncrement = true
	var dataField schema.Field
	dataField.FieldName = "data"
	dataField.FieldType = "text"
	dataField.Comment = "数据"
	dataField.MockType = MockTypeEnumToString[DICT]
	dataField.MockParams = strconv.Itoa(int(id))
	fieldList = append(fieldList, idField, dataField)
	tableSchema.FieldList = fieldList
	all, err := s.GenerateResolver.GenerateAll(tableSchema)
	if err != nil {
		return nil, fmt.Errorf("failed generate all: %v", err)
	}
	return all, nil
}

// ValidAndHandleDict 检验词条
func (s *DictService) ValidAndHandleDict(ctx context.Context, dict *models.Dict, add bool) error {
	if dict == nil {
		return fmt.Errorf("请求参数错误")
	}
	content := dict.Content
	name := dict.Name
	reviewStatus := dict.ReviewStatus
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
		// 对 content 进行转换
		// 对中文和英文逗号进行分割
		pattern := regexp.MustCompile("[,，]")
		words := pattern.Split(content, -1)
		// 移除开头和结尾空格
		for i, word := range words {
			words[i] = strings.TrimSpace(word)
		}
		wordsList := make([]string, 0, len(words))
		for _, word := range words {
			if word != "" {
				wordsList = append(wordsList, word)
			}
		}
		content, err := json.Marshal(wordsList)
		if err != nil {
			return fmt.Errorf("cannot marshal wordsList: %v", err)
		}
		dict.Content = string(content)
		if reviewStatus >= 0 && !GetReviewStatus(reviewStatus) {
			return fmt.Errorf("请求参数错误")
		}
	}
	return nil
}
