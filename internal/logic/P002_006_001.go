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

var dicC2P6V1PatToseki = make(map[string]bool)

//dicC2P6V1PatCost -- annual -- ckey -- [cost: dpc, ika-inn, ika-gai, pha, den]
var dicC2P6V1PatCost = make(map[string]map[string][]int)

//dicC2P6V1DisMap -- annual -- ckey -- DisMap
var dicC2P6V1DisMap = make(map[string]map[string][]bool)

//dicC2P6V1DisCost -- annual -- ckey -- sybcd -- cost
var dicC2P6V1DisCost = make(map[string]map[string]map[string]int)

//dicC2P6V1wk -- annual -- ckey -- cost
var dicC2P6V1wk = make(map[string]map[string][]string)

func preLoadingTosekiAndTopicsC2P6V1(ck, flgToseki string) {
	if flgToseki == "1" {
		dicC2P6V1PatToseki[ck] = true
	}
}

func loadingExpenseC2P6V1(ck, mnKensaku, sort, nyugai, honn, sinryoYm, jitsuDates, seikyuYm string, gaku int) {
	//年度：DPC, 医科入院, 医科外来,調剤,歯科
	if _, ok := dicC2P6V1PatToseki[ck]; !ok {
		return
	}
	ann := calcReceAnnual(seikyuYm)
	if _, ok := dicC2P6V1PatCost[ann]; !ok {
		dicC2P6V1PatCost[ann] = make(map[string][]int)
	}
	if _, ok := dicC2P6V1PatCost[ann][ck]; !ok {
		dicC2P6V1PatCost[ann][ck] = []int{0, 0, 0, 0, 0}
	}
	idx := 2
	if sort == "ＤＰＣ" {
		idx = 0
	} else if sort == "医科" {
		if nyugai == "2" { //外来
			idx = 2
		} else { //入院
			idx = 1
		}
	} else if sort == "調剤" {
		idx = 3
	} else if sort == "歯科" {
		idx = 4
	}
	dicC2P6V1PatCost[ann][ck][idx] += gaku
}

func loadingDisPdmC2P6V1(ck, mnKensaku, cd119, sybcd, flgDoubt string, gaku int) {
	da, okA := DicExp[mnKensaku]
	if !okA {
		return
	}
	//honn := da[2-1]
	//sinryoYm := da[3-1]
	seikyuYm := da[5-1]
	ann := calcReceAnnual(seikyuYm)

	//対象患者でなければ抜けます
	if _, ok := dicC2P6V1PatCost[ann][ck]; !ok {
		return
	}
	//まず主要病名の有無をとります。
	if _, ok := dicC2P6V1DisMap[ann]; !ok {
		dicC2P6V1DisMap[ann] = make(map[string][]bool)
	}
	if _, ok := dicC2P6V1DisMap[ann][ck]; !ok {
		dicC2P6V1DisMap[ann][ck] = []bool{false, false, false, false, false, false, false, false, false, false}
	}
	mm := []string{"G0010", //糖尿病
		"",      //インシュリン療法(病名でない)
		"G0070", //糖尿病性神経症
		"G0080", //糖尿病性網膜症
		"G0170", //動脈閉塞
		"G0020", //高血圧症
		"G0050", //高尿酸血症
		"G0160", //虚血性心疾患
		"G0120", //脳血管疾患
		"G9001", //１型糖尿病
	}
	flg := 0
	d := dicC2P6V1DisMap[ann][ck]
	for idx, m := range mm {
		if m == "" {
			continue
		}
		if _, ok := MstMiz[m][sybcd]; ok {
			//疑い病名でなければ
			if flgDoubt != "1" {
				d[idx] = true
			}
			flg++
		}
	}
	//主要病名でなければ、額を取っていく。
	if flg == 0 {
		if _, ok := dicC2P6V1DisCost[ann]; !ok {
			dicC2P6V1DisCost[ann] = make(map[string]map[string]int)
		}
		if _, ok := dicC2P6V1DisCost[ann][ck]; !ok {
			dicC2P6V1DisCost[ann][ck] = make(map[string]int)
		}
		d := dicC2P6V1DisCost[ann][ck]
		d[sybcd] += gaku
	}
}

func loadingTosekiTopicC2P6V1(ck, mnKensaku,
	flgToseki, flgInsulin, flgMngDiabetes,
	flgMngBP, flgMngFat, flgSmoking,
	flgYoboToseki, flgTestHbA1c, flgTestFat string) {

	da, okA := DicExp[mnKensaku]
	db, okB := DicPsn[ck]
	if !okA || !okB {
		return
	}
	honn := da[2-1]
	seikyuYm := da[5-1]
	gend := db[1-1]
	ymdB := db[2-1]
	ann := calcReceAnnual(seikyuYm)

	if flgInsulin != "0" {
		dc, okC := dicC2P6V1DisMap[ann][ck]
		if okC {
			dc[2-1] = true
		}
	}
	if flgToseki != "0" {
		if _, ok0 := dicC2P6V1wk[ann][ck]; ok0 {
			return
		}
		if _, ok1 := dicC2P6V1wk[ann]; !ok1 {
			dicC2P6V1wk[ann] = make(map[string][]string)
		}
		//annual--ckey で透析患者だった。
		disCost := 0
		if costs, ok := dicC2P6V1PatCost[ann][ck]; ok {
			disCost = costs[0] + costs[1] + costs[2] + costs[3] + costs[4]
		}
		age := common.AgeAt(ymdB, ann+"0401")
		dicC2P6V1wk[ann][ck] = []string{honn, gend, strconv.Itoa(age), strconv.Itoa(disCost)}
	}

}

func opSummaryC2P6V1(logicOutdir string) {
	for ann := range dicC2P6V1wk {
		opSummaryC2P6V1Main(ann, logicOutdir)
	}
}

func opSummaryC2P6V1Main(ann, logicOutdir string) {
	ofnm := logicOutdir + "Result_C2P6V1_" + ann + ".csv"
	oHandle, _ := os.Create(ofnm)
	defer oHandle.Close()
	writer := bufio.NewWriter(transform.NewWriter(oHandle, japanese.ShiftJIS.NewEncoder()))

	a := List{}
	for k, data := range dicC2P6V1wk[ann] {
		e := Entry{k, common.Atoi(data[4-1], 0)}
		a = append(a, e)
	}
	sort.Sort(sort.Reverse(a))

	for _, ent := range a {
		ck := ent.name
		arr := dicC2P6V1wk[ann][ck]
		pcost := dicC2P6V1PatCost[ann][ck]
		arr = append(arr,
			strconv.Itoa(pcost[0]+pcost[1]),
			strconv.Itoa(pcost[2]),
			strconv.Itoa(pcost[3]))

		arr = append(arr, dspEtcDis(ann, ck)...)
		if mp, ok := dicC2P6V1DisMap[ann][ck]; !ok {
			for i := 0; i < 10; i++ {
				arr = append(arr, "")
			}
		} else {
			for i := 0; i < 10; i++ {
				if mp[i] {
					arr = append(arr, "●")
				} else {
					arr = append(arr, "")
				}
			}

		}
		writer.WriteString(strings.Join(arr, common.Comma) + common.CrLf)
	}
	writer.Flush()
}

func dspEtcDis(ann, ck string) []string {
	d, ok := dicC2P6V1DisCost[ann][ck]
	if !ok {
		return []string{"", "", "", "", ""}
	}
	a := List{}
	for k, v := range d {
		e := Entry{k, v}
		a = append(a, e)
	}
	sort.Sort(sort.Reverse(a))
	var ans []string
	for idx, ent := range a {
		if idx > 4 {
			break
		}
		syb := ent.name
		dsc := ""
		if nm, ok := MstB[syb]; ok {
			dsc = nm[0]
		}
		if syb == "0000999" {
			dsc = ""
		}
		ans = append(ans, dsc)
	}
	remain := 5 - a.Len()
	if remain > 0 {
		for i := 0; i < remain; i++ {
			ans = append(ans, "")
		}
	}
	return ans
}

//Entry -- entry for sort
type Entry struct {
	name  string
	value int
}

//List -- list of entry
type List []Entry

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	if l[i].value == l[j].value {
		return (l[i].name < l[j].name)
	}
	return (l[i].value < l[j].value)
}
