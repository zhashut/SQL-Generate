package consts

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/9
 * Time: 20:13
 * Description: 用户常量
 */

type UserConst string

const (
	USER_LOGIN_STATE = "userLoginState" // 用户登录态键
	DEFAULT_ROLE     = "user"           // 默认权限
	ADMIN_ROLE       = "admin"          // 管理员权限
	SYSTEM_USER_ID   = 0                // 系统用户 id（虚拟用户）
	SALT             = "shut"           // 盐值
)
