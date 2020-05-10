package ana

import (
	"os"
	"strconv"
	"strings"

	"github.com/okayu3/gosmartana/pkg/ckey"
	"github.com/okayu3/gosmartana/pkg/common"
	"github.com/okayu3/gosmartana/pkg/rece"
)

var topicKeys = []string{"TOSEKI_SI_CODE", "TOSEKI_SI_CODE_ATHOME", "FUKUMAKU_SI_CODE",
	"INSULIN_SI_CODE", "MTB_MNG_DIABETES", "MTB_MNG_BP",
	"MTB_MNG_FAT", "SMOKING", "PRECAUT_TOSEKI", "HBA1C_SI_CODE", "FAT_SI_CODE"}

const cntTopicSorts = 11

//MakeBasicsMed -- 医科レセプトからExpense,SV,toseki作成
func MakeBasicsMed(one [][]string, fnm string, args []interface{}) int {
	var arr []string
	var sort, mnKensaku, ymdB, gend, name, kananame string
	var prf, pnt, ircd, irnm, seikyuYm string
	var sinryoYm, innDate, shubetu string
	var jitsuDates, ten string
	var act, tokki string
	var ck, iKi, iBan string
	var sCntDis, icd10 string
	var aCnt, aIcd10 []string
	topicCheck := make([]int, cntTopicSorts, cntTopicSorts)

	if one == nil {
		return -1
	}
	outHandle := args[0].(*os.File)
	outHandleSV := args[1].(*os.File)
	outHandleTopic := args[2].(*os.File)
	aOutHandlesPDM := args[3].([]*os.File)
	flgHO := 0
	cntDis := 0

	kbn := "医科"

	for _, arr = range one {
		if (arr == nil) || (len(arr) <= 0) {
			continue
		}
		sort = arr[rece.RSort]
		if sort == "MN" {
			mnKensaku = arr[rece.RMNkensaku]
		} else if sort == "IR" {
			prf = arr[rece.RIRprf]
			pnt = arr[rece.RIRpnttbl]
			ircd = arr[rece.RIRircode]
			irnm = arr[rece.RIRirname]
			seikyuYm = common.YmW2g(arr[rece.RIRseikyuYm])
		} else if sort == "YK" {
			prf = arr[rece.RYKprf]
			pnt = arr[rece.RYKpnttbl]
			ircd = arr[rece.RYKykcode]
			irnm = arr[rece.RYKykname]
			seikyuYm = common.YmW2g(arr[rece.RYKseikyuYm])
		} else if sort == "RE" {
			//共通(のちにrefactoringすること)
			if mnKensaku == common.Empty {
				mnKensaku = arr[rece.RREkensaku]
			}
			ymdB = common.YmdW2g(arr[rece.RREbirth])
			gend = arr[rece.RREgender]
			name = arr[rece.RREname]
			if len(arr) > rece.RREkananame {
				kananame = arr[rece.RREkananame]
			}
			//独自
			sinryoYm = common.YmW2g(arr[rece.RREsinryoYm])
			innDate = common.YmdW2g(arr[rece.RREinnDate])
			shubetu = arr[rece.RREshubetu]
			tokki = arr[rece.RREtokki]
		} else if sort == "HO" {
			flgHO = 1
			//共通(のちにrefactoringすること)
			iKi = arr[rece.RHOinsKigo]
			iBan = arr[rece.RHOinsBango]
			ck = ckey.Get(iKi, iBan, ymdB, gend, name, kananame)
			//独自
			jitsuDates = arr[rece.RHOjitudate]
			ten = arr[rece.RHOten]
		} else if sort == "KO" {
			if flgHO == 0 {
				//共通(のちにrefactoringすること)
				iKi = arr[rece.RKOftn]
				iBan = arr[rece.RKOrcv]
				ck = ckey.Get(iKi, iBan, ymdB, gend, name, kananame)
				//独自
				jitsuDates = arr[rece.RKOjitudate]
				ten = arr[rece.RKOten]
			}
		} else if sort == "SY" {
			cntDis++
			sCntDis, icd10 = opSaveSV(outHandleSV, ck, mnKensaku, cntDis, gend, arr)
			aCnt = append(aCnt, sCntDis)
			aIcd10 = append(aIcd10, icd10)
		} else if sort == "SI" {
			act = arr[rece.RSIactCd]
			opCheckActTopic(act, topicCheck)
		} else if sort == "SK" {
			act = arr[rece.RSKactCd]
			opCheckActTopic(act, topicCheck)
		} else if sort == "CD" {
			act = arr[rece.RCDdpcRececd]
			opCheckActTopic(act, topicCheck)
		}
	}
	opSaveExpense(outHandle, ck, mnKensaku, gend, ymdB, shubetu,
		jitsuDates, ten, kbn, sinryoYm, innDate,
		prf, pnt, ircd, irnm, seikyuYm)

	opSaveTosekiAndSTopic(outHandleTopic, ck, mnKensaku, tokki, topicCheck)

	OpSavePDMData(aOutHandlesPDM, mnKensaku, jitsuDates, ten, gend, aCnt, aIcd10)
	return 1
}

func opCheckActTopic(act string, topicCheck []int) {
	for i, kk := range topicKeys {
		if _, ok := MstSTopic[kk][act]; ok {
			topicCheck[i]++
		}
	}
}

func opSaveTosekiAndSTopic(outHandle *os.File, ck, mnKensaku, tokki string, topicCheck []int) {
	if common.IsLongCareRece(tokki) {
		if topicCheck[0] > 0 {
			topicCheck[0] = 1
		}
		if topicCheck[1] > 0 {
			topicCheck[1] = 1
		}
	} else {
		topicCheck[0] = 0
		topicCheck[1] = 0
	}
	var cntTopics int
	for _, v := range topicCheck {
		cntTopics += v
	}
	if cntTopics <= 0 {
		return
	}
	var vals []string
	vals = append(vals, ck)
	vals = append(vals, mnKensaku)
	for _, v := range topicCheck {
		vals = append(vals, strconv.Itoa(v))
	}
	oneTosekiSTopic := strings.Join(vals, common.Comma)
	outHandle.WriteString(oneTosekiSTopic + "\n")
}
