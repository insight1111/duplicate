package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Path   string
	SHA256 string
}

func dirList(startDir string) (result []File, err error) {
	_result, err := ioutil.ReadDir(startDir)
	pwd, _ := os.Getwd()
	if err != nil {
		return
	}
	files := []File{}
	for _, file := range _result {
		if file.IsDir() {
			_t, _ := dirList(filepath.Join(startDir, file.Name()))
			files = append(files, _t...)
			continue
		}
		path := filepath.Join(pwd, startDir, file.Name())
		f := File{
			Path:   path,
			SHA256: getSHA(path),
		}
		files = append(files, f)
	}
	result = files
	return
}

func getSHA(filename string) string {
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

func main() {
	result, _ := dirList("testdir")

	fmt.Println(result)
	// fmt.Println(getSHA("testdir/a.txt"))
}
