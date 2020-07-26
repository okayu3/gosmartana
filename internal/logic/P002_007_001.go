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

//dicC2P7V1 // annual -- kbn -- agerange -- population
var dicC2P7V1 = make(map[string]map[string]map[int]int)
var dicC2P7V1wk = make(map[string]map[string]int)

func loadingTosekiTopicC2P7V1(ck, mnKensaku,
	flgToseki, flgInsulin, flgMngDiabetes,
	flgMngBP, flgMngFat, flgSmoking,
	flgYoboToseki, flgTestHbA1c, flgTestFat string) {

	if flgToseki == "0" {
		return
	}
	da, okA := DicExp[mnKensaku]
	db, okB := DicPsn[ck]
	if !okA || !okB {
		return
	}
	honn := da[2-1]
	sinryoYm := da[3-1]
	seikyuYm := da[5-1]
	gend := db[1-1]
	ymdB := db[2-1]
	ann := calcReceAnnual(seikyuYm)
	ageRange := calcAgeRange(ymdB, strconv.Itoa(common.AnnualAtYm(sinryoYm))+"04")

	if _, ok0 := dicC2P7V1wk[ann][ck]; ok0 {
		return
	}
	if _, ok1 := dicC2P7V1wk[ann]; !ok1 {
		dicC2P7V1wk[ann] = make(map[string]int)
	}
	dicC2P7V1wk[ann][ck] = 1

	if _, ok2 := dicC2P7V1[ann]; !ok2 {
		dicC2P7V1[ann] = make(map[string]map[int]int)
	}
	d := dicC2P7V1[ann]

	kk := []string{"0", gend, honn + ":0", honn + ":" + gend}
	for _, k := range kk {
		if _, ok := d[k]; !ok {
			d[k] = make(map[int]int)
		}
		if _, ok := d[k][ageRange]; !ok {
			d[k][ageRange] = 0
		}
		if _, ok := d[k][-5]; !ok {
			d[k][-5] = 0
		}
		d[k][ageRange]++
		d[k][-5]++
	}

}

func opSummaryC2P7V1(logicOutdir string) {
	for ann := range dicC2P7V1 {
		opSummaryC2P7V1Main(ann, logicOutdir)
	}
}

func opSummaryC2P7V1Main(ann, logicOutdir string) {
	ofnm := logicOutdir + "Result_C2P7V1_" + ann + ".csv"
	oHandle, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_CREATE, 0666)
	defer oHandle.Close()
	writer := bufio.NewWriter(transform.NewWriter(oHandle, japanese.ShiftJIS.NewEncoder()))

	pgDsc := []string{"全体", "男性", "女性", "本人(全体)",
		"本人(男性)", "本人(女性)", "本人外(全体)", "本人外(男性)",
		"本人外(女性)"}
	pgFlg := []string{"0", "1", "2", "1:0", "1:1", "1:2", "2:0", "2:1", "2:2"}
	title := "年齢層,人数,透析人数"

	for idx, pg := range pgFlg {
		writer.WriteString("@" + pgDsc[idx] + common.CrLf)
		writer.WriteString(title + common.CrLf)
		for i := -1; i <= 14; i++ {
			ageRange := i * 5
			pop, ok0 := DicPop[pg][ageRange]
			if !ok0 {
				pop = 0
			}
			v1, ok1 := dicC2P7V1[ann][pg][ageRange]
			if !ok1 {
				v1 = 0
			}
			wk := []int{ageRange, pop, v1}
			wk2 := []string{}
			for _, one := range wk {
				wk2 = append(wk2, strconv.Itoa(one))
			}
			writer.WriteString(strings.Join(wk2, common.Comma) + common.CrLf)
		}
	}
	writer.Flush()
}
