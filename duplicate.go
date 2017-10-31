// 指定フォルダより重複ファイルを探しリストにする
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/insight1111/duplicate/utils"

	"github.com/ktat/go-pager"
	"gopkg.in/cheggaaa/pb.v1"
)

// var SHA chan string

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
	// pwd, _ := os.Getwd()
	if err != nil {
		return
	}
	count := len(_result)
	bar := pb.StartNew(count)
	files := []File{}
	for _, file := range _result {
		bar.Increment()
		fi, err := os.Lstat(filepath.Join(startDir, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		if file.IsDir() {
			_t, _ := dirList2(filepath.Join(startDir, file.Name()))
			files = append(files, _t...)
			continue
		} else if fi.Mode()&os.ModeSymlink != 0 {
			continue
		} else if file.Name()[0:2] == "~$" {
			continue
		} else if filepath.Ext(file.Name()) == ".lnk" {
			continue
		}
		// path := filepath.Join(pwd, startDir, file.Name())
		path := filepath.Join(startDir, file.Name())
		// sha:=make(chan string)
		// go utils.GetCreateTime(path)
		// s:=<-sha
		f := File{
			Path:   path,
			SHA256: utils.GetFileSHA(path),
			Size:   file.Size(),
			Create: utils.GetCreateTime(path),
		}
		files = append(files, f)
		time.Sleep(time.Millisecond)
	}
	bar.FinishPrint("調査完了しました...")
	result = files
	return
}
func dirList2(startDir string) (result []File, err error) {
	_result, err := ioutil.ReadDir(startDir)
	// pwd, _ := os.Getwd()
	if err != nil {
		return
	}
	files := []File{}
	for _, file := range _result {
		fi, err := os.Lstat(filepath.Join(startDir, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		if file.IsDir() {
			_t, _ := dirList2(filepath.Join(startDir, file.Name()))
			files = append(files, _t...)
			continue
		} else if fi.Mode()&os.ModeSymlink != 0 {
			continue
		} else if file.Name()[0:2] == "~$" {
			continue
		} else if filepath.Ext(file.Name()) == ".lnk" {
			continue
		}
		// path := filepath.Join(pwd, startDir, file.Name())
		path := filepath.Join(startDir, file.Name())
		sha := make(chan string)
		go utils.GetFileSHA2(path, sha)
		s := <-sha

		f := File{
			Path:   path,
			SHA256: s,
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

func replaceSymlink(dups []DupFiles) {
	for _, d := range dups {
		o := d.Original
		dir, filename := filepath.Split(o)
		orig := regexp.MustCompile(`(original)`)
		var ren string
		if !orig.MatchString(filename) {
			ext := filepath.Ext(filename)
			bname := filename[:len(filename)-len(ext)]
			ren = bname + "(original)" + ext
		} else {
			ren = filename
		}
		renpath := filepath.Join(dir, ren)
		os.Rename(o, renpath)
		for _, f := range d.Files {
			if f.Path != o {
				// ext := filepath.Ext(f.Path)
				// bname := f.Path[:len(f.Path)-len(ext)]
				dir, bname := filepath.Split(f.Path)
				// fmt.Print(bname)
				r := filepath.Join(dir, bname+".lnk")
				os.Remove(r)
				err := os.Symlink(renpath, r)
				if err != nil {
					log.Fatal(err)
				}
				os.Remove(f.Path)
			}
		}
	}
}

func main() {
	// fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	// 	flag.Usage = func() {
	// 		fmt.Fprintf(os.Stderr, `
	// Usage of %s:
	//   %s [OPTIONS] AGRS...
	// Options\n`, os.Args[0], os.Args[0])
	// 		flag.PrintDefaults()
	// 	}
	var (
		filePath = flag.String("p", ".", "directory path")
		silent   = flag.Bool("y", false, "no confirm, always select yes.")
	)
	// filePath := os.Args[1]
	flag.Parse()
	// fmt.Printf("filepath %#v y-option %#v", filePath, *silent)
	// pwd, _ := os.Getwd()
	list, err := dirList(*filePath)
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
	fmt.Printf("ファイル一覧を表示しますか？(y/n)")
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	if stdin.Text() == "y" {
		var p pager.Pager
		p.Init()
		r := []string{}
		for _, d := range dup {
			r = append(r, d.Original+"("+strconv.Itoa(len(d.Files))+")")
		}
		p.SetContent(strings.Join(r, "\n"))
		if p.PollEvent() == false {
			p.Close()
		}

	}
	if !*silent {
		fmt.Print("重複ファイルをショートカットに置き換えますか？(y/n)")
		stdin = bufio.NewScanner(os.Stdin)
		stdin.Scan()
		if stdin.Text() != "y" {
			return
		}
	}
	fmt.Println("処理を開始します")
	replaceSymlink(dup)
	fmt.Println("完了しました")
}
