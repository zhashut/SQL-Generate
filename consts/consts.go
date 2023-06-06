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

// GetEnumByValue 根据 value 获取枚举
func GetEnumByValue(value string) MockTypeEnum {
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
