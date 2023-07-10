package ses

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"sql_generate/consts"
	"sql_generate/models"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/7/8
 * Time: 18:13
 * Description: 设置session
 */

// SetSession 设置session
func SetSession(c *gin.Context, data interface{}) {
	s := sessions.Default(c)
	user := data.(*models.User)
	marshal, _ := json.Marshal(user)
	s.Set(consts.USER_LOGIN_STATE, marshal)
	_ = s.Save()
}

// GetSession 获取session
func GetSession(c *gin.Context, key string) *models.User {
	s := sessions.Default(c)
	var user *models.User
	value := s.Get(key)
	if value == nil {
		return nil
	}
	u := value.([]byte)
	if err := json.Unmarshal(u, &user); err != nil {
		return nil
	}
	return user
}

// DeleteSession 删除session
func DeleteSession(c *gin.Context, key string) {
	s := sessions.Default(c)
	s.Delete(key)
	_ = s.Save()
}
