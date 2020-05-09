package common

import "os"

//FileExists -- ファイルの存在を確認します。
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
