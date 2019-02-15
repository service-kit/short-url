package http

import (
	"github.com/service-kit/short-url/common"
	"github.com/service-kit/short-url/storage"
	"github.com/service-kit/short-url/util"
	"go.uber.org/zap"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func handleShortUrlRequest(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.RequestURI)
	if "/" != r.RequestURI {
		short_url := r.RequestURI[1:]
		if strings.Contains(short_url, "cache/") && strings.Contains(short_url, ".jpg") {
			downJpg(w, short_url)
			return
		}
		if common.FAVICON_ICO == short_url {
			downFaviconIco(w)
			return
		}
		logger.Info("short url request", zap.String("short url", short_url))
		original_url, err := storage.GetInstance().GetOriginalUrl(short_url)
		if nil != err || "" == original_url {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		logger.Info("redirect to original url", zap.String("original url", original_url))
		http.Redirect(w, r, original_url, http.StatusMovedPermanently)
		return
	}
	r.ParseForm()
	var form url.Values
	if http.MethodGet == r.Method {
		form = r.Form
	} else if http.MethodPost == r.Method {
		form = r.PostForm
	}
	original_url := form.Get("original_url")
	if "" == original_url {
		logger.Info("get index html")
		GetInstance().outputHTML(w, r, "./html/index.html")
		return
	}
	short_url := form.Get("short_url")
	if "" == short_url {
		logger.Info("use system create short url")
		short_url = util.BuildShortUrl(original_url)
	} else {
		res, _ := storage.GetInstance().GetOriginalUrl(short_url)
		if original_url != res {
			w.Write([]byte(short_url + " has exist"))
			return
		}
	}
	short_url_info := new(common.ShortUrlInfo)
	short_url_info.OriginalUrl = original_url
	short_url_info.ShortUrl = short_url
	storage.GetInstance().StorageShortUrlInfo(short_url_info)
	logger.Info("register", zap.Any("param", form))
	fullShortUrl := GetInstance().shortUrlHeader + short_url
	jpgData, err := util.BuildQRCodeJpg(fullShortUrl)
	if nil != err {
		return
	}
	cacheFileName := "./cache/" + short_url + ".jpg"
	util.SaveFile("./html/"+cacheFileName, jpgData)
	err = fillRegisterResultHtml(w, original_url, fullShortUrl, cacheFileName)
	if nil != err {
		logger.Error("fill register result html err", zap.Error(err))
	}
}

func downJpg(w http.ResponseWriter, url string) error {
	cacheFile, err := os.OpenFile("./html/"+url, os.O_RDONLY, 0777)
	if nil != err {
		return err
	}
	w.Header().Set("Content-Type", "image/jpg")
	_, err = io.Copy(w, cacheFile)
	return err
}

func downFaviconIco(w http.ResponseWriter) error {
	f, err := os.OpenFile("./html/favicon.ico", os.O_RDONLY, 0777)
	if nil != err {
		return err
	}
	w.Header().Set("Content-Type", "image/icon")
	_, err = io.Copy(w, f)
	return err
}

func fillRegisterResultHtml(w http.ResponseWriter, oriUrl, shortUrl, qrjpg string) error {
	return fillHtmlData(w, map[string]string{"ORIURL": oriUrl, "SHORTURL": shortUrl, "QRJPG": qrjpg}, "./html/register_result.html")
}

func fillHtmlData(w http.ResponseWriter, data map[string]string, htmls ...string) error {
	t, err := template.ParseFiles(htmls...)
	if nil != err {
		logger.Error("template parse files err", zap.Error(err))
		return err
	}
	err = t.Execute(w, data)
	if nil != err {
		logger.Error("template execute err", zap.Error(err))
	}
	return err
}
