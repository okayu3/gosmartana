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

//dicC1P2V1 // annual // kbn // agerange // cd19 // gaku
var dicC1P2V1 = make(map[string]map[string]map[int]map[string]int)

func loadingExpenseC1P2V1(ck, mnKensaku, sort, nyugai, honn, sinryoYm, jitsuDates, seikyuYm string, gaku int) {
	if sort != "歯科" {
		return
	}
	loadingDisPdmC1P2V1(ck, mnKensaku, "8099", gaku)
}

func loadingDisPdmC1P2V1(ck, mnKensaku, cd119 string, gaku int) {
	da, okA := DicExp[mnKensaku]
	db, okB := DicPsn[ck]
	if !okA || !okB {
		return
	}
	honn := da[2-1]
	sinryoYm := da[3-1]
	sort := da[4-1]
	seikyuYm := da[5-1]
	ann := calcReceAnnual(seikyuYm)
	if _, ok2 := dicC1P2V1[ann]; !ok2 {
		dicC1P2V1[ann] = make(map[string]map[int]map[string]int)
	}
	d := dicC1P2V1[ann]

	gend := db[1-1]
	ymdB := db[2-1]
	ageRange := calcAgeRange(ymdB, sinryoYm)
	cd19 := "99"
	if cd119 != common.Empty {
		cd19 = cd119[0:2]
		if cd19 == "16" { //周産期に発生した病態(P00-P96) を、妊娠・出産にまとめる。
			cd19 = "15"
		} else if cd19 == "20" { //19と20を外因系でまとめる
			cd19 = "19"
		} else if (cd19 == "18") || (cd19 == "21") || (cd19 == "22") {
			//18, 21, 22 を「99:その他」でまとめる
			cd19 = "99"
		}
	}
	if sort == "歯科" { //ここにはこない。
		cd19 = "80"
	}
	kk := []string{"0", gend, honn + ":0", honn + ":" + gend}
	for _, k := range kk {
		if _, ok := d[k]; !ok {
			d[k] = make(map[int]map[string]int)
		}
		if _, ok := d[k][ageRange]; !ok {
			d[k][ageRange] = make(map[string]int)
		}
		if _, ok := d[k][-5]; !ok {
			d[k][-5] = make(map[string]int)
		}
		d[k][ageRange][cd19] += gaku
		d[k][-5][cd19] += gaku
	}
}

func opSummaryC1P2V1(logicOutdir string) {
	for ann := range dicC1P2V1 {
		opSummaryC1P2V1Main(ann, logicOutdir)
	}
}

func opSummaryC1P2V1Main(ann, logicOutdir string) {
	ofnm := logicOutdir + "Result_C1P2V1_" + ann + ".csv"
	oHandle, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_CREATE, 0666)
	defer oHandle.Close()
	writer := bufio.NewWriter(transform.NewWriter(oHandle, japanese.ShiftJIS.NewEncoder()))

	pgDsc := []string{"全体", "男性", "女性", "本人(全体)",
		"本人(男性)", "本人(女性)", "本人外(全体)", "本人外(男性)",
		"本人外(女性)"}
	pgFlg := []string{"0", "1", "2", "1:0", "1:1", "1:2", "2:0", "2:1", "2:2"}
	dsDsc := []string{
		"感染症", "新生物", "血液・免疫系", "内分泌・代謝",
		"精神系", "神経系", "眼系", "耳系",
		"循環器系", "呼吸器系", "消化器系", "皮膚系",
		"筋骨格系", "尿路・性器", "妊娠・出産", "先天性・染色体系",
		"その他", "外因系", "歯科"}
	dsFlg := []string{"01", "02", "03", "04", "05", "06", "07", "08",
		"09", "10", "11", "12", "13", "14", "15", "17", "99", "19", "80"}

	for idx, pg := range pgFlg {
		writer.WriteString("@" + pgDsc[idx] + common.CrLf)
		writer.WriteString("年齢," + strings.Join(dsDsc, common.Comma) + common.CrLf)
		for i := -1; i <= 14; i++ {
			ageRange := i * 5
			r1, ok1 := dicC1P2V1[ann][pg][ageRange]
			wk := []string{strconv.Itoa(ageRange)}
			if !ok1 {
				for range dsFlg {
					wk = append(wk, "0")
				}
			} else {
				for _, ds := range dsFlg {
					gaku, ok2 := r1[ds]
					if !ok2 {
						gaku = 0
					}
					wk = append(wk, strconv.Itoa(gaku))
				}
			}
			writer.WriteString(strings.Join(wk, common.Comma) + common.CrLf)
		}
	}
	writer.Flush()
}
