package config

import (
	"strconv"
	"strings"
	"sync"
)

type ConfigManager struct {
	conf ServiceConfig
}

var m *ConfigManager
var once sync.Once

func GetInstance() *ConfigManager {
	once.Do(func() {
		m = &ConfigManager{}
	})
	return m
}

func (self *ConfigManager) InitManager() error {
	return self.conf.Init("./conf/short_url_conf.ini")
}

func (self ConfigManager) GetConfig(confName string) (string, error) {
	return self.conf.GetConfig(confName)
}

func (self ConfigManager) GetInt(confName string) (int, error) {
	value, err := self.conf.GetConfig(confName)
	if nil != err {
		return 0, err
	}
	return strconv.Atoi(value)
}

func (self ConfigManager) GetConfigArray(configName string) ([]string, error) {
	str, err := self.conf.GetConfig(configName)
	if nil != err {
		return nil, err
	}
	return strings.Split(str, ","), err
}
