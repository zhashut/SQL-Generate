package schema

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/5
 * Time: 14:46
 * Description: 表概要
 */

// TableSchema 表概要
type TableSchema struct {
	DBName       string   `form:"dbName" json:"dbName"`             // 库名
	TableName    string   `form:"tableName" json:"tableName"`       // 表名
	TableComment string   `form:"tableComment" json:"tableComment"` // 表注释
	MockNum      int32    `form:"mockNum" json:"mockNum"`           // 模拟数据条数
	FieldList    []*Field `form:"fieldList" json:"fieldList"`       // 列表信息
}

// Field 列信息列表
type Field struct {
	FieldName     string      `form:"fieldName" json:"fieldName"`         // 字段名
	FieldType     string      `form:"fieldType" json:"fieldType"`         // 字段类型
	DefaultValue  string      `form:"defaultValue" json:"defaultValue"`   // 默认值
	NotNull       bool        `form:"notNull" json:"notNull"`             // 是否非空
	Comment       string      `form:"comment" json:"comment"`             // 注释（字段中文名）
	PrimaryKey    bool        `form:"primaryKey" json:"primaryKey"`       // 是否为主键
	AutoIncrement bool        `form:"autoIncrement" json:"autoIncrement"` // 是否自增
	MockType      string      `form:"mockType" json:"mockType"`           // 模拟类型（随机、图片、规则、词库）
	MockParams    interface{} `form:"mockParams" json:"mockParams"`       // 模拟参数,词库和递增会传int类型,其他传string,所以用interface{}
	OnUpdate      string      `form:"onUpdate" json:"onUpdate"`           // 附加条件
}
