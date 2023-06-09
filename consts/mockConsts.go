package consts

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/5
 * Time: 23:13
 * Description: No Description
 */

type MockType int32

//const (
//	NONE MockType = iota
//	INCREASE
//	FIXED
//	RANDOM
//	RULE
//	DICT
//)

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

//func (m MockType) String() string {
//	return [...]string{"不模拟", "递增", "固定", "随机", "规则", "词库"}[m]
//}

type FieldTypeEnum string

const (
	TINYINT    FieldTypeEnum = "tinyint"
	SMALLINT   FieldTypeEnum = "smallint"
	MEDIUMINT  FieldTypeEnum = "mediumint"
	INT        FieldTypeEnum = "int"
	BIGINT     FieldTypeEnum = "bigint"
	FLOAT      FieldTypeEnum = "float"
	DOUBLE     FieldTypeEnum = "double"
	DECIMAL    FieldTypeEnum = "decimal"
	DATE       FieldTypeEnum = "date"
	TIME       FieldTypeEnum = "time"
	YEAR       FieldTypeEnum = "year"
	DATETIME   FieldTypeEnum = "datetime"
	TIMESTAMP  FieldTypeEnum = "timestamp"
	CHAR       FieldTypeEnum = "char"
	VARCHAR    FieldTypeEnum = "varchar"
	TINYTEXT   FieldTypeEnum = "tinytext"
	TEXT       FieldTypeEnum = "text"
	MEDIUMTEXT FieldTypeEnum = "mediumtext"
	LONGTEXT   FieldTypeEnum = "longtext"
	TINYBLOB   FieldTypeEnum = "tinyblob"
	BLOB       FieldTypeEnum = "blob"
	MEDIUMBLOB FieldTypeEnum = "mediumblob"
	LONGBLOB   FieldTypeEnum = "longblob"
	BINARY     FieldTypeEnum = "binary"
	VARBINARY  FieldTypeEnum = "varbinary"
)

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
