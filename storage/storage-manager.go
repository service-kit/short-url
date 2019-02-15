package storage

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/service-kit/short-url/common"
	"github.com/service-kit/short-url/config"
	"github.com/service-kit/short-url/log"
	"github.com/service-kit/short-url/redis"
	"go.uber.org/zap"
	"sync"
)

type StorageManager struct {
	MysqlParam  string
	mysqlSwitch bool
}

var m *StorageManager
var once sync.Once
var logger *zap.Logger

func GetInstance() *StorageManager {
	once.Do(func() {
		m = &StorageManager{}
	})
	return m
}

func (self *StorageManager) InitManager() error {
	logger = log.GetInstance().GetLogger()
	swi, _ := config.GetInstance().GetInt("MYSQL_SWITCH")
	if common.SWITHC_ON != swi {
		self.mysqlSwitch = false
		return nil
	}
	self.mysqlSwitch = true
	err := self.loadConfig()
	if nil != err {
		return err
	}
	return self.initDB()
}

func (self *StorageManager) initDB() error {
	db, err := self.getDBCon()
	if nil != err {
		return err
	}
	defer db.Close()
	if !db.HasTable(&common.ShortUrlInfo{}) {
		db.CreateTable(&common.ShortUrlInfo{})
	}
	return db.Error
}

func (self *StorageManager) getDBCon() (*gorm.DB, error) {
	if !self.mysqlSwitch {
		return nil, errors.New("mysql switch off")
	}
	return gorm.Open("mysql", self.MysqlParam)
}

func (self *StorageManager) exist(data interface{}) bool {
	db, err := self.getDBCon()
	if nil != err {
		return false
	}
	defer db.Close()
	return !db.First(data).RecordNotFound()
}

func (self *StorageManager) insert(data interface{}) error {
	db, err := self.getDBCon()
	if nil != err {
		return err
	}
	defer db.Close()
	return db.Create(data).Error
}

func (self *StorageManager) update(data interface{}) error {
	db, err := self.getDBCon()
	if nil != err {
		return err
	}
	defer db.Close()
	return db.Save(data).Error
}

func (self *StorageManager) delete(data interface{}) error {
	db, err := self.getDBCon()
	if nil != err {
		return err
	}
	defer db.Close()
	return db.Delete(data).Error
}

func (self *StorageManager) query(data interface{}) error {
	db, err := self.getDBCon()
	if nil != err {
		return err
	}
	defer db.Close()
	return db.Where(data).First(data).Error
}

func (self *StorageManager) loadConfig() error {
	addr, err := config.GetInstance().GetConfig("DB_ADDR")
	if nil != err {
		return err
	}
	user, err := config.GetInstance().GetConfig("DB_USER")
	if nil != err {
		return err
	}
	passwd, err := config.GetInstance().GetConfig("DB_PASSWD")
	if nil != err {
		return err
	}
	dbname, err := config.GetInstance().GetConfig("DB_DBNAME")
	if nil != err {
		return err
	}
	self.MysqlParam = GenerateMysqlParam(addr, user, passwd, dbname)
	return err
}

func GenerateMysqlParam(addr, user, passwd, dbName string) string {
	return user + ":" + passwd + "@tcp(" + addr + ")/" + dbName + "?charset=utf8"
}

func (self *StorageManager) raw(sql string, out interface{}) error {
	db, err := self.getDBCon()
	if nil != err {
		return err
	}
	defer db.Close()
	return db.Exec(sql).Find(out).Error
}

func (self *StorageManager) selectWithOrderAndLimit(cond, order string, limit int, out interface{}) error {
	db, err := self.getDBCon()
	if nil != err {
		return err
	}
	defer db.Close()
	return db.Where(cond).Order(order).Limit(limit).Find(out).Error
}

func (self *StorageManager) selectAll(out interface{}) error {
	db, err := self.getDBCon()
	if nil != err {
		return err
	}
	defer db.Close()
	return db.Find(out).Error
}

func (self StorageManager) generateShortUrlKey(clientID string) string {
	return "short_url:" + clientID
}

func (self *StorageManager) StorageShortUrlInfo(short_url *common.ShortUrlInfo) (bool, error) {
	if self.mysqlSwitch {
		exist := self.exist(short_url)
		if exist {
			return true, errors.New("account exist")
		}
		err := self.insert(short_url)
		if nil != err {
			logger.Error("storage account to db err", zap.Error(err))
			return false, err
		}
	}
	original_url, _ := redis.GetInstance().GetStringValue(self.generateShortUrlKey(short_url.ShortUrl))
	if "" != original_url {
		return false, errors.New("account exist")
	}
	return false, redis.GetInstance().SetStringValue(self.generateShortUrlKey(short_url.ShortUrl), short_url.OriginalUrl)
}

func (self *StorageManager) GetOriginalUrl(short_url string) (string, error) {
	original_url, err := redis.GetInstance().GetStringValue(self.generateShortUrlKey(short_url))
	if nil == err && "" != original_url {
		return original_url, err
	}
	if !self.mysqlSwitch {
		return "", errors.New("not register")
	}
	short_url_info := new(common.ShortUrlInfo)
	short_url_info.ShortUrl = short_url
	err = self.query(short_url_info)
	if nil != err {
		logger.Error("query short url info from db err", zap.Error(err))
		return "", err
	}
	if "" != short_url_info.ShortUrl {
		logger.Info("sync short url info to redis ", zap.String("id", short_url_info.OriginalUrl))
		err = redis.GetInstance().SetStringValue(self.generateShortUrlKey(short_url_info.OriginalUrl), short_url_info.ShortUrl)
		if nil != err {
			logger.Error("sync to redis short url info err", zap.Error(err))
		}
	}
	return short_url_info.ShortUrl, nil
}

func (self *StorageManager) LoadAllShortUrlData() ([]common.ShortUrlInfo, error) {
	if !self.mysqlSwitch {
		return nil, nil
	}
	var infos []common.ShortUrlInfo
	err := self.selectAll(&infos)
	if nil != err {
		return nil, err
	}
	return infos, nil
}
