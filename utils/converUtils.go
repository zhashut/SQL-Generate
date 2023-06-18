package utils

import "strings"

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/18
 * Time: 21:31
 * Description: 转换相关的工具方法
 */

// CamelCase 将蛇形命名转为驼峰命名
func CamelCase(s string) string {
	words := strings.Split(s, "_")
	for i, w := range words {
		words[i] = strings.Title(strings.ToLower(w))
	}
	return strings.Join(words, "")
}
