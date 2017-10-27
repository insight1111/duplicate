package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

// GetFileSHA 指定ファイルのSHA256を算出
func GetFileSHA(filename string) string {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		log.Println("sha error:", err)
		return ""
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Println(err)
	}
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

//GetFileSHA2 指定ファイルのSHA256を算出(goroutineバージョン)
func GetFileSHA2(filename string, sha chan string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Println(err)
	}
	// return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	sha <- strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
