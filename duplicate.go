package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func dupList(fileList []File) (result [][]File, err error) {
	return
}

func main() {
	result, _ := dirList("testdir")
	list, _ := dupList(result)
	fmt.Println(result)
	fmt.Println(list)
	// fmt.Println(getSHA("testdir/a.txt"))
}
