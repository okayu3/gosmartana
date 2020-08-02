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

//dicC3P10V1 // annual -- ck -- bitPerCancer
var dicC3P10V1 = make(map[string]map[string][]bool)

//dicC3P10V1honn // annual -- ck -- honn
var dicC3P10V1honn = make(map[string]map[string]string)

func loadingDisPdmC3P10V1(ck, mnKensaku, cd119, icd10, sybcd, flgDoubt string, gaku int) {
	if flgDoubt != "0" && flgDoubt != "" {
		return
	}
	if !strings.HasPrefix(icd10, "C") {
		return
	}
	//ここで癌であることが判明。

	da, okA := DicExp[mnKensaku]
	if !okA {
		return
	}
	honn := da[2-1]
	seikyuYm := da[5-1]
	ann := calcReceAnnual(seikyuYm)

	if _, ok := dicC3P10V1[ann]; !ok {
		dicC3P10V1[ann] = make(map[string][]bool)
		dicC3P10V1honn[ann] = make(map[string]string)
	}
	if _, ok := dicC3P10V1[ann][ck]; !ok {
		dicC3P10V1[ann][ck] = make([]bool, CstCancerNum)
	}

	if idx, ok := MstCancer[cd119]; ok {
		dicC3P10V1[ann][ck][idx] = true
	}
	if cd119 == "0210" {
		if idx, ok := MstCancer[icd10[0:3]]; ok {
			dicC3P10V1[ann][ck][idx] = true
		}
	}
	dicC3P10V1honn[ann][ck] = honn
}

func opSummaryC3P10V1(logicOutdir string) {
	for ann := range dicC3P10V1 {
		opSummaryC3P10V1Main(ann, logicOutdir)
	}
}

func opSummaryC3P10V1Main(ann, logicOutdir string) {
	ofnm := logicOutdir + "Result_C3P10V1_" + ann + ".csv"
	oHandle, _ := os.Create(ofnm)
	defer oHandle.Close()
	writer := bufio.NewWriter(transform.NewWriter(oHandle, japanese.ShiftJIS.NewEncoder()))

	var prefix string
	sumDic := make(map[string]map[int][]int)
	ageAtYm := ann + "04"
	for ck, mp := range dicC3P10V1[ann] {
		honn, ok := dicC3P10V1honn[ann][ck]
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
				sumDic[k][ageRange] = make([]int, CstCancerNum+1)
			}
			if _, ok := sumDic[k][-5]; !ok {
				sumDic[k][-5] = make([]int, CstCancerNum+1)
			}
			for i := 0; i < CstCancerNum; i++ {
				if mp[i] {
					sumDic[k][ageRange][i]++
					sumDic[k][-5][i]++
				}
			}
			sumDic[k][ageRange][CstCancerNum]++
			sumDic[k][-5][CstCancerNum]++
		}

	}
	pgDsc := []string{"全体", "男性", "女性", "一般本人(全体)",
		"一般本人(男性)", "一般本人(女性)", "一般家族(全体)", "一般家族(男性)",
		"一般家族(女性)", "特退任継(全体)", "特退任継(男性)", "特退任継(女性)"}
	pgFlg := []string{"0", "1", "2", "1:0", "1:1", "1:2", "2:0", "2:1", "2:2", "tn:0", "tn:1", "tn:2"}

	disDsc := make([]string, CstCancerNum+1)
	disDsc[0] = "全体,"
	for k, v := range MstCancer {
		if _, ok := MstCancerDsc[k]; !ok {
			continue
		}
		disDsc[v+1] = k + "," + MstCancerDsc[k]
	}

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
				at = CstCancerNum
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
