package logic

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/okayu3/gosmartana/pkg/common"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// dicC2P4V1DisMap -- annual -- ckey -- honn
var dicC2P4V1DisMap = make(map[string][]map[string]string)
var mmC2P4V1 = []string{
	"D0010", //糖尿病
	"D0020", //高血圧症
	"D0040", //高脂血症
	"D0050", //高尿酸血症
	"D0060", //肝障害
	"D0065", //動脈硬化
	"D0070", //糖尿病性神経症
	"D0080", //糖尿病性網膜症
	"D0090", //糖尿病性腎症
	"D0100", //痛風腎
	"D0110", //高血圧性腎臓障害
	"D0120", //脳血管疾患
	"D0130", //脳出血
	"D0140", //脳梗塞
	"D0150", //その他の脳血管疾患
	"D0160", //虚血性心疾患
	"D0170", //動脈閉塞
	"D0180", //大動脈疾患
}

func loadingDisPdmC2P4V1(ck, mnKensaku, cd119, sybcd, flgDoubt string, gaku int) {
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
	if _, ok := dicC2P4V1DisMap[ann]; !ok {
		dicC2P4V1DisMap[ann] = make([](map[string]string), len(mmC2P4V1))
		for i := 0; i < len(mmC2P4V1); i++ {
			dicC2P4V1DisMap[ann][i] = make(map[string]string)
		}
	}

	d := dicC2P4V1DisMap[ann]
	for idx, m := range mmC2P4V1 {
		if m == "" {
			continue
		}
		if _, ok := MstMiz[m][sybcd]; ok {
			d[idx][ck] = honn
		}
	}

}

func opSummaryC2P4V1(logicOutdir string) {
	for ann := range dicC2P4V1DisMap {
		opSummaryC2P4V1Main(ann, logicOutdir)
	}
}

func opSummaryC2P4V1Main(ann, logicOutdir string) {
	ofnm := logicOutdir + "Result_C2P4V1_" + ann + ".csv"
	oHandle, _ := os.Create(ofnm)
	defer oHandle.Close()
	writer := bufio.NewWriter(transform.NewWriter(oHandle, japanese.ShiftJIS.NewEncoder()))

	var prefix string
	sumDic := make(map[string][]int)
	mmCnt := len(mmC2P4V1)
	for idx, oneDis := range dicC2P4V1DisMap[ann] {
		for ck, honn := range oneDis {
			db, okA := DicPsn[ck]
			if !okA {
				continue
			}
			gend := db[1-1]
			sort := db[3-1]

			if sort == "0" {
				prefix = honn
			} else {
				prefix = "tn"
			}
			kk := []string{"0", gend, prefix + ":0", prefix + ":" + gend}
			for _, k := range kk {
				if _, ok := sumDic[k]; !ok {
					sumDic[k] = make([]int, mmCnt)
				}
				sumDic[k][idx]++
			}
		}
	}
	pgDsc := []string{"全体", "男性", "女性", "一般本人(全体)",
		"一般本人(男性)", "一般本人(女性)", "一般家族(全体)", "一般家族(男性)",
		"一般家族(女性)", "特退任継(全体)", "特退任継(男性)", "特退任継(女性)"}
	pgFlg := []string{"0", "1", "2", "1:0", "1:1", "1:2", "2:0", "2:1", "2:2", "tn:0", "tn:1", "tn:2"}

	mmDsc := map[string]string{
		"D0010": "糖尿病",
		"D0020": "高血圧症",
		"D0040": "高脂血症",
		"D0050": "高尿酸血症",
		"D0060": "肝障害",
		"D0065": "動脈硬化",
		"D0070": "糖尿病性神経症",
		"D0080": "糖尿病性網膜症",
		"D0090": "糖尿病性腎症",
		"D0100": "痛風腎",
		"D0110": "高血圧性腎臓障害",
		"D0120": "脳血管疾患",
		"D0130": "脳出血",
		"D0140": "脳梗塞",
		"D0150": "その他の脳血管疾患",
		"D0160": "虚血性心疾患",
		"D0170": "動脈閉塞",
		"D0180": "大動脈疾患",
	}

	for idx, pg := range pgFlg {
		if _, ok := sumDic[pg]; !ok {
			continue
		}
		writer.WriteString("@" + pgDsc[idx] + common.CrLf)
		wk0 := []string{"加入者数", strconv.Itoa(DicPop[pg][-5])}
		writer.WriteString(strings.Join(wk0, common.Comma) + common.CrLf)
		for ii, one := range mmC2P4V1 {
			ddsc := mmDsc[one]
			ratio := fmt.Sprintf("%f", float64(sumDic[pg][ii])/float64(DicPop[pg][-5]))
			wk := []string{one, ddsc, strconv.Itoa(sumDic[pg][ii]), ratio}
			writer.WriteString(strings.Join(wk, common.Comma) + common.CrLf)
		}
	}
	writer.Flush()
}
