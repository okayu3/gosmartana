package logic

import (
	"log"
	"strconv"

	"github.com/okayu3/gosmartana/pkg/common"
	"github.com/okayu3/gosmartana/pkg/rece"
)

//DicPop  -- Population of Age Range By 5years
var DicPop = make(map[string]map[int]int)

//DicPsn -- Person Data
var DicPsn = make(map[string][]string)

//DicExp -- Expense Data
var DicExp = make(map[string][]string)

//RunLogic -- logic running
//  param: outDir := "C:/Users/woodside3/go/output/"
func RunLogic(outDir string, settingDir string) {
	loadPopulation(settingDir + "setting_population.csv")
	loadPerson(outDir + "person.csv")
	loadExpense(outDir + "expense.csv")
	loadDisPdm(outDir + "diseasePdm.csv")
	opSummary(outDir)
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

func loadPerson(fnmPsnMst string) {
	if !common.FileExists(fnmPsnMst) {
		log.Println("Person Mst Not Found:" + fnmPsnMst)
		return
	}
	common.LoadCSV(fnmPsnMst, func(a []string, lineno int) {
		ck := a[0]
		gend := a[4]
		ymdB := a[5]
		DicPsn[ck] = []string{gend, ymdB}
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

		gaku := common.Atoi(ten, 0) * rece.GakuForTen
		DicExp[mnKensaku] = []string{ck, honn, sinryoYm, sort}

		loadingExpenseC1P1V1(ck, mnKensaku, sort, nyugai, honn, sinryoYm, jitsuDates, gaku)
		loadingExpenseC1P2V1(ck, mnKensaku, sort, nyugai, honn, sinryoYm, jitsuDates, gaku)
		loadingExpenseC1P3V1(ck, mnKensaku, sort, nyugai, honn, sinryoYm, jitsuDates, gaku)

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
		cd119 := a[10-1]
		gaku := common.Atoi(a[20-1], 0)
		if gaku == 0 {
			return
		}
		//for 001_002_001 From Here
		loadingDisPdmC1P2V1(ck, mnKensaku, cd119, gaku)
		//for 001_002_001 Till Here
	}, common.ModeCsvUTF8)
}

func opSummary(outDir string) {
	common.MakeDir(outDir + "logic")
	logicOutdir := outDir + "logic/"

	opSummaryC1P1V1(logicOutdir)
	opSummaryC1P2V1(logicOutdir)
	opSummaryC1P3V1(logicOutdir)
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
