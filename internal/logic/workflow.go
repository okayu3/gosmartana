package logic

import (
	"log"
	"strconv"
	"strings"

	"github.com/okayu3/gosmartana/pkg/common"
	"github.com/okayu3/gosmartana/pkg/rece"
	"golang.org/x/text/width"
)

//MstB -- disease Mst
var MstB = make(map[string][]string)

//MstMiz -- Mizushima-Group DisGroup Definition
var MstMiz = make(map[string]map[string]bool)

//MstCancer -- Cancer to Idx
var MstCancer = make(map[string]int)

//MstCancerDsc -- Cancer Nm
var MstCancerDsc = make(map[string]string)

//CstCancerNum -- Cancer Sort Num
var CstCancerNum = 0

//DicPop  -- Population of Age Range By 5years
var DicPop = make(map[string]map[int]int)

//DicPsn -- Person Data
var DicPsn = make(map[string][]string)

//DicExp -- Expense Data
var DicExp = make(map[string][]string)

//RunLogic -- logic running
//  param: outDir := "C:/Users/woodside3/go/output/"
func RunLogic(mstDir, outDir, settingDir, tokutaiIki, ninkeiIki string) {
	loadMstB(mstDir + "2020/b/b_20200301.txt")
	loadMstMizushima(mstDir + "mst_mizushima.csv")
	loadMstCancer(mstDir+"mst_cancer0.csv", mstDir+"mst_cancer1.csv")
	loadPopulation(settingDir + "setting_population.csv")
	loadPerson(outDir+"person.csv", tokutaiIki, ninkeiIki)
	preLoadTosekiAndTopics(outDir + "tosekiTopic.csv")
	loadExpense(outDir + "expense.csv")
	loadDisPdm(outDir + "diseasePdm.csv")
	loadTosekiAndTopics(outDir + "tosekiTopic.csv")
	opSummary(outDir)
}

func loadMstB(fnmB string) {
	if !common.FileExists(fnmB) {
		log.Println("Disease(B) Mst Not Found:" + fnmB)
		return
	}
	common.LoadCSV(fnmB, func(a []string, lineno int) {
		sybcd := a[3-1]
		sybnm := a[6-1]
		MstB[sybcd] = []string{sybnm}
	}, common.ModeCsvSJIS)
}
func loadMstMizushima(fnmMiz string) {
	if !common.FileExists(fnmMiz) {
		log.Println("Mizushima DisGroup Mst Not Found:" + fnmMiz)
		return
	}
	common.LoadCSV(fnmMiz, func(a []string, lineno int) {
		grCd1 := a[1-1]
		grCd2 := a[2-1]
		sybCd := a[6-1]
		if _, ok := MstMiz[grCd1]; !ok {
			MstMiz[grCd1] = make(map[string]bool)
		}
		if _, ok := MstMiz[grCd2]; !ok {
			MstMiz[grCd2] = make(map[string]bool)
		}
		MstMiz[grCd1][sybCd] = true
		MstMiz[grCd2][sybCd] = true
	}, common.ModeCsvSJIS)
}

func loadMstCancer(fnmCd119, fnmIcd10 string) {
	if !common.FileExists(fnmCd119) {
		log.Println("Cancer Mst(0:cd119) Not Found:" + fnmCd119)
		return
	}
	if !common.FileExists(fnmIcd10) {
		log.Println("Cancer Mst(1:ICD10) Not Found:" + fnmIcd10)
		return
	}
	var idx = 0
	common.LoadCSV(fnmCd119, func(a []string, lineno int) {
		cd := a[1-1]
		nm := a[2-1]
		MstCancer[cd] = idx
		nm = strings.Replace(nm, "の悪性新生物", "", -1)
		nm = strings.Replace(nm, "＜腫瘍＞", "", -1)
		MstCancerDsc[cd] = nm
		idx++
	}, common.ModeCsvSJIS)
	common.LoadCSV(fnmIcd10, func(a []string, lineno int) {
		cd := a[1-1]
		nm := a[2-1]
		MstCancer[cd] = idx
		nm = strings.Replace(nm, "の悪性新生物", "", -1)
		nm = strings.Replace(nm, "＜腫瘍＞", "", -1)
		MstCancerDsc[cd] = nm
		idx++
	}, common.ModeCsvSJIS)
	CstCancerNum = idx
}

func loadPopulation(fnmPop string) {
	if !common.FileExists(fnmPop) {
		log.Println("Population Mst Not Found:" + fnmPop)
		return
	}
	m := map[string]string{"計": "0",
		"男性":     "1",
		"女性":     "2",
		"被保険者計":  "1:0",
		"被保険者男性": "1:1",
		"被保険者女性": "1:2",
		"被扶養者計":  "2:0",
		"被扶養者男性": "2:1",
		"被扶養者女性": "2:2",
		"特退任継計":  "tn:0",
		"特退任継男性": "tn:1",
		"特退任継女性": "tn:2",
	}
	common.LoadCSV(fnmPop, func(a []string, lineno int) {
		k, ok := m[a[0]]
		if !ok {
			return
		}
		for i := 0; i <= 15; i++ {
			if _, ok := DicPop[k]; !ok {
				DicPop[k] = make(map[int]int)
			}
			DicPop[k][(i-1)*5] = common.Atoi(a[i+1], 0)
		}
	}, common.ModeCsvSJIS)

}

func loadPerson(fnmPsnMst, tokutaiIki, ninkeiIki string) {
	if !common.FileExists(fnmPsnMst) {
		log.Println("Person Mst Not Found:" + fnmPsnMst)
		return
	}
	tokutaiW := width.Widen.String(tokutaiIki)
	ninkeiW := width.Widen.String(ninkeiIki)

	common.LoadCSV(fnmPsnMst, func(a []string, lineno int) {
		ck := a[0]
		iKi := a[2-1]
		gend := a[5-1]
		ymdB := a[6-1]
		sort := "0"
		if iKi == tokutaiW || iKi == tokutaiIki {
			sort = "1"
		} else if iKi == ninkeiW || iKi == ninkeiIki {
			sort = "2"
		}
		DicPsn[ck] = []string{gend, ymdB, sort}
	}, common.ModeCsvUTF8)
}

func loadExpense(fnmExpense string) {
	if !common.FileExists(fnmExpense) {
		log.Println("Expense Data Not Found:" + fnmExpense)
		return
	}
	common.LoadCSV(fnmExpense, func(a []string, lineno int) {
		ck := a[1-1]
		mnKensaku := a[2-1]
		nyugai := a[6-1]
		honn := a[7-1]
		jitsuDates := a[8-1]
		ten := a[9-1]
		sort := a[10-1]
		sinryoYm := a[11-1]
		seikyuYm := a[17-1]

		gaku := common.Atoi(ten, 0) * rece.GakuForTen
		DicExp[mnKensaku] = []string{ck, honn, sinryoYm, sort, seikyuYm}

		loadingExpenseC1P1V1(ck, mnKensaku, sort, nyugai, honn, sinryoYm, jitsuDates, seikyuYm, gaku)
		loadingExpenseC1P2V1(ck, mnKensaku, sort, nyugai, honn, sinryoYm, jitsuDates, seikyuYm, gaku)
		loadingExpenseC1P3V1(ck, mnKensaku, sort, nyugai, honn, sinryoYm, jitsuDates, seikyuYm, gaku)
		loadingExpenseC2P6V1(ck, mnKensaku, sort, nyugai, honn, sinryoYm, jitsuDates, seikyuYm, gaku)

	}, common.ModeCsvUTF8)
}

func loadDisPdm(fnmDisPdm string) {
	if !common.FileExists(fnmDisPdm) {
		log.Println("Disease(with Pdm) Data Not Found:" + fnmDisPdm)
		return
	}
	common.LoadCSV(fnmDisPdm, func(a []string, lineno int) {
		ck := a[1-1]
		mnKensaku := a[2-1]
		sybcd := a[7-1]
		cd119 := a[10-1]
		icd10 := a[11-1]
		flgDoubt := a[17-1]
		gaku := common.Atoi(a[20-1], 0)
		if gaku == 0 {
			return
		}
		loadingDisPdmC1P2V1(ck, mnKensaku, cd119, gaku)
		loadingDisPdmC2P4V1(ck, mnKensaku, cd119, sybcd, flgDoubt, gaku)
		loadingDisPdmC2P5V1(ck, mnKensaku, cd119, sybcd, flgDoubt, gaku)
		loadingDisPdmC2P6V1(ck, mnKensaku, cd119, sybcd, flgDoubt, gaku)
		loadingDisPdmC3P10V1(ck, mnKensaku, cd119, icd10, sybcd, flgDoubt, gaku)
		loadingDisPdmC3P11V1(ck, mnKensaku, cd119, icd10, sybcd, flgDoubt, gaku)
		loadingDisPdmC4P15V1(ck, mnKensaku, cd119, icd10, sybcd, flgDoubt, gaku)

	}, common.ModeCsvUTF8)
}

func preLoadTosekiAndTopics(fnmTosekiTopic string) {
	if !common.FileExists(fnmTosekiTopic) {
		log.Println("Toseki & Topics Data Not Found:" + fnmTosekiTopic)
		return
	}
	common.LoadCSV(fnmTosekiTopic, func(a []string, lineno int) {
		ck := a[1-1]
		flgToseki := a[3-1]
		preLoadingTosekiAndTopicsC2P6V1(ck, flgToseki)
	}, common.ModeCsvUTF8)
}

func loadTosekiAndTopics(fnmTosekiTopic string) {
	if !common.FileExists(fnmTosekiTopic) {
		log.Println("Toseki & Topics Data Not Found:" + fnmTosekiTopic)
		return
	}
	common.LoadCSV(fnmTosekiTopic, func(a []string, lineno int) {
		ck := a[1-1]
		mnKensaku := a[2-1]
		flgToseki := a[3-1]
		flgInsulin := a[6-1]
		flgMngDiabetes := a[7-1]
		flgMngBP := a[8-1]
		flgMngFat := a[9-1]
		flgSmoking := a[10-1]
		flgYoboToseki := a[11-1]
		flgTestHbA1c := a[12-1]
		flgTestFat := a[13-1]
		//for 002_007_001 From Here
		loadingTosekiTopicC2P4V1(ck, mnKensaku, flgToseki, flgInsulin, flgMngDiabetes, flgMngBP, flgMngFat, flgSmoking, flgYoboToseki, flgTestHbA1c, flgTestFat)
		loadingTosekiTopicC2P6V1(ck, mnKensaku, flgToseki, flgInsulin, flgMngDiabetes, flgMngBP, flgMngFat, flgSmoking, flgYoboToseki, flgTestHbA1c, flgTestFat)
		loadingTosekiTopicC2P7V1(ck, mnKensaku, flgToseki, flgInsulin, flgMngDiabetes, flgMngBP, flgMngFat, flgSmoking, flgYoboToseki, flgTestHbA1c, flgTestFat)
		//for 002_007_001 Till Here
	}, common.ModeCsvUTF8)
}

func opSummary(outDir string) {
	common.MakeDir(outDir + "logic")
	logicOutdir := outDir + "logic/"

	opSummaryC1P1V1(logicOutdir)
	opSummaryC1P2V1(logicOutdir)
	opSummaryC1P3V1(logicOutdir)
	opSummaryC2P4V1(logicOutdir)
	opSummaryC2P5V1(logicOutdir)
	opSummaryC2P6V1(logicOutdir)
	opSummaryC2P7V1(logicOutdir)
	opSummaryC3P10V1(logicOutdir)
	opSummaryC3P11V1(logicOutdir)
	opSummaryC4P15V1(logicOutdir)
}

func calcAgeRange(ymdB, sinryoYm string) int {
	age := common.AgeAt(ymdB, strconv.Itoa(common.AnnualAtYm(sinryoYm))+"0401")
	if age >= 75 {
		age = 70
	}
	if age < 0 {
		age = 0
	}
	ageRange := (age / 5) * 5
	return ageRange
}

func calcReceAnnual(seikyuYm string) string {
	var i int
	var err0 error
	i, err0 = strconv.Atoi(seikyuYm)
	if err0 != nil {
		return "1900"
	}
	//year
	y := i / 100
	//month
	m := ((i%100 - 1) % 12) + 1
	if m < 5 {
		y--
	}
	return strconv.Itoa(y)
}
