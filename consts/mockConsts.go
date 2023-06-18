package consts

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/5
 * Time: 23:13
 * Description: No Description
 */

type MockType int32

type MockTypeEnum string

const (
	NONE     MockTypeEnum = "不模拟"
	INCREASE MockTypeEnum = "递增"
	FIXED    MockTypeEnum = "固定"
	RANDOM   MockTypeEnum = "随机"
	RULE     MockTypeEnum = "规则"
	DICT     MockTypeEnum = "词库"
)

var MockTypeEnumToString = map[MockTypeEnum]string{
	NONE:     "不模拟",
	INCREASE: "递增",
	FIXED:    "固定",
	RANDOM:   "随机",
	RULE:     "规则",
	DICT:     "词库",
}

var MockTypeStringToEnum = map[string]MockTypeEnum{
	"不模拟": NONE,
	"递增":  INCREASE,
	"固定":  FIXED,
	"随机":  RANDOM,
	"规则":  RULE,
	"词库":  DICT,
}

// GetMockTypeEnumByValue 根据 value 获取枚举，默认返回 NONE
func GetMockTypeEnumByValue(value string) MockTypeEnum {
	for mockNum, mockString := range MockTypeEnumToString {
		if value == mockString {
			return mockNum
		}
	}
	return NONE
}

type FieldTypeEnum struct {
	Value    string
	JavaType string
	GoType   string
}

var (
	TINYINT    = FieldTypeEnum{"tinyint", "Integer", "int"}
	SMALLINT   = FieldTypeEnum{"smallint", "Integer", "int"}
	MEDIUMINT  = FieldTypeEnum{"mediumint", "Integer", "int"}
	INT        = FieldTypeEnum{"int", "Integer", "int"}
	BIGINT     = FieldTypeEnum{"bigint", "Long", "int64"}
	FLOAT      = FieldTypeEnum{"float", "Double", "float64 "}
	DOUBLE     = FieldTypeEnum{"double", "Double", "float64 "}
	DECIMAL    = FieldTypeEnum{"decimal", "BigDecimal", "float64"}
	DATE       = FieldTypeEnum{"date", "Date", "time.Time"}
	TIME       = FieldTypeEnum{"time", "Time", "time.Time"}
	YEAR       = FieldTypeEnum{"year", "Integer", "int"}
	DATETIME   = FieldTypeEnum{"datetime", "Date", "time.Time"}
	TIMESTAMP  = FieldTypeEnum{"timestamp", "Long", "int64"}
	CHAR       = FieldTypeEnum{"char", "String", "string"}
	VARCHAR    = FieldTypeEnum{"varchar", "String", "string"}
	TINYTEXT   = FieldTypeEnum{"tinytext", "String", "string"}
	TEXT       = FieldTypeEnum{"text", "String", "string"}
	MEDIUMTEXT = FieldTypeEnum{"mediumtext", "String", "string"}
	LONGTEXT   = FieldTypeEnum{"longtext", "String", "string"}
	TINYBLOB   = FieldTypeEnum{"tinyblob", "byte[]", "[]byte"}
	BLOB       = FieldTypeEnum{"blob", "byte[]", "[]byte"}
	MEDIUMBLOB = FieldTypeEnum{"mediumblob", "byte[]", "[]byte"}
	LONGBLOB   = FieldTypeEnum{"longblob", "byte[]", "[]byte"}
	BINARY     = FieldTypeEnum{"binary", "byte[]", "[]byte"}
	VARBINARY  = FieldTypeEnum{"varbinary", "byte[]", "[]byte"}
)

var FieldTypeEnumToString = map[FieldTypeEnum][]string{
	TINYINT:    {"tinyint", "Integer", "int"},
	SMALLINT:   {"smallint", "Integer", "int"},
	MEDIUMINT:  {"mediumint", "Integer", "int"},
	INT:        {"int", "Integer", "int"},
	BIGINT:     {"bigint", "Long", "int64"},
	FLOAT:      {"float", "Double", "float64 "},
	DOUBLE:     {"double", "Double", "float64 "},
	DECIMAL:    {"decimal", "BigDecimal", "float64"},
	DATE:       {"date", "Date", "time.Time"},
	TIME:       {"time", "Time", "time.Time"},
	YEAR:       {"year", "Integer", "int"},
	DATETIME:   {"datetime", "Date", "time.Time"},
	TIMESTAMP:  {"timestamp", "Long", "int64"},
	CHAR:       {"char", "String", "string"},
	VARCHAR:    {"varchar", "String", "string"},
	TINYTEXT:   {"tinytext", "String", "string"},
	TEXT:       {"text", "String", "string"},
	MEDIUMTEXT: {"mediumtext", "String", "string"},
	LONGTEXT:   {"longtext", "String", "string"},
	TINYBLOB:   {"tinyblob", "byte[]", "[]byte"},
	BLOB:       {"blob", "byte[]", "[]byte"},
	MEDIUMBLOB: {"mediumblob", "byte[]", "[]byte"},
	LONGBLOB:   {"longblob", "byte[]", "[]byte"},
	BINARY:     {"binary", "byte[]", "[]byte"},
	VARBINARY:  {"varbinary", "byte[]", "[]byte"},
}

var FieldTypeEnumStruct = map[FieldTypeEnum]struct{}{
	DATETIME:   {},
	TIMESTAMP:  {},
	DATE:       {},
	TIME:       {},
	CHAR:       {},
	VARCHAR:    {},
	TINYTEXT:   {},
	TEXT:       {},
	MEDIUMTEXT: {},
	TINYBLOB:   {},
	BLOB:       {},
	MEDIUMBLOB: {},
	LONGBLOB:   {},
	BINARY:     {},
	VARBINARY:  {},
}

type MockParamsRandomTypeEnum string

const (
	STRING     MockParamsRandomTypeEnum = "字符串"
	NAME       MockParamsRandomTypeEnum = "人名"
	CITY       MockParamsRandomTypeEnum = "城市"
	URL        MockParamsRandomTypeEnum = "网址"
	EMAIL      MockParamsRandomTypeEnum = "邮箱"
	IP         MockParamsRandomTypeEnum = "IP"
	INTEGER    MockParamsRandomTypeEnum = "整数"
	DECIMALS   MockParamsRandomTypeEnum = "小数"
	UNIVERSITY MockParamsRandomTypeEnum = "大学"
	DATES      MockParamsRandomTypeEnum = "日期"
	TIMESTAMPS MockParamsRandomTypeEnum = "时间戳"
	PHONE      MockParamsRandomTypeEnum = "手机号"
)

var MockParamsRandomTypeEnumToString = map[MockParamsRandomTypeEnum]string{
	STRING:     "字符串",
	NAME:       "人名",
	CITY:       "城市",
	URL:        "网址",
	EMAIL:      "邮箱",
	IP:         "IP",
	INTEGER:    "整数",
	DECIMALS:   "小数",
	UNIVERSITY: "大学",
	DATES:      "日期",
	TIMESTAMPS: "时间戳",
	PHONE:      "手机号",
}

var MockParamsRandomTypeStringToEnum = map[string]MockParamsRandomTypeEnum{
	"字符串": STRING,
	"人名":  NAME,
	"城市":  CITY,
	"网址":  URL,
	"邮箱":  EMAIL,
	"IP":  IP,
	"整数":  INTEGER,
	"小数":  DECIMALS,
	"大学":  UNIVERSITY,
	"日期":  DATES,
	"时间戳": TIMESTAMPS,
	"手机号": PHONE,
}

// GetMockParamsRandomTypeByValue 根据 value 获取枚举，默认返回 STRING
func GetMockParamsRandomTypeByValue(value string) MockParamsRandomTypeEnum {
	for mockNum, mockString := range MockParamsRandomTypeEnumToString {
		if value == mockString {
			return mockNum
		}
	}
	return STRING
}

// GetFieldTypeEnumByValue 根据 value 获取枚举，默认返回 TEXT
func GetFieldTypeEnumByValue(value string) FieldTypeEnum {
	for mockNum, mockString := range FieldTypeEnumToString {
		if value == mockString[SQLIndex] {
			return mockNum
		}
	}
	return TEXT
}
