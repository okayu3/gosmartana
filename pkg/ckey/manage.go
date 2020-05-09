package ckey

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/okayu3/gosmartana/pkg/common"
)

//PsnMap -- 組合員のCkeyマップ. key は 検索パターン
var PsnMap = make(map[string]string)

//PsnInfoMap -- 組合員の情報マップ
var PsnInfoMap = make(map[string][]string)

//psnNewly -- 今回発見された被保険者のマップ
var psnNewly = make(map[string]int)

//fnmPsnMst -- PsnMst のファイル名
var fnmPsnMst string

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
	psnNewly[neuCkey] = 1
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

//LoadPersonMst -- Personファイルの読み込み
func LoadPersonMst(fnm string) {
	var ck, iKi, iBan, ymdB, gend, nm, knm string
	var keys []string
	//var sei, mei, ksei, kmei, iBranch string
	fnmPsnMst = fnm

	if !common.FileExists(fnm) {
		return
	}
	common.LoadCSV(fnmPsnMst, func(a []string, lineno int) {
		ck = a[0]
		iKi = a[1]
		iBan = a[2]
		//iBranch = a[3]
		gend = a[4]
		ymdB = a[5]
		nm = a[6]
		knm = a[7]
		//sei = a[8]
		//mei = a[9]
		//ksei = a[10]
		//kmei = a[11]
		keys = makeKeys(iKi, iBan, ymdB, gend, nm, knm)
		for _, one := range keys {
			PsnMap[one] = ck
		}
		PsnInfoMap[ck] = []string{ck, iKi, iBan, ymdB, gend, nm, knm}
	}, common.ModeCsvUTF8)
}

//UpdatePersonMst -- 新規で発見された人の情報だけ保存
func UpdatePersonMst() {
	var ck, iKi, iBan, ymdB, gend, nm, knm string
	var sei, mei, ksei, kmei string
	outHandlePsn, _ := os.OpenFile(fnmPsnMst, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer outHandlePsn.Close()

	neuCks := make([]string, 0, len(psnNewly))
	for k := range psnNewly {
		neuCks = append(neuCks, k)
	}
	sort.Strings(neuCks)
	for _, neuCkey := range neuCks {
		a := PsnInfoMap[neuCkey]
		ck = a[0]
		iKi = a[1]
		iBan = a[2]
		ymdB = a[3]
		gend = a[4]
		nm = a[5]
		knm = a[6]
		sei, mei = common.DevideName(nm)
		ksei, kmei = common.DevideName(knm)
		iBranch := common.Empty
		onePsn := strings.Join([]string{ck, iKi, iBan, iBranch,
			gend, ymdB,
			nm, knm, sei, mei, ksei, kmei,

			common.Empty}, common.Comma)
		outHandlePsn.WriteString(onePsn + "\n")
	}
}
