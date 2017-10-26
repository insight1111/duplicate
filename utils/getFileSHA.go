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
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
func GetFileSHA2(filename string, sha chan string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	// return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	sha <- strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
