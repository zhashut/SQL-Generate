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
 * Description:  * SQL 生成器
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

type SQLDialect interface {
	WrapFieldName(name string) string       // 封装字段名
	ParseFieldName(fieldName string) string // 解析字段名
	WrapTableName(name string) string       // 封装表名
	ParseTableName(tableName string) string // 解析表名
}

func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{}
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
		createSQLStr, err := s.buildCreateFieldSQL(&field)
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
	// 填充模板
	// 构造模板
	result := fmt.Sprintf("%s\n"+
		"create table if not exists %s\n"+
		"(\n"+
		"%s\n"+
		") %s;", tablePrefixComment, tableName, fieldStr, tableSuffixComment)
	zap.S().Info("sql result: ", result)
	return result, nil
}

// 生成创建字段的 SQL
func (s *SQLBuilder) buildCreateFieldSQL(field *schema.Field) (string, error) {
	if field == nil {
		return "", fmt.Errorf("请求参数错误")
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
	fieldStrSlice := strings.Builder{}
	// 字段名
	fieldStrSlice.WriteString(fieldName)
	// 字段类型
	fieldAppend(&fieldStrSlice, EMPTY, fieldType)
	// 默认值
	if defaultValue == "" {
		// TODO getValueStr 未实现
		fieldAppend(&fieldStrSlice, EMPTY, DEFAULT, getValueStr(field, defaultValue))
	}
	// 是否非空
	tmpValue := NOT_NULL
	if notNil {
		tmpValue = NULL
	}
	fieldAppend(&fieldStrSlice, EMPTY, tmpValue)
	// 是否自增
	if autoIncrement {
		fieldAppend(&fieldStrSlice, EMPTY, AUTO_INCREMENT)
	}
	// 附加条件
	if onUpdate != "" {
		fieldAppend(&fieldStrSlice, EMPTY, ON_UPDATE, onUpdate)
	}
	// 注释
	if comment != "" {
		fieldAppend(&fieldStrSlice, EMPTY, fmt.Sprintf("comment '%s'", comment))
	}
	// 是否为主键
	if primaryKey {
		fieldAppend(&fieldStrSlice, EMPTY, PRIMARY_KEY)
	}
	return fieldStrSlice.String(), nil
}

// BuildInsertSQL TODO 构造插入数据 SQL
// e.g. INSERT INTO report (id, content) VALUES (1, "瑶瑶最好了")
func (s *SQLBuilder) BuildInsertSQL(tableSchema *schema.TableSchema, dataList interface{}) (string, error) {
	return "", nil
}

func fieldAppend(fieldStrSlice *strings.Builder, fields ...string) {
	for _, field := range fields {
		fieldStrSlice.WriteString(field)
	}
}

// 根据列的属性获取值字符串
func getValueStr(field *schema.Field, value interface{}) string {
	if field == nil || value == nil {
		return "''"
	}

	fieldType := FieldTypeEnum(field.FieldType)
	fieldValue := value.(string)
	switch fieldType {
	case DATETIME:
	case TIMESTAMP:
		if fieldValue != CURRENT_TIMESTAMP {
			return fieldValue
		} else {
			return fmt.Sprintf("'%s'", fieldValue)
		}
	case DATE:
	case TIME:
	case CHAR:
	case VARCHAR:
	case TINYTEXT:
	case TEXT:
	case MEDIUMTEXT:
	case LONGTEXT:
	case TINYBLOB:
	case BLOB:
	case MEDIUMBLOB:
	case LONGBLOB:
	case BINARY:
	case VARBINARY:
		return fmt.Sprintf("'%s'", fieldValue)
	default:
		return fieldValue
	}
	return fieldValue
}
