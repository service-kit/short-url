package common

const (
	CODE    = "code"
	ERROR   = "error"
	TOKEN   = "token"
	SUCCESS = "success"
	FAIL    = "fail"
)

const (
	ERROR_VERIFY_NOT_PASS  = "verify not pass"
	ERROR_CAN_NOT_REGISTER = "can not register"
	ERROR_EXIST            = "register exist"
)

type ShortUrlInfo struct {
	OriginalUrl string `gorm:"primary_key;auto_increment:false"`
	ShortUrl    string `gorm:"primary_key;auto_increment:false"`
}

const (
	SWITHC_ON  = 1
	SWITHC_OFF = 0
)

const (
	LL_DEBUG = "debug"
	LL_INFO  = "info"
	LL_ERROR = "error"
)

const (
	SHORT_URL_HEADER = "http://127.0.0.1/"
	FAVICON_ICO      = "favicon.ico"
)
