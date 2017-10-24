// 指定フォルダより重複ファイルを探しリストにする
package main

import (
	"bufio"
	"duplicate/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
			SHA256: utils.GetFileSHA(path),
			Size:   file.Size(),
			Create: utils.GetCreateTime(path),
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

func main() {
	list, err := dirList("testdir")
	if err != nil {
		log.Fatal(err)
	}
	dup := dupList(list)
	tf := 0
	fz := int64(0)
	wz := int64(0)
	for _, d := range dup {
		tf += len(d.Files)
		fz += d.TotalSize
		wz += d.WasteSize
	}

	fmt.Printf("トータル重複ファイル数: %v ファイルサイズ: %vbyte (うち無駄なファイルサイズ: %vbyte)\n", tf, fz, wz)
	fmt.Print("重複ファイルをショートカットに置き換えますか？(y/n)")
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	if stdin.Text() != "y" {
		return
	}
	fmt.Println("yes!")
	err = os.Symlink(`C:\Users\saito0924\go\src\duplicate\testdir\a.txt`, `C:\Users\saito0924\go\src\duplicate\testdir\a.lnk`)
	if err != nil {
		log.Fatal(err)
	}
}
