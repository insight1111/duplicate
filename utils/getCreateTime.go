package utils

import (
	"log"
	"os"
	"syscall"
	"time"
)

// GetCreateTime ファイル名を受け取り、その作成日時を返す
func GetCreateTime(file string) time.Time {
	f, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}
	stat := f.Sys().(*syscall.Win32FileAttributeData).CreationTime.Nanoseconds()
	return time.Unix(0, stat)
}
