package server

import (
	"context"
	"fmt"
	"sql_generate/global"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: zhashut
 * Date: 2024/12/14
 * Time: 17:40
 * Description: No Description
 */

func GetAddress(ctx context.Context) (string, error) {
	address := fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port)
	return address, nil
}
