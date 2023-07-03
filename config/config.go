package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
)

type Config struct {
	Server Server `yaml:"server"`
	Mysql  Mysql  `yaml:"mysql"`
}

type Server struct {
	ServerAddr string `yaml:"server_addr"`
	ListenPort int32  `yaml:"listen_port"`
}

type Mysql struct {
	Path         string `yaml:"path"`
	Port         int32  `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Db           string `yaml:"db"`
	Config       string `json:"config" yaml:"config"`               // 高级配置
	MaxIdleConns int    `json:"maxIdleConns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `json:"maxOpenConns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
}

func (m *Mysql) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", m.Username, m.Password, m.Path, m.Port, m.Db, m.Config)
	//return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Db + "?" + m.Config
}

func ReadConfig(configPath string) (*Config, error) {
	var vConfig Config
	buff, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buff, &vConfig)
	if err != nil {
		return nil, err
	}
	return &vConfig, nil
}

func InitConfig(file string) *Config {
	klog.Infof("read config from %s", file)
	cfg, err := ReadConfig(file)
	if err != nil {
		klog.Fatalf("failed to read config from %s, err %s", file, err)
	}
	klog.Infof("config listed as :%+v", cfg)
	return cfg

}
