package ana

import (
	"os"
	"strconv"
	"strings"

	"github.com/okayu3/gosmartana/pkg/ckey"
	"github.com/okayu3/gosmartana/pkg/common"
	"github.com/okayu3/gosmartana/pkg/rece"
)

const (
	gakuForTen = 10
)

//MstB -- Bマスタのdictionary
var MstB = make(map[string][]string)

//LoadDisB -- bマスタのロード
func LoadDisB(fnm string, fnmCd119 string) {
	var sybcd, sybnm, icd10, cd119 string
	common.LoadCSV(fnm, func(arr []string, lineno int) {
		sybcd = arr[3-1]
		sybnm = arr[6-1]
		icd10 = arr[16-1]
		MstB[sybcd] = []string{sybnm, icd10, common.Empty}
	}, common.ModeCsvSJIS)
	common.LoadCSV(fnmCd119, func(arr []string, lineno int) {
		sybcd = arr[1-1]
		cd119 = arr[2-1]
		if _, ok := MstB[sybcd]; ok {
			MstB[sybcd][3-1] = cd119
		}
	}, common.ModeCsvSJIS)

}

//MakeSVMed -- 医科レセプトの SV作成
func MakeSVMed(one [][]string, fnm string, args []interface{}) int {
	var arr []string
	var sort, mnKensaku, ymdB, gend, name, kananame string
	var ck, iKi, iBan string
	cntDis := 0
	kananame = common.Empty
	svHandle := args[0].(*os.File)
	for _, arr = range one {
		if (arr == nil) || (len(arr) <= 0) {
			continue
		}
		sort = arr[rece.RSort]
		if sort == "MN" {
			mnKensaku = arr[rece.RMNkensaku]
		} else if sort == "RE" {
			if mnKensaku == common.Empty {
				mnKensaku = arr[rece.RREkensaku]
			}
			ymdB = common.YmdW2g(arr[rece.RREbirth])
			gend = arr[rece.RREgender]
			name = arr[rece.RREname]
			if len(arr) > rece.RREkananame {
				kananame = arr[rece.RREkananame]
			}
		} else if sort == "HO" {
			iKi = arr[rece.RHOinsKigo]
			iBan = arr[rece.RHOinsBango]
			ck = ckey.Get(iKi, iBan, ymdB, gend, name, kananame)
		} else if sort == "SY" {
			cntDis++
			opSaveSV(svHandle, ck, mnKensaku, cntDis, gend, arr)
		}
	}
	return 1
}

func opSaveSV(svHandle *os.File, ck string, mnKensaku string, cntDis int, gend string, arr []string) {
	lineno := arr[2-1]
	recSeq := arr[3-1]
	recDsc := arr[4-1]
	sybcd := arr[5-1]
	innDate := arr[6-1]
	tenki := arr[7-1]
	affix := arr[8-1]
	disnm := arr[9-1]
	flgMain := arr[10-1]

	disInfo, ok := MstB[sybcd]
	if (disnm == common.Empty) && ok {
		disnm = disInfo[0]
	}
	icd10 := disInfo[1]
	cd119 := disInfo[2]
	flgDoubt := "0"
	if common.IsDoubtDisease(affix) {
		flgDoubt = "1"
	}
	//flgHandw := common.Empty
	flgHandw := "0"
	if sybcd == "0000999" {
		flgHandw = "1"
	}
	gaku := common.Empty

	oneSv := strings.Join([]string{ck, mnKensaku, strconv.Itoa(cntDis),
		lineno, recDsc, recSeq, sybcd, affix, affix, cd119, icd10,
		disnm, flgMain, innDate, tenki, flgDoubt, flgHandw, gend, gaku},
		common.Comma)
	svHandle.WriteString(oneSv + "\n")
}
