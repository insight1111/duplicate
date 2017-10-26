package main

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

func TestDirList(t *testing.T) {
	pwd, _ := os.Getwd()
	result, _ := dirList(filepath.Join(pwd, "testdir"))
	dir, _ := os.Getwd()

	if reflect.TypeOf(result).Kind() != reflect.Slice {
		t.Errorf("結果は構造体のスライスではありません。実際は %sです", reflect.TypeOf(result))
	}

	// ファイル数は5(再帰調査)
	expect := 5
	if len(result) != expect {
		t.Errorf("受け取ったファイル数が違います　get:%#v expect:%#v", len(result), expect)
	}

	// 結果にはd.txtを含む
	expect = 0
	r := regexp.MustCompile(`d\.txt`)
	for _, file := range result {
		if r.MatchString(file.Path) {
			expect = 1
		}
	}
	if expect == 0 {
		t.Errorf("結果にd.txtが含まれません")
	}

	// 結果はフルパスで入る
	expect = 0
	for _, file := range result {
		if file.Path == filepath.Join(dir, "testdir", "a.txt") {
			expect = 1
		}
	}
	if expect == 0 {
		t.Errorf("結果がフルパスではありません")
	}

	// 結果にはSHA256が含まれている
	expect = 0
	for _, file := range result {
		if file.SHA256 != "" {
			expect = 1
		}
	}
	if expect == 0 {
		t.Errorf("結果にSHA256が含まれていません")
	}

	// a.txtのSHA256が合致している
	expect = 0
	aSHA := "8E4621379786EF42A4FEC155CD525C291DD7DB3C1FDE3478522F4F61C03FD1BD"
	rSHA := ""
	for _, file := range result {
		r = regexp.MustCompile(`a\.txt`)
		if r.MatchString(file.Path) {
			rSHA = file.SHA256
			if rSHA == aSHA {
				expect = 1
			}
		}
	}
	if expect == 0 {
		t.Errorf("a.txtのSHA256がマッチしていません r:%v e:%v", rSHA, aSHA)
	}

}

func TestDupList(t *testing.T) {
	result, _ := dirList("testdir")
	dup := dupList(result)

	if len(dup) == 0 {
		t.Errorf("dupに何も入っていません")
	}

	// 重複ファイルはきちんと結果に入っているか
	expect := 0
	r1 := regexp.MustCompile(`a\.txt`)
	r2 := regexp.MustCompile(`d\.txt`)
	for _, file := range dup[0].Files {
		if r1.MatchString(file.Path) {
			expect++
		} else if r2.MatchString(file.Path) {
			expect++
		}
	}
	if expect < 2 {
		t.Errorf("a.txtまたはd.txtが重複結果に入っていません get %#v", dup[0])
	}

	// オリジナルはa.txt
	if r1.MatchString(dup[0].Original) == false {
		t.Errorf("オリジナルはa.txtではありません get %#v", dup[0].Original)
	}
	// fmt.Printf("%+v\n", dup[0])

}

func BenchmarkDirList(b *testing.B) {
	b.ResetTimer()
	_, _ = dirList("c:/downloads")
}
func BenchmarkDirListChannel(b *testing.B) {
	b.ResetTimer()
	_, _ = dirList2("c:/downloads")
}
