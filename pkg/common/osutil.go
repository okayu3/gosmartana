package common

import "os"

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
