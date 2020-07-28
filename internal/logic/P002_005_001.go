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

//dicC2P5V1DisMap -- annual -- DisSort -- ckey -- 1
var dicC2P5V1DisMap = make(map[string][]map[string]string)

func loadingDisPdmC2P5V1(ck, mnKensaku, cd119, sybcd, flgDoubt string, gaku int) {
	if flgDoubt != "0" && flgDoubt != "" { //疑い病名なら抜けます
		return
	}
	da, okA := DicExp[mnKensaku]
	if !okA {
		return
	}
	honn := da[2-1]
	seikyuYm := da[5-1]
	ann := calcReceAnnual(seikyuYm)

	//まず主要病名の有無をとります。
	if _, ok := dicC2P5V1DisMap[ann]; !ok {
		dicC2P5V1DisMap[ann] = [](map[string]string){
			make(map[string]string),
			make(map[string]string),
			make(map[string]string),
			make(map[string]string),
			make(map[string]string),
			make(map[string]string),
			make(map[string]string),
			make(map[string]string),
		}
	}
	mm := []string{"G0010", //糖尿病
		"G0040", //高脂血症
		"G0050", //高尿酸血症
		"",      //精神・神経
		"G0020", //高血圧
		"G0120", //脳内出血等
		"G0160", //虚血性心疾患等
		"G0060", //肝臓疾患
	}
	d := dicC2P5V1DisMap[ann]
	for idx, m := range mm {
		if m == "" {
			continue
		}
		if _, ok := MstMiz[m][sybcd]; ok {
			d[idx][ck] = honn
		}
	}
	//精神・神経
	if (cd119 == "0503") || //統合失調症、統合失調症型障害及び妄想性障害
		(cd119 == "0504") || //気分［感情］障害躁うつ病を含む
		(cd119 == "0505") || //神経症性障害、ストレス関連障害及び身体表現性障害
		(cd119 == "0507") || //その他の精神及び行動の障害
		(cd119 == "0605") { //自律神経系の障害
		d[4-1][ck] = honn
	}
}

func opSummaryC2P5V1(logicOutdir string) {
	for ann := range dicC2P5V1DisMap {
		opSummaryC2P5V1Main(ann, logicOutdir)
	}
}

func opSummaryC2P5V1Main(ann, logicOutdir string) {
	ofnm := logicOutdir + "Result_C2P5V1_" + ann + ".csv"
	oHandle, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_CREATE, 0666)
	defer oHandle.Close()
	writer := bufio.NewWriter(transform.NewWriter(oHandle, japanese.ShiftJIS.NewEncoder()))

	sumDic := make(map[string][]int)

	for idx, disMap := range dicC2P5V1DisMap[ann] {
		for ck, honn := range disMap {
			db, okA := DicPsn[ck]
			if !okA {
				continue
			}
			gend := db[1-1]
			kk := []string{"0", gend, honn + ":0", honn + ":" + gend}
			for _, k := range kk {
				if _, ok := sumDic[k]; !ok {
					sumDic[k] = []int{0, 0, 0, 0, 0, 0, 0, 0}
				}
				sumDic[k][idx]++
			}
		}
	}

	pgDsc := []string{"全体", "男性", "女性", "本人(全体)",
		"本人(男性)", "本人(女性)", "本人外(全体)", "本人外(男性)",
		"本人外(女性)"}
	pgFlg := []string{"0", "1", "2", "1:0", "1:1", "1:2", "2:0", "2:1", "2:2"}
	title := "区分,人数,糖尿病,高脂血症,高尿酸血症,精神・神経,高血圧,脳内出血等,虚血性心疾患等,肝臓疾患"

	writer.WriteString(title + common.CrLf)
	for idx, pg := range pgFlg {
		arr := []string{pgDsc[idx], strconv.Itoa(DicPop[pg][-5])}
		for ii := 0; ii < 8; ii++ {
			arr = append(arr, strconv.Itoa(sumDic[pg][ii]))
		}
		writer.WriteString(strings.Join(arr, common.Comma) + common.CrLf)
	}
	writer.Flush()
}
