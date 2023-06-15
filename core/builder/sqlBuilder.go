package builder

import (
	"fmt"
	"go.uber.org/zap"
	. "sql_generate/consts"
	"sql_generate/core/schema"
	"strings"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/5
 * Time: 16:47
 * Description: SQL 生成器
 * 支持方言，策略模式
 */

const (
	DEFAULT           = "default "
	EMPTY             = " "
	NOT_NULL          = "not null"
	NULL              = "null"
	AUTO_INCREMENT    = "auto_increment"
	ON_UPDATE         = "on update "
	PRIMARY_KEY       = "primary key"
	CURRENT_TIMESTAMP = "CURRENT_TIMESTAMP"
)

type SQLBuilder struct {
	SQLDialect SQLDialect
}

// SQLDialect SQL 方言
type SQLDialect interface {
	WrapFieldName(name string) string       // 封装字段名
	ParseFieldName(fieldName string) string // 解析字段名
	WrapTableName(name string) string       // 封装表名
	ParseTableName(tableName string) string // 解析表名
}

// NewSQLBuilder TODO 这里可以设置成传参的方式来指定不同的实现，灵活变动
func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{
		SQLDialect: &Dialect{},
	}
}

// BuildCreateTableSql 构造建表SQL
func (s *SQLBuilder) BuildCreateTableSql(tableSchema *schema.TableSchema) (string, error) {
	// 构造表名
	tableName := s.SQLDialect.WrapTableName(tableSchema.TableName)
	dbName := tableSchema.DBName
	if dbName == "" {
		dbName = fmt.Sprintf("%s.%s", dbName, tableName)
	}
	// 构造表前缀注释
	tableComment := tableSchema.TableComment
	if tableComment == "" {
		tableComment = tableName
	}
	tablePrefixComment := fmt.Sprintf("-- %s", tableComment)
	// 构造表后缀注释
	tableSuffixComment := fmt.Sprintf("comment '%s'", tableComment)
	// 构造表字段
	fieldList := tableSchema.FieldList
	fieldSize := len(fieldList)
	fieldStrBuilder := strings.Builder{}
	for i := 0; i < fieldSize; i++ {
		field := fieldList[i]
		createSQLStr, err := s.BuildCreateFieldSQL(&field)
		if err != nil {
			return "", err
		}
		fieldAppend(&fieldStrBuilder, createSQLStr)
		// 最后一个字段后没有逗号和换行
		if i != fieldSize-1 {
			fieldAppend(&fieldStrBuilder, ",")
			fieldAppend(&fieldStrBuilder, "\n")
		}
	}
	fieldStr := fieldStrBuilder.String()
	// 构造模板
	result := fmt.Sprintf("%s\n"+
		"create table if not exists %s\n"+
		"(\n"+
		"%s\n"+
		") %s;", tablePrefixComment, tableName, fieldStr, tableSuffixComment)
	zap.S().Info("BuildCreateTableSql.result: ", result)
	return result, nil
}

// BuildCreateFieldSQL 生成创建字段的 SQL
func (s *SQLBuilder) BuildCreateFieldSQL(field *schema.Field) (string, error) {
	if field == nil {
		return "", fmt.Errorf("buildCreateFieldSQL: 请求参数错误")
	}
	fieldName := s.SQLDialect.WrapFieldName(field.FieldName)
	fieldType := field.FieldType
	defaultValue := field.DefaultValue
	notNil := field.NotNull
	comment := field.Comment
	onUpdate := field.OnUpdate
	primaryKey := field.PrimaryKey
	autoIncrement := field.AutoIncrement
	// <库名>.<表名> <column_name> int default 0 not null auto_increment comment '注释' primary Key
	fieldStrList := strings.Builder{}
	// 字段名
	fieldStrList.WriteString(fieldName)
	// 字段类型
	fieldAppend(&fieldStrList, EMPTY, fieldType)
	// 默认值
	if defaultValue == "" {
		fieldAppend(&fieldStrList, EMPTY, DEFAULT, getValueStr(field, defaultValue))
	}
	// 是否非空
	tmpValue := NOT_NULL
	if notNil {
		tmpValue = NULL
	}
	fieldAppend(&fieldStrList, EMPTY, tmpValue)
	// 是否自增
	if autoIncrement {
		fieldAppend(&fieldStrList, EMPTY, AUTO_INCREMENT)
	}
	// 附加条件
	if onUpdate != "" {
		fieldAppend(&fieldStrList, EMPTY, ON_UPDATE, onUpdate)
	}
	// 注释
	if comment != "" {
		fieldAppend(&fieldStrList, EMPTY, fmt.Sprintf("comment '%s'", comment))
	}
	// 是否为主键
	if primaryKey {
		fieldAppend(&fieldStrList, EMPTY, PRIMARY_KEY)
	}
	zap.S().Infof("buildCreateFieldSQL.fieldStrList: %s", fieldStrList.String())
	return fieldStrList.String(), nil
}

// BuildInsertSQL 构造插入数据 SQL
// e.g. INSERT INTO report (id, content) VALUES (1, "最好了")
func (s *SQLBuilder) BuildInsertSQL(tableSchema *schema.TableSchema, dataList []map[string]interface{}) (string, error) {
	// 构造表名
	tableName := s.SQLDialect.WrapTableName(tableSchema.TableName)
	dbName := tableSchema.DBName
	if dbName != "" {
		tableName = fmt.Sprintf("%s.%s", dbName, tableName)
	}
	// 构造表字段
	fieldList := tableSchema.FieldList
	// 过滤不模拟的字段
	tmpList := make([]schema.Field, 0)
	for _, field := range fieldList {
		typeEnum := MockTypeStringToEnum[field.MockType]
		if typeEnum != NONE {
			tmpList = append(tmpList, field)
		}
	}
	fieldList = tmpList
	// 拼接插入语句
	resultStringBuilder := strings.Builder{}
	total := len(dataList)
	for i := 0; i < total; i++ {
		dataRow := dataList[i]
		keyStr := s.getKeyStrWithJoin(fieldList)
		valueStr := s.getValueStrWithJoin(dataRow, fieldList)
		// 构造并填充模板
		result := fmt.Sprintf("insert into %s (%s) values (%s);", tableName, keyStr, valueStr)
		resultStringBuilder.WriteString(result)
		// 最后一个字段后没有换行
		if i != total-1 {
			resultStringBuilder.WriteString("\n")
		}
	}
	zap.S().Infof("BuildInsertSQL.resultStringBuilder: %s", resultStringBuilder.String())
	return resultStringBuilder.String(), nil
}

// 封装拼接字符串方法
func fieldAppend(fieldStrList *strings.Builder, fields ...string) {
	for _, field := range fields {
		fieldStrList.WriteString(field)
	}
}

// 获取字段键
func (s *SQLBuilder) getKeyStrWithJoin(fieldList []schema.Field) string {
	builder := strings.Builder{}
	for i, field := range fieldList {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(s.SQLDialect.WrapFieldName(field.FieldName))
	}
	return builder.String()
}

// 获取字段值
func (s *SQLBuilder) getValueStrWithJoin(dataRow map[string]interface{}, fieldList []schema.Field) string {
	builder := strings.Builder{}
	for i, field := range fieldList {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(getValueStr(&field, dataRow[field.FieldName]))
	}
	return builder.String()
}

// 根据列的属性获取值字符串
func getValueStr(field *schema.Field, value interface{}) string {
	if field == nil || value == nil {
		return "''"
	}

	fieldType := GetFieldTypeEnumByValue(field.FieldType)
	fieldValue := value.(string)
	if _, ok := FieldTypeEnumStruct[fieldType]; ok {
		return fmt.Sprintf("'%s'", fieldValue)
	} else {
		return fieldValue
	}
}
