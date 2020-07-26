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

//dicC1P3V1 // annual -- kbn -- agerange -- [gakuidx] gaku
var dicC1P3V1 = make(map[string]map[string]map[int][]int)

func loadingExpenseC1P3V1(ck, mnKensaku, sort, nyugai, honn, sinryoYm, jitsuDates, seikyuYm string, gaku int) {
	db, okA := DicPsn[ck]
	if !okA {
		return
	}
	ann := calcReceAnnual(seikyuYm)
	if _, ok2 := dicC1P3V1[ann]; !ok2 {
		dicC1P3V1[ann] = make(map[string]map[int][]int)
	}
	d := dicC1P3V1[ann]

	gend := db[1-1]
	ymdB := db[2-1]
	ageRange := calcAgeRange(ymdB, sinryoYm)
	gakuIdx := 0
	if sort == "調剤" {
		gakuIdx = 1
	} else if sort == "歯科" {
		gakuIdx = 2
	}
	kk := []string{"0", gend, honn + ":0", honn + ":" + gend}
	for _, k := range kk {
		if _, ok := d[k]; !ok {
			d[k] = make(map[int][]int)
		}
		if _, ok := d[k][ageRange]; !ok {
			d[k][ageRange] = []int{0, 0, 0}
		}
		if _, ok := d[k][-5]; !ok {
			d[k][-5] = []int{0, 0, 0}
		}
		d[k][ageRange][gakuIdx] += gaku
		d[k][-5][gakuIdx] += gaku
	}
}

func opSummaryC1P3V1(logicOutdir string) {
	for ann := range dicC1P3V1 {
		opSummaryC1P3V1Main(ann, logicOutdir)
	}
}
func opSummaryC1P3V1Main(ann, logicOutdir string) {
	ofnm := logicOutdir + "Result_C1P3V1_" + ann + ".csv"
	oHandle, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_CREATE, 0666)
	defer oHandle.Close()
	writer := bufio.NewWriter(transform.NewWriter(oHandle, japanese.ShiftJIS.NewEncoder()))

	pgDsc := []string{"全体", "男性", "女性", "本人(全体)",
		"本人(男性)", "本人(女性)", "本人外(全体)", "本人外(男性)",
		"本人外(女性)"}
	pgFlg := []string{"0", "1", "2", "1:0", "1:1", "1:2", "2:0", "2:1", "2:2"}
	title := "年齢層,人数,総額(医科＋調剤),医科,調剤,歯科,一人当たり総額,一人当たり[医科＋調剤],一人当たり医科,一人当たり調剤,一人当たり歯科"

	for idx, pg := range pgFlg {
		writer.WriteString("@" + pgDsc[idx] + common.CrLf)
		writer.WriteString(title + common.CrLf)
		for i := -1; i <= 14; i++ {
			ageRange := i * 5
			pop, ok0 := DicPop[pg][ageRange]
			if !ok0 {
				pop = 0
			}
			r1, ok1 := dicC1P3V1[ann][pg][ageRange]
			wk := []int{ageRange}
			if !ok1 {
				wk = append(wk, pop, 0, 0, 0, 0, 0, 0, 0, 0, 0)
			} else {
				med := r1[0]
				pha := r1[1]
				den := r1[2]
				sumAll := med + pha + den
				sum := med + pha
				popDeno := pop
				if popDeno == 0 {
					popDeno = 1
				}
				wk = append(wk, pop, sum, med, pha, den, sumAll/popDeno, sum/popDeno,
					med/popDeno, pha/popDeno, den/popDeno)
			}
			wk2 := []string{}
			for _, one := range wk {
				wk2 = append(wk2, strconv.Itoa(one))
			}
			writer.WriteString(strings.Join(wk2, common.Comma) + common.CrLf)
		}
	}
	writer.Flush()
}
