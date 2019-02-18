package data

import (
	"github.com/service-kit/short-url/common"
	"github.com/service-kit/short-url/log"
	"github.com/service-kit/short-url/storage"
	"go.uber.org/zap"
	"sync"
	"time"
)

type DataManager struct {
	cacheLock      sync.RWMutex
	shortUrlMap    map[string]string
	originalUrlMap map[string]string
}

var m *DataManager
var once sync.Once
var logger *zap.Logger

func GetInstance() *DataManager {
	once.Do(func() {
		m = &DataManager{}
	})
	return m
}

func (self *DataManager) InitManager() (err error) {
	logger = log.GetInstance().GetLogger()
	if nil != err {
		return
	}
	self.loadShortUrl()
	self.checkSyncDB()
	return
}

func (self *DataManager) checkSyncDB() {
	go func() {
		for {
			self.syncToDB()
			time.Sleep(time.Second)
		}
	}()
}

func (self *DataManager) syncToDB() {

}

func (self *DataManager) loadShortUrl() error {
	urls, err := storage.GetInstance().LoadAllShortUrlData()
	if nil != err {
		return err
	}
	for _, short_url := range urls {
		self.shortUrlMap[short_url.ShortUrl] = short_url.OriginalUrl
		self.originalUrlMap[short_url.OriginalUrl] = short_url.ShortUrl
	}
	return err
}

func (self *DataManager) GetOriginalUrl(short_url string) (string, error) {
	self.cacheLock.RLock()
	original_url := self.originalUrlMap[short_url]
	if "" != original_url {
		self.cacheLock.RUnlock()
		return original_url, nil
	}
	original_url, err := storage.GetInstance().GetOriginalUrl(short_url)
	if nil != err {
		return original_url, err
	}
	self.addToCache(original_url, short_url)
	return original_url, err
}

func (self *DataManager) GetShortUrl(original_url string) (string, error) {
	self.cacheLock.RLock()
	short_url := self.originalUrlMap[original_url]
	if "" != short_url {
		self.cacheLock.RUnlock()
		return short_url, nil
	}
	original_url, err := storage.GetInstance().GetOriginalUrl(short_url)
	if nil != err {
		return short_url, err
	}
	self.addToCache(original_url, short_url)
	return short_url, err
}

func (self *DataManager) addToCache(original_url, short_url string) {
	self.cacheLock.Lock()
	defer self.cacheLock.Unlock()
	self.originalUrlMap[original_url] = short_url
	self.shortUrlMap[short_url] = original_url
}

func (self *DataManager) AddNewShortUrl(original_url, short_url string) error {
	self.addToCache(original_url, short_url)
	short_url_info := new(common.ShortUrlInfo)
	short_url_info.ShortUrl = short_url
	short_url_info.OriginalUrl = original_url
	_, err := storage.GetInstance().StorageShortUrlInfo(short_url_info)
	if nil != err {
		return err
	}
	return nil
}
