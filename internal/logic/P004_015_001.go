package logic

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/okayu3/gosmartana/pkg/common"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

//dicC4P15V1 // annual -- ck -- bitPerCancer
var dicC4P15V1 = make(map[string]map[string][]bool)

//dicC4P15V1honn // annual -- ck -- honn
var dicC4P15V1honn = make(map[string]map[string]string)

var disGrpNum = 3

func loadingDisPdmC4P15V1(ck, mnKensaku, cd119, icd10, sybcd, flgDoubt string, gaku int) {
	if flgDoubt != "0" && flgDoubt != "" {
		return
	}
	flg := -1
	if strings.HasPrefix(icd10, "F2") {
		flg = 0
	}
	if strings.HasPrefix(icd10, "F3") {
		flg = 1
	}
	if strings.HasPrefix(icd10, "F4") {
		flg = 2
	}
	if flg == -1 {
		return
	}
	//ここでメンタル系疾患であることが判明。

	da, okA := DicExp[mnKensaku]
	if !okA {
		return
	}
	honn := da[2-1]
	seikyuYm := da[5-1]
	ann := calcReceAnnual(seikyuYm)

	if _, ok := dicC4P15V1[ann]; !ok {
		dicC4P15V1[ann] = make(map[string][]bool)
		dicC4P15V1honn[ann] = make(map[string]string)
	}
	if _, ok := dicC4P15V1[ann][ck]; !ok {
		dicC4P15V1[ann][ck] = make([]bool, disGrpNum)
	}
	dicC4P15V1[ann][ck][flg] = true
	dicC4P15V1honn[ann][ck] = honn
}

func opSummaryC4P15V1(logicOutdir string) {
	for ann := range dicC4P15V1 {
		opSummaryC4P15V1Main(ann, logicOutdir)
	}
}

func opSummaryC4P15V1Main(ann, logicOutdir string) {
	ofnm := logicOutdir + "Result_C4P15V1_" + ann + ".csv"
	oHandle, _ := os.Create(ofnm)
	defer oHandle.Close()
	writer := bufio.NewWriter(transform.NewWriter(oHandle, japanese.ShiftJIS.NewEncoder()))

	var prefix string
	sumDic := make(map[string]map[int][]int)
	ageAtYm := ann + "04"
	for ck, mp := range dicC4P15V1[ann] {
		honn, ok := dicC4P15V1honn[ann][ck]
		if !ok {
			continue
		}
		db, okA := DicPsn[ck]
		if !okA {
			continue
		}
		gend := db[1-1]
		ymdB := db[2-1]
		ageRange := calcAgeRange(ymdB, ageAtYm)
		sort := db[3-1]

		if sort == "0" {
			prefix = honn
		} else {
			prefix = "tn"
		}
		kk := []string{"0", gend, prefix + ":0", prefix + ":" + gend}
		for _, k := range kk {
			if _, ok := sumDic[k]; !ok {
				sumDic[k] = make(map[int][]int)
			}
			if _, ok := sumDic[k][ageRange]; !ok {
				sumDic[k][ageRange] = make([]int, disGrpNum+1)
			}
			if _, ok := sumDic[k][-5]; !ok {
				sumDic[k][-5] = make([]int, disGrpNum+1)
			}
			for i := 0; i < disGrpNum; i++ {
				if mp[i] {
					sumDic[k][ageRange][i]++
					sumDic[k][-5][i]++
				}
			}
			sumDic[k][ageRange][disGrpNum]++
			sumDic[k][-5][disGrpNum]++
		}

	}
	pgDsc := []string{"全体", "男性", "女性", "一般本人(全体)",
		"一般本人(男性)", "一般本人(女性)", "一般家族(全体)", "一般家族(男性)",
		"一般家族(女性)", "特退任継(全体)", "特退任継(男性)", "特退任継(女性)"}
	pgFlg := []string{"0", "1", "2", "1:0", "1:1", "1:2", "2:0", "2:1", "2:2", "tn:0", "tn:1", "tn:2"}

	disDsc := make([]string, disGrpNum+1)
	disDsc[0] = "全体,"
	disDsc[1] = "F20-F29,"
	disDsc[2] = "F30-F39,"
	disDsc[3] = "F40-F48,"

	for idx, pg := range pgFlg {
		if _, ok := sumDic[pg]; !ok {
			continue
		}
		writer.WriteString("@" + pgDsc[idx] + common.CrLf)
		wk0 := []string{"加入者数,"}
		for i := 0; i <= 15; i++ {
			wk0 = append(wk0, strconv.Itoa(DicPop[pg][(i-1)*5]))
		}
		writer.WriteString(strings.Join(wk0, common.Comma) + common.CrLf)

		for didx, ddsc := range disDsc {
			wk := []string{ddsc}
			at := didx - 1
			if didx == 0 {
				at = disGrpNum
			}
			for i := -1; i <= 14; i++ {
				ageRange := i * 5
				if _, ok := sumDic[pg][ageRange]; !ok {
					wk = append(wk, "0")
				} else {
					wk = append(wk, strconv.Itoa(sumDic[pg][ageRange][at]))
				}
			}
			writer.WriteString(strings.Join(wk, common.Comma) + common.CrLf)
		}
	}
	writer.Flush()
}
