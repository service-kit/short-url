package config

import (
	"errors"
	"github.com/Unknwon/goconfig"
	"github.com/service-kit/short-url/util"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var logger *zap.Logger

type ServiceConfig struct {
	conf            *goconfig.ConfigFile
	confName        string
	isLoadSucc      bool
	fileLastModTime int64
	lastCheckTime   int64
}

func (self *ServiceConfig) Init(confFile string) error {
	self.confName = confFile
	var err error = nil
	self.conf, err = goconfig.LoadConfigFile(confFile)
	if nil != err {
		return err
	}
	self.isLoadSucc = true
	self.fileLastModTime, _ = util.GetFileModTime(self.confName)
	self.CheckConfigFile()
	return err
}

func (self ServiceConfig) GetConfig(confName string) (string, error) {
	if !self.isLoadSucc {
		return "", errors.New("GetConfig fail! Is not load success!!!")
	}
	return self.conf.GetValue("", confName)
}

func (self ServiceConfig) GetInt(confName string) (int, error) {
	if !self.isLoadSucc {
		return 0, errors.New("GetConfig fail! Is not load success!!!")
	}
	val, err := self.conf.GetValue("", confName)
	if nil != err {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (self ServiceConfig) GetBool(confName string) (bool, error) {
	if !self.isLoadSucc {
		return false, errors.New("GetConfig fail! Is not load success!!!")
	}
	val, err := self.conf.GetValue("", confName)
	if nil != err {
		return false, err
	}
	return strconv.ParseBool(val)
}

func (self *ServiceConfig) CheckConfigFile() {
	go func() {
		for {
			fileModTime, err := util.GetFileModTime(self.confName)
			if nil != err {
				continue
			}
			if fileModTime != self.fileLastModTime {
				self.conf.Reload()
				self.fileLastModTime = fileModTime
			}
			time.Sleep(time.Second)
		}
	}()
}
