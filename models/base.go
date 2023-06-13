package models

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/12
 * Time: 17:58
 * Description: 分页请求
 */

type PageRequest struct {
	Pages     int64  `json:"current"`   // 当前页号
	PageSize  int64  `json:"pageSize"`  // 页面大小
	SortField string `json:"sortField"` // 排序字段
	SortOrder string `json:"sortOrder"` // 排序顺序(默认升序)
}

// OnlyIDRequest 传递单个 id 的请求-通用
type OnlyIDRequest struct {
	ID int64 `form:"id" json:"id"`
}
