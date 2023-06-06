package main

import (
	"fmt"
	"sql_generate/consts"
	"sql_generate/utils"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/5/14
 * Time: 22:10
 * Description: No Description
 */

func main() {
	value := utils.GetEnumByValue(consts.NONE)
	fmt.Println(value)
}
