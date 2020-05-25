package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// 系统配置信息
type Configurations struct {
	Server *ServerConfig
	DB []*DBConfig
	Mail *MailConfig
	Admin *AdminConfig
}

// 服务器配置
type ServerConfig struct {
	Port string
}

// 数据源
type DBConfig  struct {
	Name string
	Dialect string
	DSN string
}

// 邮箱配置
type MailConfig struct {
	Addr string
	User string
	Password string
}

// 管理员初始信息
type AdminConfig struct {
	Email string
	Password string
	NickName string `yaml:"nick_name"`
	Phone string
}

var cfg *Configurations

func Load() error {
	data, err := ioutil.ReadFile("conf/lchat.yaml")
	if err != nil {
		return err
	}
	cfg = &Configurations{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}
	return nil
}

// 获取加载后的配置信息
func Get() *Configurations {
	return cfg
}
