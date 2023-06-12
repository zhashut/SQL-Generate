package db

import (
	"context"
	"fmt"
	"sql_generate/global"
	"sql_generate/models"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/11
 * Time: 19:17
 * Description: No Description
 */

func AddTableInfo(ctx context.Context, t *models.TableInfo) (bool, error) {
	if res := global.DB.Where("name = ? and userId = ?", t.Name, t.UserId).First(&models.TableInfo{}); res.RowsAffected > 0 {
		return false, fmt.Errorf("表名已存在")
	}
	global.DB.Save(&t)
	return true, nil
}
