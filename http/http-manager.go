package http

import (
	"github.com/service-kit/short-url/common"
	"github.com/service-kit/short-url/config"
	"github.com/service-kit/short-url/log"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
)

type HttpManager struct {
	addr           string
	shortUrlHeader string
	wg             *sync.WaitGroup
}

var m *HttpManager
var once sync.Once
var logger *zap.Logger

func GetInstance() *HttpManager {
	once.Do(func() {
		m = &HttpManager{}
	})
	return m
}

func (self *HttpManager) InitManager(wg *sync.WaitGroup) error {
	logger = log.GetInstance().GetLogger()
	self.wg = wg
	var err error = nil
	self.addr, err = config.GetInstance().GetConfig("SHORT_URL_HTTP_ADDR")
	if nil != err {
		return err
	}
	self.shortUrlHeader, err = config.GetInstance().GetConfig("SHORT_URL_HEADER")
	if "" == self.shortUrlHeader {
		logger.Warn("short url header nil")
		self.shortUrlHeader = common.SHORT_URL_HEADER
	}
	self.wg.Add(1)
	http.HandleFunc("/", handleShortUrlRequest)
	go self.startHttpServer()
	return nil
}

func (self *HttpManager) startHttpServer() {
	logger.Info("Start Http Server", zap.String("addr", self.addr))
	defer self.wg.Done()
	http.ListenAndServe(self.addr, nil)
}

func (self *HttpManager) outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
