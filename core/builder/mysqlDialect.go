package builder

import "fmt"

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/6
 * Time: 15:52
 * Description: mysql 方言
 */

type MySQLDialect struct{}

func NewMySQLDialect() SQLDialect {
	dialect, err := NewSQLDialectFactory().GetDialect(MSQL_DIALECT)
	if err != nil {
		panic(err)
	}
	return dialect
}

// WrapFieldName 封装字段名
func (d *MySQLDialect) WrapFieldName(name string) string {
	return fmt.Sprintf("`%s`", name)
}

// ParseFieldName 解析字段名
func (d *MySQLDialect) ParseFieldName(fieldName string) string {
	if fieldName[0] == '`' && fieldName[len(fieldName)-1] == '`' {
		fieldName = fieldName[1 : len(fieldName)-1]
		return fieldName
	}
	return fieldName
}

// WrapTableName 封装表名
func (d *MySQLDialect) WrapTableName(name string) string {
	return fmt.Sprintf("`%s`", name)
}

// ParseTableName 解析表名
func (d *MySQLDialect) ParseTableName(tableName string) string {
	if tableName[0] == '`' && tableName[len(tableName)-1] == '`' {
		tableName = tableName[1 : len(tableName)-1]
		return tableName
	}
	return tableName
}
