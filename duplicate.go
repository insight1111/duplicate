// 指定フォルダより重複ファイルを探しリストにする
package main

import (
	"duplicate/sha"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// File ファイルの構造体。
// Path ファイルパス
// SHA256 そのファイルのSHA256
type File struct {
	Path   string
	SHA256 string
	Size   int64
	Create time.Time
}

// DupFiles 重複ファイルの構造体
// Files 重複ファイルの集合
// TotalSize 重複ファイルの全ファイルサイズ
// WasteSize 無駄になっているサイズ
// Original 重複ファイル群の中のオリジナルファイル名(creation timeが一番若いもの)
type DupFiles struct {
	Files     []File
	TotalSize int64
	WasteSize int64
	Original  string
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
			Size:   file.Size(),
			Create: getCreateTime(path),
		}
		files = append(files, f)
	}
	result = files
	return
}

func dupList(fileList []File) []DupFiles {
	shaMapList := map[string][]File{}
	dups := []DupFiles{}
	for _, file := range fileList {
		shaMapList[file.SHA256] = append(shaMapList[file.SHA256], file)
	}
	for _, value := range shaMapList {
		if len(value) > 1 {
			d := makeDupFiles(value)
			dups = append(dups, d)
			// result = append(result, value)
		}
	}
	return dups
}

func makeDupFiles(files []File) DupFiles {
	d := DupFiles{}
	d.Files = files
	mintime := time.Now()
	var original File
	for _, f := range files {
		c := f.Create
		if c.Before(mintime) {
			original = f
			mintime = f.Create
		}
	}
	d.Original = original.Path
	totalSize := int64(0)
	for _, f := range files {
		totalSize += f.Size
	}
	d.TotalSize = totalSize
	d.WasteSize = totalSize - original.Size

	return d
}

func getCreateTime(file string) time.Time {
	f, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}
	stat := f.Sys().(*syscall.Win32FileAttributeData).CreationTime.Nanoseconds()
	return time.Unix(0, stat)
}

func main() {
	// result, _ := dirList("testdir")
	// list := dupList(result)
	// fmt.Println(result)
	// fmt.Println(list)
	// fmt.Println(sha.GetFileSHA("testdir/a.txt"))
	// // fi, err := os.Stat("testdir/a.txt")
	// // stat := fi.Sys().(*syscall.Win32FileAttributeData)
	// // if err != nil {
	// // 	fmt.Println(err)
	// // 	return
	// // }
	// //
	// // fmt.Println(time.Unix(0, stat.CreationTime.Nanoseconds()))
	// fmt.Println(time.Unix(1267867237, 0))
	// fmt.Println(getCreateTime("testdir/a.txt"))
	// f := File{}
	// fmt.Println(f.Size)
	// a1 := getCreateTime("c:/users/saito0924/go/src/duplicate/testdir/a.txt")
	// a2 := getCreateTime("c:/users/saito0924/go/src/duplicate/testdir/c.txt")
	// fmt.Println(a1, a2)
	// fmt.Println(a1.Before(a2))
	t := time.Now()
	fmt.Println(t)
}
