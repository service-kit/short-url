package util

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src rand.Source

func init() {
	src = rand.NewSource(time.Now().UnixNano())
}

func RandString(n int) string {
	b := make([]byte, n)
	rand_times := n / 63
	if n%63 != 0 {
		rand_times += 1
	}
	for rt := 0; rt < rand_times; rt++ {
		for i, cache, remain := n-63*rt-1, src.Int63(), letterIdxMax; i >= 0; {
			if remain == 0 {
				cache, remain = src.Int63(), letterIdxMax
			}
			if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
				b[i] = letterBytes[idx]
				i--
			}
			cache >>= letterIdxBits
			remain--
		}
	}
	return string(b)
}

func NewUUID() string {
	uid := uuid.New()
	return uid.String()
}
