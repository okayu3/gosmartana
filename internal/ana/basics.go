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

var cntTopicSorts int = len(topicKeys)

//MakeBasicsMED -- 医科レセプトからExpense,SV,toseki作成
func MakeBasicsMED(one [][]string, fnm string, args []interface{}) int {
	var arr []string
	var sort, mnKensaku, ymdB, gend, name, kananame string
	var prf, pnt, ircd, irnm, seikyuYm string
	var sinryoYm, innDate, shubetu string
	var jitsuDates, ten string
	var act, tokki string
	var ck, iKi, iBan string
	var sCntDis, icd10 string
	var aCnt, aIcd10 []string
	if one == nil {
		return -1
	}
	dicDrug := make(map[string][]float64)
	topicCheck := make([]int, cntTopicSorts, cntTopicSorts)
	outHandle := args[0].(*os.File)
	outHandleSV := args[1].(*os.File)
	outHandleTopic := args[2].(*os.File)
	aOutHandlesPDM := args[3].([]*os.File)
	outHandleGene := args[4].(*os.File)
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
		} else if sort == "IY" {
			collectGeneric(arr[rece.RIYmedDrugCd], arr[rece.RIYmedAmount], arr[rece.RIYmedCount], sinryoYm, dicDrug)
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
		prf, pnt, ircd, irnm, seikyuYm, common.Empty)

	opSaveTosekiAndSTopic(outHandleTopic, ck, mnKensaku, tokki, topicCheck)

	OpSavePDMData(aOutHandlesPDM, mnKensaku, jitsuDates, ten, gend, aCnt, aIcd10)

	nyugai := "2"
	if (shubetu != common.Empty) && (int(shubetu[3]-'0')%2 == 1) {
		nyugai = "1"
	}
	opSaveGeneric(outHandleGene, ck, mnKensaku, dicDrug, sinryoYm, nyugai)
	return 1
}

//MakeBasicsDEN -- 医科レセプトからExpense,SV,toseki作成
func MakeBasicsDEN(one [][]string, fnm string, args []interface{}) int {
	var arr []string
	var sort, mnKensaku, ymdB, gend, name, kananame string
	var prf, pnt, ircd, irnm, seikyuYm string
	var sinryoYm, innDate, shubetu string
	var jitsuDates, ten string
	var act, tokki string
	var ck, iKi, iBan string
	var sCntDis, icd10 string
	var aCnt, aIcd10 []string
	if one == nil {
		return -1
	}
	dicDrug := make(map[string][]float64)
	topicCheck := make([]int, cntTopicSorts, cntTopicSorts)
	outHandle := args[0].(*os.File)
	outHandleSV := args[1].(*os.File)
	outHandleTopic := args[2].(*os.File)
	aOutHandlesPDM := args[3].([]*os.File)
	outHandleGene := args[4].(*os.File)
	flgHO := 0
	cntDis := 0

	kbn := "歯科"

	for _, arr = range one {
		if (arr == nil) || (len(arr) <= 0) {
			continue
		}
		sort = arr[rece.RSort]
		if sort == "MN" {
			//歯科のみ
			mnKensaku = arr[rece.RMNkensakuDen]
			irnm = arr[rece.RMNirnameDen]
		} else if sort == "IR" {
			prf = arr[rece.RIRprf]
			pnt = arr[rece.RIRpnttbl]
			ircd = arr[rece.RIRircode]
			//歯科のみ
			seikyuYm = common.YmW2g(arr[rece.RIRseikyuYmDen])
		} else if sort == "RE" {
			//共通(のちにrefactoringすること)
			if mnKensaku == common.Empty {
				mnKensaku = arr[rece.RREdenKensaku]
			}
			ymdB = common.YmdW2g(arr[rece.RREbirth])
			gend = arr[rece.RREgender]
			name = arr[rece.RREname]
			//歯科のみ
			if len(arr) > rece.RREdenkananame {
				kananame = arr[rece.RREdenkananame]
			}
			//独自
			sinryoYm = common.YmW2g(arr[rece.RREsinryoYm])
			innDate = common.YmdW2g(arr[rece.RREinnDate])
			shubetu = arr[rece.RREshubetu]
			//歯科のみ
			tokki = arr[rece.RREtokkiDen]
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
		} else if sort == "IY" {
			collectGeneric(arr[rece.RIYmedDrugCd], arr[rece.RIYmedAmount], arr[rece.RIYmedCount], sinryoYm, dicDrug)
		} else if sort == "SI" {
			act = arr[rece.RSIactCd]
			opCheckActTopic(act, topicCheck)
		} else if sort == "SK" {
			act = arr[rece.RSKactCd]
			opCheckActTopic(act, topicCheck)
		}
	}
	opSaveExpense(outHandle, ck, mnKensaku, gend, ymdB, shubetu,
		jitsuDates, ten, kbn, sinryoYm, innDate,
		prf, pnt, ircd, irnm, seikyuYm, common.Empty)

	opSaveTosekiAndSTopic(outHandleTopic, ck, mnKensaku, tokki, topicCheck)

	OpSavePDMData(aOutHandlesPDM, mnKensaku, jitsuDates, ten, gend, aCnt, aIcd10)

	nyugai := "2"
	if (shubetu != common.Empty) && (int(shubetu[3]-'0')%2 == 1) {
		nyugai = "1"
	}
	opSaveGeneric(outHandleGene, ck, mnKensaku, dicDrug, sinryoYm, nyugai)
	return 1
}

//MakeBasicsPHA -- 調剤レセプトからExpense,SV,toseki作成
func MakeBasicsPHA(one [][]string, fnm string, args []interface{}) int {
	var arr []string
	var sort, mnKensaku, ymdB, gend, name, kananame string
	var prf, pnt, ircd, irnm, seikyuYm string
	var reqIrNo string
	var sinryoYm, innDate, shubetu string
	var jitsuDates, ten string
	var ck, iKi, iBan string
	var count float64

	if one == nil {
		return -1
	}
	outHandle := args[0].(*os.File)
	outHandleGene := args[1].(*os.File)
	dicDrug := make(map[string][]float64)

	flgHO := 0

	kbn := "調剤"
	preSort := common.Empty

	for _, arr = range one {
		if (arr == nil) || (len(arr) <= 0) {
			continue
		}
		preSort = sort
		sort = arr[rece.RSort]
		if sort == "MN" {
			mnKensaku = arr[rece.RMNkensaku]
		} else if sort == "YK" {
			prf = arr[rece.RYKprf]
			pnt = arr[rece.RYKpnttbl]
			ircd = arr[rece.RYKykcode]
			irnm = arr[rece.RYKykname]
			seikyuYm = common.YmW2g(arr[rece.RYKseikyuYm])
		} else if sort == "RE" {
			//共通(のちにrefactoringすること)
			if mnKensaku == common.Empty {
				//調剤独自
				mnKensaku = arr[rece.RREphaKensaku]
			}
			ymdB = common.YmdW2g(arr[rece.RREbirth])
			gend = arr[rece.RREgender]
			name = arr[rece.RREname]
			if len(arr) > rece.RREphakananame {
				//調剤独自
				kananame = arr[rece.RREphakananame]
			}
			//独自
			sinryoYm = common.YmW2g(arr[rece.RREsinryoYm])
			innDate = common.YmdW2g(arr[rece.RREinnDate])
			shubetu = arr[rece.RREshubetu]
			reqIrNo = arr[rece.RREprfPha] + arr[rece.RREpnttblPha] + arr[rece.RREircodePha]
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
		} else if (sort == "SH") || (sort == "RC") {
			count = 0.0
		} else if sort == "CZ" {
			if preSort != "CZ" {
				count = 0.0
			}
			if w, err := strconv.ParseFloat(arr[rece.RCZsuryo], 64); err == nil {
				count += w
			}
		} else if sort == "IY" {
			sCount := strconv.FormatFloat(count, 'f', -1, 64)
			collectGeneric(arr[rece.RIYdrugCd], arr[rece.RIYamount], sCount, sinryoYm, dicDrug)
		}
	}
	opSaveExpense(outHandle, ck, mnKensaku, gend, ymdB, shubetu,
		jitsuDates, ten, kbn, sinryoYm, innDate,
		prf, pnt, ircd, irnm, seikyuYm, reqIrNo)

	opSaveGeneric(outHandleGene, ck, mnKensaku, dicDrug, sinryoYm, "2")
	return 1
}

//MakeBasicsDPC -- DPCレセプトからExpense,SV,toseki,generic作成
func MakeBasicsDPC(one [][]string, fnm string, args []interface{}) int {
	var arr []string
	var sort, mnKensaku, ymdB, gend, name, kananame string
	var prf, pnt, ircd, irnm, seikyuYm string
	var sinryoYm, innDate, shubetu string
	var jitsuDates, ten string
	var act, tokki, betsu string
	var ck, iKi, iBan string
	var sCntDis, icd10, sybcd string
	var aCnt, aIcd10 []string
	var buInnDate, buTenki string
	if one == nil {
		return -1
	}
	sokatsu := common.Empty
	dicDrug := make(map[string][]float64)
	dicDis := make(map[string]int)
	topicCheck := make([]int, cntTopicSorts, cntTopicSorts)
	outHandle := args[0].(*os.File)
	outHandleSV := args[1].(*os.File)
	outHandleTopic := args[2].(*os.File)
	aOutHandlesPDM := args[3].([]*os.File)
	outHandleGene := args[4].(*os.File)
	flgHO := 0
	cntDis := 0

	kbn := "ＤＰＣ"

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
		} else if sort == "RE" {
			//総括区分がカラだったらExpenseに生かさないこと。
			sokatsu = arr[rece.RREdpcSokatsu]
			if sokatsu == "0" || sokatsu == "1" {
				//共通(のちにrefactoringすること)
				if mnKensaku == common.Empty {
					//DPC独自
					mnKensaku = arr[rece.RREdpcKensaku]
				}
				ymdB = common.YmdW2g(arr[rece.RREbirth])
				gend = arr[rece.RREgender]
				name = arr[rece.RREname]
				//DPC独自
				if len(arr) > rece.RREdpckananame {
					kananame = arr[rece.RREdpckananame]
				}
				//独自
				sinryoYm = common.YmW2g(arr[rece.RREsinryoYm])
				innDate = common.YmdW2g(arr[rece.RREinnDate])
				shubetu = arr[rece.RREshubetu]
				tokki = arr[rece.RREtokki]
			}
		} else if sort == "HO" {
			flgHO = 1
			//共通(のちにrefactoringすること)
			//独自
			if sokatsu == "0" || sokatsu == "1" {
				iKi = arr[rece.RHOinsKigo]
				iBan = arr[rece.RHOinsBango]
				ck = ckey.Get(iKi, iBan, ymdB, gend, name, kananame)
				jitsuDates = arr[rece.RHOjitudate]
				ten = arr[rece.RHOten]
			}
		} else if sort == "KO" {
			if sokatsu == "0" || sokatsu == "1" {
				if flgHO == 0 {
					//共通(のちにrefactoringすること)
					iKi = arr[rece.RKOftn]
					iBan = arr[rece.RKOrcv]
					ck = ckey.Get(iKi, iBan, ymdB, gend, name, kananame)
					//独自
					jitsuDates = arr[rece.RKOjitudate]
					ten = arr[rece.RKOten]
				}
			}
		} else if sort == "SY" {
			sybcd = arr[rece.RSYsybcd]
			if _, ok := dicDis[sybcd]; !ok {
				dicDis[sybcd] = 1
				cntDis++
				sCntDis, icd10 = opSaveSV(outHandleSV, ck, mnKensaku, cntDis, gend, arr)
				aCnt = append(aCnt, sCntDis)
				aIcd10 = append(aIcd10, icd10)
			}
		} else if sort == "BU" {
			buInnDate = arr[rece.RBUinnDate]
			buTenki = arr[rece.RBUtenki]
		} else if sort == "SB" {
			sybcd = arr[rece.RSBsybcd]
			if _, ok := dicDis[sybcd]; !ok {
				dicDis[sybcd] = 1
				cntDis++
				sCntDis, icd10 = opSaveSVatSB(outHandleSV, ck, mnKensaku, cntDis, gend, arr,
					buInnDate, buTenki)
				aCnt = append(aCnt, sCntDis)
				aIcd10 = append(aIcd10, icd10)
			}
		} else if sort == "IY" {
			collectGeneric(arr[rece.RIYmedDrugCd], arr[rece.RIYmedAmount], arr[rece.RIYmedCount], sinryoYm, dicDrug)
		} else if sort == "SI" {
			act = arr[rece.RSIactCd]
			opCheckActTopic(act, topicCheck)
		} else if sort == "SK" {
			act = arr[rece.RSKactCd]
			opCheckActTopic(act, topicCheck)
		} else if sort == "CD" {
			act = arr[rece.RCDdpcRececd]
			opCheckActTopic(act, topicCheck)
			betsu = arr[rece.RCDdpcBetu]
			if betsu == "21" || betsu == "22" || betsu == "23" || betsu == "24" ||
				betsu == "25" || betsu == "26" || betsu == "27" || betsu == "28" ||
				betsu == "31" || betsu == "32" || betsu == "33" {
				collectGeneric(act, arr[rece.RCDdpcAmount], arr[rece.RCDdpcCount], sinryoYm, dicDrug)
			}
		}
	}
	opSaveExpense(outHandle, ck, mnKensaku, gend, ymdB, shubetu,
		jitsuDates, ten, kbn, sinryoYm, innDate,
		prf, pnt, ircd, irnm, seikyuYm, common.Empty)

	opSaveTosekiAndSTopic(outHandleTopic, ck, mnKensaku, tokki, topicCheck)

	OpSavePDMData(aOutHandlesPDM, mnKensaku, jitsuDates, ten, gend, aCnt, aIcd10)

	//必ず入院なので nyugai は "1"
	opSaveGeneric(outHandleGene, ck, mnKensaku, dicDrug, sinryoYm, "1")
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
