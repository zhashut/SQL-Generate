package builder

import "encoding/json"

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 8:39
 * Description: 数据 JSON 生成器
 */

func BuildJSON(dataList []map[string]interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(dataList, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
