package redis

import (
	"github.com/service-kit/short-url/config"
	"github.com/service-kit/short-url/log"
	"go.uber.org/zap"
	"sync"
)

type RedisManager struct {
	redisPool RedisPool
}

var m *RedisManager
var once sync.Once
var logger *zap.Logger

func GetInstance() *RedisManager {
	once.Do(func() {
		m = &RedisManager{}
	})
	return m
}

func (self *RedisManager) InitManager() error {
	logger = log.GetInstance().GetLogger()
	host, err := config.GetInstance().GetConfig("REDIS_ADDR")
	if nil != err {
		logger.Error("RedisMAnager InitManager fail! Can not find conf REDIS_HOST!!!")
		return err
	}
	passwd, err := config.GetInstance().GetConfig("REDIS_PASSWD")
	if nil != err {
		logger.Error("RedisMAnager InitManager fail! Can not find conf REDIS_PASSWD!!!")
		return err
	}
	self.redisPool.Init(host, passwd)
	return nil
}

func (self *RedisManager) SetStringValue(key, value string) (err error) {
	return self.redisPool.SetStringValue(key, value)
}

func (self *RedisManager) SetStringValueWithExpireTime(key, value string, expireTime int64) (err error) {
	return self.redisPool.SetStringValueWithExpireTime(key, value, expireTime)
}

func (self *RedisManager) SetMultiValueWithExpireTime(kvMap map[string]string, ktMap map[string]int64) (err error) {
	return self.redisPool.SetMultiValueWithExpireTime(kvMap, ktMap)
}

func (self *RedisManager) GetStringValue(key string) (out string, err error) {
	return self.redisPool.GetStringValue(key)
}

func (self *RedisManager) GetMultiValue(keys ...interface{}) (out map[string]string, err error) {
	return self.redisPool.GetMultiValue(keys)
}

func (self *RedisManager) GetKeyExpire(key string) (out int64, err error) {
	return self.redisPool.GetKeyExpire(key)
}
