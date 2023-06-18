package consts

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/19
 * Time: 1:45
 * Description: 代码生成器常量/枚举
 */

const (
	GoTemplate = `// {{.Comment}}
type {{.EntityName}} struct {
  {{- range .FieldList }}
  // {{.Comment}}
  {{.Name}} {{.Type}}
  {{- end }}
}`
	JavaTemplate = `/**
 * {{.Comment}}
 */
public class {{.EntityName}} {
  {{- range .FieldList }}
  /** {{.Comment}} */
  private {{.Type}} {{.Name}};
  {{- end }}
}`
	JavaField = "Java"
	GoField   = "Go"
)

const (
	SQLIndex = iota
	JavaIndex
	GoIndex
)
