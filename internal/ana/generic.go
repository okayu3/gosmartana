package ana

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/okayu3/gosmartana/pkg/common"
	"github.com/okayu3/gosmartana/pkg/rece"
)

type geneFlgRange struct {
	fromYM int
	tillYM int
	stat   string
}

type yakkaRange struct {
	fromYM int
	tillYM int
	yakka  float64
}

var mstGeneFlg = make(map[string][]geneFlgRange)
var mstDrgYjcd = make(map[string]string)
var mstDrgYakka = make(map[string][]yakkaRange)
var mstDrgCheap = make(map[string][]yakkaRange)
var mstDrgExpsv = make(map[string][]yakkaRange)

//LoadMstGeneric -- ジェネリック関連マスタの読み込み
func LoadMstGeneric(fnmAnyDrg, fnmGeneFlg, fnmCheap, fnmExpensv string) {
	var yjcd, drgcd string
	common.LoadCSV(fnmAnyDrg, func(arr []string, lineno int) {
		yjcd = arr[0]
		drgcd = arr[1]
		mstDrgYjcd[drgcd] = yjcd
		for i := 2; i < len(arr); i++ {
			mstDrgYakka[drgcd] = append(mstDrgYakka[drgcd], newYakkaRange(arr[i]))
		}
	}, common.ModeCsvSJIS)
	common.LoadCSV(fnmGeneFlg, func(arr []string, lineno int) {
		drgcd = arr[1]
		for i := 2; i < len(arr); i++ {
			mstGeneFlg[drgcd] = append(mstGeneFlg[drgcd], newGeneFlgRange(arr[i]))
		}
	}, common.ModeCsvSJIS)
	common.LoadCSV(fnmCheap, func(arr []string, lineno int) {
		yjcd = arr[0]
		for i := 1; i < len(arr); i++ {
			mstDrgCheap[yjcd] = append(mstDrgCheap[yjcd], newYakkaRange(arr[i]))
		}
	}, common.ModeCsvSJIS)
	common.LoadCSV(fnmExpensv, func(arr []string, lineno int) {
		yjcd = arr[0]
		for i := 1; i < len(arr); i++ {
			mstDrgExpsv[yjcd] = append(mstDrgExpsv[yjcd], newYakkaRange(arr[i]))
		}
	}, common.ModeCsvSJIS)
}

func yakkaAt(drcd, ym string) float64 {
	atym, err := strconv.Atoi(ym)
	if err != nil {
		//現時点のYMをとって、そのYYYYMMを整数化する。
		t := time.Now()
		atym, _ = strconv.Atoi(t.Format("200601"))
	}
	if (atym > 197001) && (atym < 201304) {
		atym = 201304
	}
	sli, ok := mstDrgYakka[drcd]
	if !ok {
		return 0
	}
	for _, v := range sli {
		if (v.fromYM <= atym) && (atym <= v.tillYM) {
			return v.yakka
		}
	}
	return 0
}

func geneFlgAt(drcd, ym string) string {
	atym, err := strconv.Atoi(ym)
	if err != nil {
		//現時点のYMをとって、そのYYYYMMを整数化する。
		t := time.Now()
		atym, _ = strconv.Atoi(t.Format("200601"))
	}
	if (atym > 197001) && (atym < 201304) {
		atym = 201304
	}
	sli, ok := mstGeneFlg[drcd]
	if !ok {
		return "0"
	}
	for _, v := range sli {
		if (v.fromYM <= atym) && (atym <= v.tillYM) {
			return v.stat
		}
	}
	return "0"
}

func geneCheapExpsvAt(drcd, ym string) (float64, float64) {
	atym, err := strconv.Atoi(ym)
	if err != nil {
		//現時点のYMをとって、そのYYYYMMを整数化する。
		t := time.Now()
		atym, _ = strconv.Atoi(t.Format("200601"))
	}
	if (atym > 197001) && (atym < 201304) {
		atym = 201304
	}
	yjcd, ok0 := mstDrgYjcd[drcd]
	if !ok0 {
		return 0, 0
	}
	var costC, costE float64
	sli, ok := mstDrgCheap[yjcd[0:9]]
	if !ok {
		costC = 0
	}
	for _, v := range sli {
		if (v.fromYM <= atym) && (atym <= v.tillYM) {
			costC = v.yakka
		}
	}
	sli, ok = mstDrgExpsv[yjcd[0:9]]
	if !ok {
		costE = 0
	}
	for _, v := range sli {
		if (v.fromYM <= atym) && (atym <= v.tillYM) {
			costE = v.yakka
		}
	}
	return costC, costE
}

func newYakkaRange(one string) yakkaRange {
	wk := strings.Split(one, common.Collon)
	d := yakkaRange{}
	d.fromYM, _ = strconv.Atoi(wk[0])
	d.tillYM, _ = strconv.Atoi(wk[1])
	d.yakka, _ = strconv.ParseFloat(wk[2], 64)
	return d
}

func newGeneFlgRange(one string) geneFlgRange {
	wk := strings.Split(one, common.Collon)
	d := geneFlgRange{}
	d.fromYM, _ = strconv.Atoi(wk[0])
	d.tillYM, _ = strconv.Atoi(wk[1])
	d.stat = wk[2]
	return d
}

func collectGeneric(drgCd, amount, count, sinryoYm string, dicDrug map[string][]float64) {
	sort := geneFlgAt(drgCd, sinryoYm)
	yakka := yakkaAt(drgCd, sinryoYm)
	if yakka == 0 {
		return
	}
	ykCheap := yakka
	ykExpsv := yakka
	if (sort == "3") || (sort == "2") {
		ykCheap, ykExpsv = geneCheapExpsvAt(drgCd, sinryoYm)
	}
	if _, ok := dicDrug[drgCd]; !ok {
		dicDrug[drgCd] = []float64{0, 0, 0, 0}
	}
	fAmount, err0 := strconv.ParseFloat(amount, 64)
	fCount, err1 := strconv.ParseFloat(count, 64)
	if (err0 != nil) || (err1 != nil) {
		return
	}
	costN := common.Round5sha(yakka*fAmount/10.0) * fCount * rece.GakuForTen
	costC := common.Round5sha(ykCheap*fAmount/10.0) * fCount * rece.GakuForTen
	costE := common.Round5sha(ykExpsv*fAmount/10.0) * fCount * rece.GakuForTen
	if _, ok := dicDrug[drgCd]; !ok {
		dicDrug[drgCd] = []float64{0, 0, 0, 0}
	}
	dicDrug[drgCd][0] += fAmount * fCount
	dicDrug[drgCd][1] += costN
	dicDrug[drgCd][2] += costC
	dicDrug[drgCd][3] += costE
}
func opSaveGeneric(outHandle *os.File, ck, mnKensaku string,
	dicDrug map[string][]float64, sinryoYm, nyugai string) {

	var yjcd, grpcd, seibuncd string
	var ok bool

	for drgCd, v := range dicDrug {
		w := []string{ck, mnKensaku, sinryoYm}
		w = append(w, geneFlgAt(drgCd, sinryoYm))
		yjcd, ok = mstDrgYjcd[drgCd]
		grpcd = common.Empty
		seibuncd = common.Empty
		if !ok || yjcd == common.Empty {
			yjcd = common.Empty
		} else {
			if len(yjcd) < 11 {
				//log.Logger("what yjcd" + yjcd)
				fmt.Println("what yjcd" + yjcd)
			} else {
				grpcd = yjcd[0:9]
				seibuncd = yjcd[9:11]
			}
		}
		w = append(w, drgCd, grpcd, seibuncd, nyugai)
		for _, vv := range v {
			w = append(w, fmt.Sprintf("%.3f", vv))
		}
		oneGene := strings.Join(w, common.Comma)
		outHandle.WriteString(oneGene + "\n")
	}

}
