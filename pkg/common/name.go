package common

import "strings"

//DevideName -- 氏名 を 最初の全角スペースで区切って 姓と名に分ける
func DevideName(nm string) (string, string) {
	deli := "　"
	idx := strings.Index(nm, deli)
	if idx < 0 {
		return nm, Empty
	} else if idx == 0 {
		return nm[idx+len(deli):], Empty
	}
	return nm[0:idx], nm[idx+len(deli):]
}
