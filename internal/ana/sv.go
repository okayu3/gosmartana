package ana

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/okayu3/gosmartana/pkg/common"
)

const (
	gakuForTen = 10
)

//MstB -- Bマスタのdictionary
var MstB = make(map[string][]string)

//MstHB -- HBマスタのdictionary
var MstHB = make(map[string][]string)

//MstSTopic -- 診療行為に関連するマスタ
var MstSTopic = make(map[string]map[string]int)
var wkSTopic map[string]int

//LoadDisB -- bマスタのロード
func LoadDisB(fnm string, fnmCd119 string, fnmHb string, fnmToseki string) {
	var sybcd, sybnm, icd10, cd119, sybcdPre, sybcdNxt string
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
	common.LoadCSV(fnmHb, func(arr []string, lineno int) {
		sybcdPre = arr[3-1]
		sybcdNxt = arr[4-1]
		sybnm = arr[6-1]
		MstHB[sybcdPre] = []string{sybcdNxt, sybnm}
	}, common.ModeCsvSJIS)
	wkSTopic = nil
	common.LoadCSV(fnmToseki, func(arr []string, lineno int) {
		if len(arr) == 1 {
			if arr[0] == common.Empty {
				return
			}
			kk := strings.TrimSpace(arr[0])
			if strings.HasPrefix(kk, "[") && strings.HasSuffix(kk, "]") {
				kk = kk[1 : len(kk)-1]
				MstSTopic[kk] = make(map[string]int)
				wkSTopic = MstSTopic[kk]
			}
		}
		if (wkSTopic != nil) && (len(arr) == 2) {
			wkSTopic[arr[0]] = 1
		}
	}, common.ModeCsvSJIS)

}

/* //MakeSVMed -- 医科レセプトの SV作成
func MakeSVMed(one [][]string, fnm string, args []interface{}) int {
	var arr []string
	var sort, mnKensaku, ymdB, gend, name, kananame string
	var ck, iKi, iBan string
	flgHO := 0
	cntDis := 0
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
			flgHO = 1
			iKi = arr[rece.RHOinsKigo]
			iBan = arr[rece.RHOinsBango]
			ck = ckey.Get(iKi, iBan, ymdB, gend, name, kananame)
		} else if sort == "HO" {
			if flgHO == 0 {
				//共通(のちにrefactoringすること)
				iKi = arr[rece.RKOftn]
				iBan = arr[rece.RKOrcv]
				ck = ckey.Get(iKi, iBan, ymdB, gend, name, kananame)
			}
		} else if sort == "SY" {
			cntDis++
			opSaveSV(svHandle, ck, mnKensaku, cntDis, gend, arr)
		}
	}
	return 1
}
*/
func opSaveSV(svHandle *os.File, ck string, mnKensaku string, cntDis int, gend string, arr []string) (string, string) {
	lineno := arr[2-1]
	recSeq := arr[3-1]
	recDsc := arr[4-1]
	sybcd := arr[5-1]
	innDate := arr[6-1]
	tenki := arr[7-1]
	affix := arr[8-1]
	prefix, suffix := common.DivAffix(affix)
	disnm := arr[9-1]
	flgMain := arr[10-1]

	disInfo, ok := MstB[sybcd]
	var icd10, cd119 string
	if ok {
		icd10 = disInfo[1]
		cd119 = disInfo[2]
		if disnm == common.Empty {
			disnm = disInfo[0]
		}
	} else {
		sybcd, icd10, cd119, disnm = detectSybcd(sybcd, disnm)
	}
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
	sCntDis := strconv.Itoa(cntDis)

	oneSv := strings.Join([]string{ck, mnKensaku, sCntDis,
		lineno, recDsc, recSeq, sybcd, prefix, suffix, cd119, icd10,
		disnm, flgMain, innDate, tenki, flgDoubt, flgHandw, gend, gaku},
		common.Comma)
	svHandle.WriteString(oneSv + "\n")
	return sCntDis, cd119
}

func opSaveSVatSB(svHandle *os.File, ck string, mnKensaku string, cntDis int,
	gend string, arr []string, buInnDate, buTenki string) (string, string) {
	lineno := arr[2-1]
	recSeq := arr[3-1]
	recDsc := arr[4-1]
	sybcd := arr[5-1]
	innDate := buInnDate
	tenki := buTenki
	affix := arr[6-1]
	prefix, suffix := common.DivAffix(affix)
	disnm := arr[7-1]
	icd10 := arr[8-1]
	flgMain := common.Empty
	if arr[9-1] == "11" {
		flgMain = "1"
	}

	disInfo, ok := MstB[sybcd]
	var cd119 string
	if ok {
		cd119 = disInfo[2]
		if disnm == common.Empty {
			disnm = disInfo[0]
		}
	} else {
		sybcd, icd10, cd119, disnm = detectSybcd(sybcd, disnm)
	}
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
	sCntDis := strconv.Itoa(cntDis)

	oneSv := strings.Join([]string{ck, mnKensaku, sCntDis,
		lineno, recDsc, recSeq, sybcd, prefix, suffix, cd119, icd10,
		disnm, flgMain, innDate, tenki, flgDoubt, flgHandw, gend, gaku},
		common.Comma)
	svHandle.WriteString(oneSv + "\n")
	return sCntDis, cd119
}

func detectSybcd(sybcd, disnm string) (string, string, string, string) {
	var icd10, cd119 string
	disInfoHB, ok := MstHB[sybcd]
	if !ok {
		fmt.Printf("sybcd:%s sybnm:%s not found yeah\n", sybcd, disnm)
		return sybcd, icd10, cd119, disnm
	}
	sybcdNxt := disInfoHB[0]
	if disnm == common.Empty {
		disnm = disInfoHB[1]
	}
	if sybcdNxt == common.Empty {
		return sybcd, icd10, cd119, disnm
	}
	disInfo, okN := MstB[sybcdNxt]
	if okN {
		icd10 = disInfo[1]
		cd119 = disInfo[2]
	}
	return sybcdNxt, icd10, cd119, disnm
}
