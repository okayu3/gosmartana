package ckey

import (
	"fmt"
	"strings"

	"github.com/okayu3/gosmartana/pkg/common"
)

//PsnMap -- 組合員のCkeyマップ. key は 検索パターン
var PsnMap = make(map[string]string)

//PsnInfoMap -- 組合員の情報マップ
var PsnInfoMap = make(map[string][]string)

//Get 結束キーの取得
func Get(iKi, iBan, ymdB, gend, name, kananame string) string {
	keys := makeKeys(iKi, iBan, ymdB, gend, name, kananame)
	for _, one := range keys {
		if _, ok := PsnMap[one]; ok {
			return PsnMap[one]
		}
	}
	seq := len(PsnMap) + 1
	neuCkey := "M01" + fmt.Sprintf("%010d", seq)
	for _, one := range keys {
		PsnMap[one] = neuCkey
	}
	PsnInfoMap[neuCkey] = []string{neuCkey, iKi, iBan, ymdB, gend, name, kananame}
	return neuCkey
}

func makeKeys(iKi, iBan, ymdB, gend, name, kananame string) []string {
	return []string{
		strings.Join([]string{iKi, iBan, ymdB, gend, name}, common.Collon),
		strings.Join([]string{iKi, iBan, ymdB, gend, kananame}, common.Collon),
		strings.Join([]string{iBan, ymdB, gend, name}, common.Collon),
		strings.Join([]string{iBan, ymdB, gend, kananame}, common.Collon),
		strings.Join([]string{iBan, ymdB, name}, common.Collon),
		strings.Join([]string{iBan, ymdB, kananame}, common.Collon),
		strings.Join([]string{iBan, ymdB, gend}, common.Collon),
	}
}
