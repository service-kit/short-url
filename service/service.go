package service

import (
	"github.com/service-kit/short-url/config"
	"github.com/service-kit/short-url/http"
	"github.com/service-kit/short-url/log"
	"github.com/service-kit/short-url/redis"
	"github.com/service-kit/short-url/storage"
	"go.uber.org/zap"
	"sync"
	"time"
)

var wg sync.WaitGroup
var logger *zap.Logger

func StartService() {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("service err", zap.Any("panic recover", e))
		}
		logger.Error("service start error,may shut dowm after 3 seconds")
		log.GetInstance().FinishProcess()
		time.Sleep(3 * time.Second)
	}()
	err := initManager()
	if nil != err {
		logger.Error("init manager err", zap.Error(err))
	} else {
		wg.Wait()
	}
}

func initManager() error {
	var err error = nil
	err = log.GetInstance().InitManager()
	if nil != err {
		return err
	}
	logger = log.GetInstance().GetLogger()
	err = config.GetInstance().InitManager()
	if nil != err {
		return err
	}
	err = redis.GetInstance().InitManager()
	if nil != err {
		return err
	}
	err = storage.GetInstance().InitManager()
	if nil != err {
		return err
	}
	wg.Add(1)
	err = http.GetInstance().InitManager(&wg)
	if nil != err {
		return err
	}
	return err
}
