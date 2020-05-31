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

//dicC1P1V1
var dicC1P1V1 = make(map[string][]map[string]int)

func loadingExpenseC1P1V1(ck, mnKensaku, sort, nyugai, honn, sinryoYm, jitsuDates string, gaku int) {
	db, okA := DicPsn[ck]
	if !okA {
		return
	}
	gend := db[1-1]
	kk := []string{"0", gend, honn + ":0", honn + ":" + gend}
	ss := "99"
	if sort == "医科" {
		if nyugai == "2" { //外来
			ss = "1:2"
		} else {
			ss = "1:1"
		}
	} else if sort == "ＤＰＣ" {
		ss = "1:3"
	} else if sort == "歯科" {
		ss = "2"
	} else if sort == "調剤" {
		ss = "3"
	}
	for _, k := range kk {
		if _, ok := dicC1P1V1[k]; !ok {
			dicC1P1V1[k] = append(dicC1P1V1[k],
				make(map[string]int), make(map[string]int),
				make(map[string]int), make(map[string]int))
		}
		dicC1P1V1[k][0][ss]++
		dicC1P1V1[k][1][ss] += common.Atoi(jitsuDates, 0)
		dicC1P1V1[k][2][ss] += gaku
	}
}

func opSummaryC1P1V1(logicOutdir string) {
	ofnm := logicOutdir + "Result_C1P1V1.csv"
	oHandle, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_CREATE, 0666)
	defer oHandle.Close()
	writer := bufio.NewWriter(transform.NewWriter(oHandle, japanese.ShiftJIS.NewEncoder()))
	as := []string{"1:2", "1:1", "1:3", "2", "3"}

	pgDsc := []string{"全体", "男性", "女性", "本人(全体)",
		"本人(男性)", "本人(女性)", "本人外(全体)", "本人外(男性)",
		"本人外(女性)"}
	pgFlg := []string{"0", "1", "2", "1:0", "1:1", "1:2", "2:0", "2:1", "2:2"}
	title := "種類,医科外来,医科入院,医科ＤＰＣ,歯科,調剤"
	srtDsc := []string{"件数", "診療日数", "決定金額"}

	for idx, pg := range pgFlg {
		writer.WriteString("@" + pgDsc[idx] + common.CrLf)
		writer.WriteString(title + common.CrLf)
		for i := 0; i < 3; i++ {
			wk := []string{srtDsc[i]}
			for _, aa := range as {
				d, ok := dicC1P1V1[pg][i][aa]
				if !ok {
					wk = append(wk, "0")
				} else {
					wk = append(wk, strconv.Itoa(d))
				}
			}
			writer.WriteString(strings.Join(wk, common.Comma) + common.CrLf)
		}
	}
	writer.Flush()
}
