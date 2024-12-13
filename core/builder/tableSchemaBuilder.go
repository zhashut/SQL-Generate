package builder

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"math"
	"mime/multipart"
	"regexp"
	. "sql_generate/consts"
	"sql_generate/core/schema"
	"sql_generate/global"
	"sql_generate/models"
	"sql_generate/utils"
	"strconv"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/16
 * Time: 18:45
 * Description: 表概要生成器
 */

const (
	FLASE_FIELD      = "false"
	VARCHAR_VAL      = "varchar"
	VARCHAR_VALS     = "varchar(256)"
	colPrimaryKey    = 1
	colNotPrimaryKey = 0
)

type TableSchemaBuilder struct {
	SQLDialect SQLDialect
}

func NewTableSchemaBuilder() *TableSchemaBuilder {
	return &TableSchemaBuilder{
		SQLDialect: NewMySQLDialect(),
	}
}

// BuildFromAuto 智能导入
func (b *TableSchemaBuilder) BuildFromAuto(content string) (*schema.TableSchema, error) {
	if len(content) == 0 || content == "" {
		return nil, fmt.Errorf("content is empty")
	}
	// 切分单词
	pattern := regexp.MustCompile("[,，]")
	words := pattern.Split(content, -1)
	if len(words) == 0 || len(words) > 20 {
		return nil, fmt.Errorf("content no compliance: %v", words)
	}
	// 根据单词去词库里匹配列信息，未匹配到的使用默认值
	var fieldInfoList []*models.FieldInfo
	if res := global.DB.Where("name in ?", words).Or("fieldName in ?", words).Find(&fieldInfoList); res.Error != nil {
		return nil, res.Error
	}
	// 名称 => 字段信息
	// 名称 => 字段信息
	nameFieldInfoMap := make(map[string][]*models.FieldInfo)
	fieldNameFieldInfoMap := make(map[string][]*models.FieldInfo)
	for _, fieldInfo := range fieldInfoList {
		nameFieldInfoMap[fieldInfo.Name] = append(nameFieldInfoMap[fieldInfo.Name], fieldInfo)
		fieldNameFieldInfoMap[fieldInfo.FieldName] = append(fieldNameFieldInfoMap[fieldInfo.FieldName], fieldInfo)
	}
	tableSchema := &schema.TableSchema{}
	tableSchema.TableName = "my_table"
	tableSchema.TableComment = "自动生成的表"
	fieldList := make([]*schema.Field, 0)
	for _, word := range words {
		var field *schema.Field
		var infoList []*models.FieldInfo
		infoList = append(infoList, nameFieldInfoMap[word]...)
		infoList = append(infoList, fieldNameFieldInfoMap[word]...)

		if len(infoList) > 0 {
			if err := json.Unmarshal([]byte(infoList[0].Content), &field); err != nil {
				return nil, fmt.Errorf("cannot unmarshal content: %v", err)
			}
		} else {
			field = getDefaultField(word)
		}

		fieldList = append(fieldList, field)
	}
	tableSchema.FieldList = fieldList
	return tableSchema, nil
}

// BuildFromSQL 根据建表 SQL 构建
func (b *TableSchemaBuilder) BuildFromSQL(sql string) (*schema.TableSchema, error) {
	//stmt, err := sqlparser.Parse(sql)
	//if err != nil {
	//	return nil, fmt.Errorf("SQL 解析错误: %v", err)
	//}
	//
	//createTable := stmt.(*sqlparser.CreateTable)
	//tableSchema := &schema.TableSchema{
	//	DBName:    createTable.Table.Qualifier.String(),
	//	TableName: createTable.Table.Name.String(),
	//	MockNum:   10,
	//}
	//
	//var fieldList []*schema.Field
	//for _, col := range createTable.TableSpec.Columns {
	//	field := &schema.Field{}
	//	field.FieldName = col.Name.String()
	//	if col.Type.Type == VARCHAR_VAL {
	//		field.FieldType = VARCHAR_VALS
	//	} else {
	//		field.FieldType = col.Type.Type
	//	}
	//	field.MockParams = ""
	//	defaultValue := ""
	//	if col.Type.Options.Default != nil {
	//		defaultValue = getExprVal(col.Type.Options.Default)
	//	}
	//	field.DefaultValue = defaultValue
	//	field.NotNull = !*col.Type.Options.Null
	//	if col.Type.Options.Comment != nil {
	//		field.Comment = col.Type.Options.Comment.Val
	//	}
	//	if col.Type.Options.KeyOpt == sqlparser.ColumnKeyOption(colPrimaryKey) {
	//		field.PrimaryKey = true
	//	}
	//	field.AutoIncrement = col.Type.Options.Autoincrement
	//	onUpdate := ""
	//	if col.Type.Options.OnUpdate != nil {
	//		onUpdate = getExprVal(col.Type.Options.OnUpdate)
	//	}
	//	field.OnUpdate = onUpdate
	//	field.MockType = MockTypeEnumToString[NONE]
	//	fieldList = append(fieldList, field)
	//}
	//tableSchema.FieldList = fieldList

	return nil, nil
}

// BuildFromExcel 导入 Excel
func (b *TableSchemaBuilder) BuildFromExcel(file multipart.File) (*schema.TableSchema, error) {
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("表格解析错误: %v", err)
	}
	// 获取表名集合
	sheets := xlsx.GetSheetMap()
	// 只支持表格的第一个工作簿
	rows, err := xlsx.Rows(sheets[1])
	if err != nil {
		return nil, fmt.Errorf("表格无数据: %v", err)
	}
	defer rows.Close()

	var fieldList []*schema.Field
	index := 0
	for rows.Next() {
		row, _ := rows.Columns()
		if index == 0 {
			for _, val := range row {
				field := &schema.Field{
					FieldName: val,
					Comment:   val,
					FieldType: FieldTypeEnumToString[TEXT][SQLIndex],
				}
				fieldList = append(fieldList, field)
			}
		} else if index == 1 {
			for i, val := range row {
				fieldType := getFieldTypeByValue(val).(string)
				fieldList[i].FieldType = fieldType
			}
		} else {
			break
		}
		index++
	}
	if index == 0 {
		return nil, fmt.Errorf("表格无数据")
	}
	tableSchema := &schema.TableSchema{
		FieldList: fieldList,
	}
	return tableSchema, nil
}

// GetFieldTypeByValue 判断字段类型
func getFieldTypeByValue(value string) interface{} {
	if len(value) == 0 {
		return FieldTypeEnumToString[TEXT][SQLIndex]
	}
	// 布尔
	if FLASE_FIELD == value {
		return FieldTypeEnumToString[TINYINT][SQLIndex]
	}
	// 整数
	if utils.IsNumeric(value) {
		number, _ := strconv.ParseInt(value, 10, 64)
		if number > math.MaxInt {
			return FieldTypeEnumToString[BIGINT][SQLIndex]
		}
		return FieldTypeEnumToString[INT][SQLIndex]
	}
	// 日期
	if utils.IsDate(value) {
		return FieldTypeEnumToString[DATETIME][SQLIndex]
	}
	// 小数
	if utils.IsDouble(value) {
		return FieldTypeEnumToString[DOUBLE][SQLIndex]
	}
	return FieldTypeEnumToString[TEXT][SQLIndex]
}

// 获取默认字段
func getDefaultField(word string) *schema.Field {
	field := &schema.Field{
		FieldName:     word,
		FieldType:     "text",
		DefaultValue:  "",
		NotNull:       false,
		Comment:       word,
		PrimaryKey:    false,
		AutoIncrement: false,
		MockType:      "",
		MockParams:    "",
		OnUpdate:      "",
	}
	return field
}

//func getExprVal(expr sqlparser.Expr) string {
//	exprVal := ""
//	if curTimeFuncExpr, ok := expr.(*sqlparser.CurTimeFuncExpr); ok {
//		exprVal = strings.ToUpper(curTimeFuncExpr.Name.String())
//	} else {
//		exprVal = sqlparser.String(expr)	
//	}
//	return exprVal
//}
