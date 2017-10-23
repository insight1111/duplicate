// 指定フォルダより重複ファイルを探しリストにする
package main

import (
	"duplicate/sha"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// File ファイルの構造体。
// Path ファイルパス
// SHA256 そのファイルのSHA256
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
			SHA256: sha.GetFileSHA(path),
		}
		files = append(files, f)
	}
	result = files
	return
}

func dupList(fileList []File) (result [][]File) {
	shaMapList := map[string][]File{}
	for _, file := range fileList {
		shaMapList[file.SHA256] = append(shaMapList[file.SHA256], file)
	}
	for _, value := range shaMapList {
		if len(value) > 1 {
			result = append(result, value)
		}
	}
	return
}

func main() {
	result, _ := dirList("testdir")
	list := dupList(result)
	fmt.Println(result)
	fmt.Println(list)
	// fmt.Println(sha.GetFileSHA("testdir/a.txt"))
}
