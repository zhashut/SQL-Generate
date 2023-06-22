package builder

import (
	"fmt"
	"sync"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/22
 * Time: 20:24
 * Description:   SQL 方言工厂
 *				工厂 + 单例模式，降低开销
 */

const (
	MSQL_DIALECT = "mysql"
)

// SQLDialectFactory 是 SQL 方言的工厂
type SQLDialectFactory struct {
	mutex     sync.Mutex
	instances map[string]SQLDialect
}

// NewSQLDialectFactory 返回一个新的 SQL 方言工厂
func NewSQLDialectFactory() *SQLDialectFactory {
	return &SQLDialectFactory{
		instances: make(map[string]SQLDialect),
	}
}

// GetDialect 返回指定名称的 SQL 方言
func (f *SQLDialectFactory) GetDialect(className string) (SQLDialect, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	dialect, ok := f.instances[className]
	if !ok {
		newDialect, err := f.newDialect(className)
		if err != nil {
			return nil, fmt.Errorf("failed to create instance of dialect %q:%v", className, err)
		}
		f.instances[className] = newDialect
		return newDialect, nil
	}
	return dialect, nil
}

// sqlDialectTypes 是已知的 SQL 方言类型列表
func (f *SQLDialectFactory) newDialect(className string) (SQLDialect, error) {
	dialectType, ok := sqlDialectTypes[className]
	if !ok {
		return nil, fmt.Errorf("unknown sql dialect type %q", className)
	}
	newDialect, err := dialectType.CreateDialect()
	if err != nil {
		return nil, err
	}
	return newDialect, nil
}

// SQLDialectType 是 SQL 方言需要实现的结构体类型
type SQLDialectType struct {
	Name          string
	CreateDialect func() (SQLDialect, error)
}

// sqlDialectTypes 是已知的 SQL 方言类型列表
var sqlDialectTypes = map[string]SQLDialectType{
	MSQL_DIALECT: {
		Name: MSQL_DIALECT,
		CreateDialect: func() (SQLDialect, error) {
			return new(MySQLDialect), nil
		},
	},
}
