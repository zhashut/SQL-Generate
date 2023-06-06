package utils

import "sql_generate/consts"

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/6
 * Time: 16:27
 * Description: No Description
 */

func GetEnumByValue(typeEnum interface{}) string {
	return string(typeEnum.(consts.MockTypeEnum))
}
