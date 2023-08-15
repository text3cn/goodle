package types

import (
	"github.com/text3cn/goodle/container"
)

type BeforStartCallback func(c container.Container)

// DBConfig 代表数据库连接的所有配置
type DBConfig struct {
	// 以下配置关于 dsn
	WriteTimeout string `mapstructure:"write_timeout"` // 写超时时间
	Loc          string `mapstructure:"loc"`           // 时区
	Port         int    `mapstructure:"port"`          // 端口
	ReadTimeout  string `mapstructure:"read_timeout"`  // 读超时时间
	Charset      string `mapstructure:"charset"`       // 字符集
	ParseTime    bool   `mapstructure:"parse_time"`    // 是否解析时间
	Protocol     string `mapstructure:"protocol"`      // 传输协议
	Database     string `mapstructure:"database"`      // 数据库
	Collation    string `mapstructure:"collation"`     // 字符序
	Timeout      string `mapstructure:"timeout"`       // 连接超时时间
	Username     string `mapstructure:"username"`      // 用户名
	Password     string `mapstructure:"password"`      // 密码
	Driver       string `mapstructure:"driver"`        // 驱动
	Host         string `mapstructure:"host"`          // 数据库地址

	// 以下配置关于连接池
	ConnMaxIdle     int    `mapstructure:"conn_max_idle"`     // 最大空闲连接数
	ConnMaxOpen     int    `mapstructure:"conn_max_open"`     // 最大连接数
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"` // 连接最大生命周期
	ConnMaxIdletime string `mapstructure:"conn_max_idletime"` // 空闲最大生命周期

}
