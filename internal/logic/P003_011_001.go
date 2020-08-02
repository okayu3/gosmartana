package logic

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/okayu3/gosmartana/pkg/common"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

//dicC3P11V1 // annual -- ck -- bitPerCancer
var dicC3P11V1 = make(map[string]map[string][]int)

//dicC3P11V1honn // annual -- ck -- honn
var dicC3P11V1honn = make(map[string]map[string]string)

func loadingDisPdmC3P11V1(ck, mnKensaku, cd119, icd10, sybcd, flgDoubt string, gaku int) {
	//疑いもカウントする。
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

	if _, ok := dicC3P11V1[ann]; !ok {
		dicC3P11V1[ann] = make(map[string][]int)
		dicC3P11V1honn[ann] = make(map[string]string)
	}
	if _, ok := dicC3P11V1[ann][ck]; !ok {
		dicC3P11V1[ann][ck] = make([]int, CstCancerNum)
	}

	if idx, ok := MstCancer[cd119]; ok {
		dicC3P11V1[ann][ck][idx] += gaku
	}
	if cd119 == "0210" {
		if idx, ok := MstCancer[icd10[0:3]]; ok {
			dicC3P11V1[ann][ck][idx] += gaku
		}
	}
	dicC3P11V1honn[ann][ck] = honn
}

func opSummaryC3P11V1(logicOutdir string) {
	for ann := range dicC3P11V1 {
		opSummaryC3P11V1Main(ann, logicOutdir)
	}
}

func opSummaryC3P11V1Main(ann, logicOutdir string) {
	ofnm := logicOutdir + "Result_C3P11V1_" + ann + ".csv"
	oHandle, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_CREATE, 0666)
	defer oHandle.Close()
	writer := bufio.NewWriter(transform.NewWriter(oHandle, japanese.ShiftJIS.NewEncoder()))

	var prefix string
	sumDic := make(map[string]map[int][]int)
	ageAtYm := ann + "04"
	for ck, mp := range dicC3P11V1[ann] {
		honn, ok := dicC3P11V1honn[ann][ck]
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
				sumDic[k][ageRange] = make([]int, CstCancerNum)
			}
			if _, ok := sumDic[k][-5]; !ok {
				sumDic[k][-5] = make([]int, CstCancerNum)
			}
			for i := 0; i < CstCancerNum; i++ {
				sumDic[k][ageRange][i] += mp[i]
				sumDic[k][-5][i] += mp[i]
				//sumDic[k][ageRange][CstCancerNum] += mp[i]
				//sumDic[k][-5][CstCancerNum] += mp[i]
			}
		}

	}
	pgDsc := []string{"全体", "男性", "女性", "一般本人(全体)",
		"一般本人(男性)", "一般本人(女性)", "一般家族(全体)", "一般家族(男性)",
		"一般家族(女性)", "特退任継(全体)", "特退任継(男性)", "特退任継(女性)"}
	pgFlg := []string{"0", "1", "2", "1:0", "1:1", "1:2", "2:0", "2:1", "2:2", "tn:0", "tn:1", "tn:2"}

	disDsc := make([]string, CstCancerNum)
	//disDsc := make([]string, CstCancerNum+1)
	//disDsc[0] = "全体,"
	for k, v := range MstCancer {
		if _, ok := MstCancerDsc[k]; !ok {
			continue
		}
		disDsc[v] = k + "," + MstCancerDsc[k]
	}
	skipDidx := MstCancer["0210"]
	for idx, pg := range pgFlg {
		if _, ok := sumDic[pg]; !ok {
			continue
		}
		writer.WriteString("@" + pgDsc[idx] + common.CrLf)
		for i := -1; i <= 14; i++ {
			ageRange := i * 5
			darr, ok := sumDic[pg][ageRange]
			if !ok {
				continue
			}
			a := List{}
			for didx := range disDsc {
				if didx == skipDidx { //その他悪性新生物(0210)は飛ばしたい
					continue
				}
				if darr[didx] == 0 { //額が0円の場合も飛ばす。
					continue
				}
				e := Entry{strconv.Itoa(didx), darr[didx]}
				a = append(a, e)
			}
			wk := []string{strconv.Itoa(ageRange)}
			sort.Sort(sort.Reverse(a))
			for rank, ent := range a {
				if rank > 9 {
					break
				}
				didx := common.Atoi(ent.name, -1)
				if didx < 0 {
					continue
				}
				wk = append(wk, disDsc[didx], strconv.Itoa(darr[didx]))
			}
			writer.WriteString(strings.Join(wk, common.Comma) + common.CrLf)
		}
	}

	writer.Flush()
}

// for idx, pg := range pgFlg {
// 	if _, ok := sumDic[pg]; !ok {
// 		continue
// 	}
// 	writer.WriteString("@" + pgDsc[idx] + common.CrLf)
// 	wk0 := []string{"加入者数,"}
// 	for i := 0; i <= 15; i++ {
// 		wk0 = append(wk0, strconv.Itoa(DicPop[pg][(i-1)*5]))
// 	}
// 	writer.WriteString(strings.Join(wk0, common.Comma) + common.CrLf)

// 	for didx, ddsc := range disDsc {
// 		wk := []string{ddsc}
// 		at := didx - 1
// 		if didx == 0 {
// 			at = CstCancerNum
// 		}
// 		for i := -1; i <= 14; i++ {
// 			ageRange := i * 5
// 			if _, ok := sumDic[pg][ageRange]; !ok {
// 				wk = append(wk, "0")
// 			} else {
// 				wk = append(wk, strconv.Itoa(sumDic[pg][ageRange][at]))
// 			}
// 		}
// 		writer.WriteString(strings.Join(wk, common.Comma) + common.CrLf)
// 	}
// }
