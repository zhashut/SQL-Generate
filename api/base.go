package api

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/13
 * Time: 6:46
 * Description: 通用
 */

type PageInfo struct {
	Records  interface{} `json:"records"`
	Pages    int64       `json:"current"`
	PageSize int64       `json:"size"`
	Total    int64       `json:"total"`
}

// 这个函数解决了将列索引数字与 Excel 中列名的字母形式之间进行转换的问题。
// 在 DownloadDataExcel 函数中，我们将这个函数用于设置表头和数据行时指定单元格位置。
func columnIndexToExcelColumn(index int) string {
	columnName := ""
	for index >= 0 {
		columnName = string('A'+(index%26)) + columnName
		index = (index / 26) - 1
	}
	return columnName
}
