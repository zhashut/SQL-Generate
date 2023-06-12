package consts

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/12
 * Time: 18:06
 * Description: 通用常量
 */

type CommonConst string

const (
	SORT_ORDER_ASC  CommonConst = "ascend"  // 升序
	SORT_ORDER_DESC CommonConst = "descend" // 降序
)

var CommonConstToString = map[CommonConst]string{
	SORT_ORDER_ASC:  "ascend",
	SORT_ORDER_DESC: "descend",
}

var CommonStringToConst = map[string]CommonConst{
	"ascend":  SORT_ORDER_ASC,
	"descend": SORT_ORDER_DESC,
}
