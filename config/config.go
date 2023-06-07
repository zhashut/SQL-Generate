package config

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 10:18
 * Description: No Description
 */

type ServerConfig struct {
	Host        string      `mapstructure:"host" json:"host"`
	Port        int         `mapstructure:"port" json:"port"`
	MySQLConfig MySQLConfig `mapstructrue:"mysql" json:"mysql"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}
