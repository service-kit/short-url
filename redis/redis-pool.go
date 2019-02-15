package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/service-kit/short-url/config"
	"go.uber.org/zap"
	"time"
)

const (
	REDIS_POOL_MAX_IDLE     = 512
	REDIS_POOL_MAX_ACTIVE   = 1024
	REDIS_POOL_IDLE_TIMEOUT = 240
)

type RedisPool struct {
	redisPool *redis.Pool
	isInit    bool
}

func (self *RedisPool) Init(host, passwd string) {
	self.redisPool = self.newPool(host, passwd)
	self.isInit = true
}

func (self *RedisPool) newPool(host, passwd string) *redis.Pool {
	maxIdle, err := config.GetInstance().GetInt("REDIS_POOL_MAX_IDLE")
	if nil != err {
		logger.Error("get config fail", zap.String("cfg", "REDIS_POOL_MAX_IDLE"), zap.Error(err))
		maxIdle = REDIS_POOL_MAX_IDLE
	}
	maxActive, err := config.GetInstance().GetInt("REDIS_POOL_MAX_ACTIVE")
	if nil != err {
		logger.Error("get config fail", zap.String("cfg", "REDIS_POOL_MAX_ACTIVE"), zap.Error(err))
		maxActive = REDIS_POOL_MAX_ACTIVE
	}
	idleTimeout, err := config.GetInstance().GetInt("REDIS_POOL_IDLE_TIMEOUT")
	if nil != err {
		logger.Error("get config fail", zap.String("cfg", "REDIS_POOL_IDLE_TIMEOUT"), zap.Error(err))
		idleTimeout = REDIS_POOL_IDLE_TIMEOUT
	}
	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			if "" != passwd {
				if _, err := c.Do("AUTH", passwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func (self *RedisPool) setValue(key string, value interface{}) error {
	if !self.isInit {
		return errors.New(REDIS_UNAVAILABLE)
	}
	conn := self.redisPool.Get()
	if nil == conn {
		return errors.New(REDIS_UNAVAILABLE)
	}
	if nil != conn.Err() {
		return errors.New(REDIS_UNAVAILABLE)
	}
	defer conn.Close()
	ret, err := conn.Do("SET", key, value)
	if nil != err {
		return err
	}
	if ret != "OK" {
		return errors.New("setValue fail! do set return err!!!")
	}
	return err
}

func (self *RedisPool) setValueWithExpireTime(key string, value interface{}, expireTime int64) error {
	if !self.isInit {
		return errors.New(REDIS_UNAVAILABLE)
	}
	conn := self.redisPool.Get()
	if nil == conn {
		return errors.New(REDIS_UNAVAILABLE)
	}
	if nil != conn.Err() {
		return errors.New(REDIS_UNAVAILABLE)
	}
	defer conn.Close()
	ret, err := conn.Do("SET", key, value)
	if nil != err {
		return err
	}
	if ret != "OK" {
		return errors.New("setValueWithExpireTime fail! do set return err!!!")
	}
	ret, err = conn.Do("EXPIRE", key, expireTime)
	if nil != err {
		return err
	}
	if ret != int64(1) {
		return errors.New("setValueWithExpireTime fail! do expire return err!!!")
	}
	return err
}

func (self *RedisPool) SetStringValue(key, value string) error {
	err := self.setValue(key, value)
	if nil != err {
		return err
	}
	return err
}

func (self *RedisPool) SetStringValueWithExpireTime(key, value string, expireTime int64) error {
	err := self.setValueWithExpireTime(key, value, expireTime)
	if nil != err {
		return err
	}
	return err
}

func (self *RedisPool) SetMultiValueWithExpireTime(kvMap map[string]string, ktMap map[string]int64) error {
	err := self.setMultiValueWithExpireTime(kvMap, ktMap)
	if nil != err {
		return err
	}
	return err
}

func (self *RedisPool) getValue(key string) (interface{}, error) {
	if !self.isInit {
		return "", errors.New(REDIS_UNAVAILABLE)
	}
	conn := self.redisPool.Get()
	if nil == conn {
		return "", errors.New(REDIS_UNAVAILABLE)
	}
	if nil != conn.Err() {
		return "", errors.New(REDIS_UNAVAILABLE)
	}
	defer conn.Close()
	return conn.Do("GET", key)
}

func (self *RedisPool) getMultiValue(keys ...interface{}) (interface{}, error) {
	if !self.isInit {
		return "", errors.New(REDIS_UNAVAILABLE)
	}
	conn := self.redisPool.Get()
	if nil == conn {
		return "", errors.New(REDIS_UNAVAILABLE)
	}
	if nil != conn.Err() {
		return "", errors.New(REDIS_UNAVAILABLE)
	}
	defer conn.Close()
	return conn.Do("MGET", keys)
}

func (self *RedisPool) setMultiValueWithExpireTime(kvMap map[string]string, ktMap map[string]int64) error {
	if !self.isInit {
		return errors.New(REDIS_UNAVAILABLE)
	}
	conn := self.redisPool.Get()
	if nil == conn {
		return errors.New(REDIS_UNAVAILABLE)
	}
	if nil != conn.Err() {
		return errors.New(REDIS_UNAVAILABLE)
	}
	if len(kvMap) != len(ktMap) {
		return errors.New("kvMap size is not match ktMap!!!")
	}
	defer conn.Close()
	var args []interface{}
	for k, v := range kvMap {
		args = append(args, k)
		args = append(args, v)
	}
	ret, err := conn.Do("MSET", args...)
	if nil != err {
		return err
	}
	if ret != "OK" {
		return errors.New("setValue fail! do set return err!!!")
	}
	for k, t := range ktMap {
		ret, err = conn.Do("EXPIRE", k, t)
		if nil != err {
			return err
		}
		if ret != int64(1) {
			return errors.New("setValueWithExpireTime fail! do expire return err!!!")
		}
	}
	return err
}

func (self *RedisPool) setMultiValue(kvMap map[string]string) error {
	if !self.isInit {
		return errors.New(REDIS_UNAVAILABLE)
	}
	conn := self.redisPool.Get()
	if nil == conn {
		return errors.New(REDIS_UNAVAILABLE)
	}
	if nil != conn.Err() {
		return errors.New(REDIS_UNAVAILABLE)
	}
	defer conn.Close()
	var args []string = make([]string, len(kvMap)*2)
	for k, v := range kvMap {
		args = append(args, k)
		args = append(args, v)
	}
	ret, err := conn.Do("MSET", args)
	if nil != err {
		return err
	}
	if ret != "OK" {
		return errors.New("setValue fail! do set return err!!!")
	}
	return err
}

func (self *RedisPool) GetStringValue(key string) (string, error) {
	return redis.String(self.getValue(key))
}

func (self *RedisPool) GetMultiValue(keys ...interface{}) (map[string]string, error) {
	return redis.StringMap(self.getMultiValue(keys))
}

func (self *RedisPool) setKeyExpireTime(key string, expireTime int64) error {
	if !self.isInit {
		return errors.New(REDIS_UNAVAILABLE)
	}
	conn := self.redisPool.Get()
	if nil == conn {
		return errors.New(REDIS_UNAVAILABLE)
	}
	if nil != conn.Err() {
		return errors.New(REDIS_UNAVAILABLE)
	}
	defer conn.Close()
	ret, err := conn.Do("EXPIRE", key, expireTime)
	if nil != err {
		return err
	}
	if ret != "OK" {
		return errors.New("setKeyExpireTime fail! do expire return err!!!")
	}
	return err
}

func (self *RedisPool) GetKeyExpire(key string) (int64, error) {
	if !self.isInit {
		return -1, errors.New(REDIS_UNAVAILABLE)
	}
	conn := self.redisPool.Get()
	if nil == conn {
		return -1, errors.New(REDIS_UNAVAILABLE)
	}
	if nil != conn.Err() {
		return -1, errors.New(REDIS_UNAVAILABLE)
	}
	defer conn.Close()
	ret, err := conn.Do("TTL", key)
	if nil != err {
		return -1, err
	}
	if ret.(int64) < 0 {
		return -1, errors.New("setKeyExpireTime fail! do expire return err!!!")
	}
	return ret.(int64), err
}
