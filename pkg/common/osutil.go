package common

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//FileExists -- ファイルの存在を確認します。
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

//MakeDir -- パスが存在しなければフォルダを作成する。
func MakeDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
	}
}

//ListUpFiles -- path配下(サブディレクトリ含む)、prefix ではじまり、 suffixで終わるファイル一覧を取得
func ListUpFiles(root, prefix, suffix string) []string {
	var ans []string
	prefix = strings.ToLower(prefix)
	suffix = strings.ToLower(suffix)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !(info.Mode().IsDir()) {
			filenm := info.Name()
			fnm := strings.ToLower(filenm)
			if ((prefix == Empty) || (strings.HasPrefix(fnm, prefix))) &&
				((suffix == Empty) || (strings.HasSuffix(fnm, suffix))) {
				ans = append(ans, filepath.ToSlash(path))
			}
			return nil
		}
		return nil
	})
	return ans
}

//ListUpRece -- path配下(サブディレクトリ含む)、suffixで終わる レセプトファイル一覧を取得
func ListUpRece(root, suffix string) [][]string {
	ans := [][]string{{}, {}, {}, {}, {}}
	suffix = strings.ToUpper(suffix) + ".CSV"
	rMed := regexp.MustCompile(`1[_-]REC\S+[_-]MED`)
	rDpc := regexp.MustCompile(`2[_-]REC\S+[_-]DPC`)
	rDen := regexp.MustCompile(`3[_-]REC\S+[_-]DEN`)
	rPha := regexp.MustCompile(`4[_-]REC\S+[_-]PHA`)
	rNur := regexp.MustCompile(`5[_-]REC\S+[_-]NUR`)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !(info.Mode().IsDir()) {
			fnm := strings.ToUpper(info.Name())
			if strings.HasSuffix(fnm, suffix) {
				if rMed.MatchString(fnm) {
					ans[0] = append(ans[0], filepath.ToSlash(path))
				} else if rDpc.MatchString(fnm) {
					ans[1] = append(ans[1], filepath.ToSlash(path))
				} else if rDen.MatchString(fnm) {
					ans[2] = append(ans[2], filepath.ToSlash(path))
				} else if rPha.MatchString(fnm) {
					ans[3] = append(ans[3], filepath.ToSlash(path))
				} else if rNur.MatchString(fnm) {
					ans[4] = append(ans[4], filepath.ToSlash(path))
				}
			}
		}
		return nil
	})
	return ans
}
