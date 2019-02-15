package util

import "time"

func GetCurrentSeconds() int64 {
	return int64(time.Now().UnixNano() / int64(time.Second))
}

func GetCurrentNanoSeconds() int64 {
	return time.Now().UnixNano()
}

func GetCurrentMicroSeconds() int64 {
	return GetCurrentNanoSeconds() / int64(time.Microsecond)
}

func GetCurrentMilliSeconds() int64 {
	return GetCurrentNanoSeconds() / int64(time.Millisecond)
}

func CheckTimeLargeOrEqual(nowTime, checkTime, interval int64) bool {
	return nowTime-checkTime >= interval
}

func CheckTimeLessOrEqual(nowTime, checkTime, interval int64) bool {
	return nowTime-checkTime <= interval
}
