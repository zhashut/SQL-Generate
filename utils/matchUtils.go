package utils

import (
	"regexp"
	"strconv"
	"time"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/17
 * Time: 18:43
 * Description: 判断是否匹配
 */

var datePatterns = []string{"2006-01-02", "2006年01月02日", "2006-01-02 15:04:05", "2006-01-02 15:04", "2006/01/02", "2006/01/02 15:04:05", "2006/01/02 15:04", "20060102"}

// IsNumeric 判断是否是小数
func IsNumeric(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		_, err := strconv.ParseInt(value, 10, 64)
		return err == nil
	}
	return true
}

// IsDouble 判断是否是小数
func IsDouble(value string) bool {
	pattern := regexp.MustCompile("[0-9]+[.]{0,1}[0-9]*[dD]{0,1}")
	return pattern.MatchString(value)
}

// IsDate 判断是否是日期
func IsDate(value string) bool {
	if value == "" {
		return false
	}
	for _, pattern := range datePatterns {
		_, err := time.Parse(pattern, value)
		if err == nil {
			return true
		}
	}
	return false
}
