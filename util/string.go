package util

import (
	"fmt"
)

func StringFixedLen(str string, size int) string {
	if size < 1 {
		return ""
	}
	bm := make([]byte, size)
	for i := 0; i < size; i++ {
		bm[i] = '0'
	}
	bm2 := []byte(fmt.Sprintf("%x", str))
	pos := len(bm) - 1
	for i := len(bm2) - 1; i >= 0; i-- {
		bm[pos] = bm2[i]
		pos--
		if 0 == pos {
			break
		}
	}
	return string(bm)
}
