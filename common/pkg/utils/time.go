package utils

import "time"

const (
	TimeFormatSecond = "2006-01-02 15:04:05"
	TimeFormatMinute = "2006-01-02 15:04"
	TimeFormatDateV1 = "2006-01-02"
	TimeFormatDateV2 = "2006_01_02"
	TimeFormatDateV3 = "20060102150405"
	TimeFormatDateV4 = "2006/01/02 - 15:04:05.000"
)

func GetTodayUnix() int64 {
	currentTime := time.Now()
	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Unix()
}
