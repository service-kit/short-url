package util

import (
	"bytes"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/jpeg"
	"image/png"
)

func BuildQRCodePng(str string) ([]byte, error) {
	code, err := qr.Encode(str, qr.L, qr.Unicode)
	if nil != err {
		return nil, err
	}
	code, err = barcode.Scale(code, 100, 100)
	if nil != err {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = png.Encode(buf, code)
	if nil != err {
		return nil, err
	}
	return buf.Bytes(), err
}

func BuildQRCodeJpg(str string) ([]byte, error) {
	code, err := qr.Encode(str, qr.L, qr.Unicode)
	if nil != err {
		return nil, err
	}
	code, err = barcode.Scale(code, 100, 100)
	if nil != err {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, code, &jpeg.Options{100})
	if nil != err {
		return nil, err
	}
	return buf.Bytes(), err
}
