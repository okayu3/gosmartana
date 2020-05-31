package main

import (
	"fmt"
	"os"
	"time"

	"github.com/okayu3/gosmartana/internal/ana"
	"github.com/okayu3/gosmartana/pkg/ckey"
	"github.com/okayu3/gosmartana/pkg/common"
	"github.com/okayu3/gosmartana/pkg/rece"

	"github.com/pelletier/go-toml"
)

//var mstDir = "C:/task/prj/YG01/mst"

//Main -- smartana main logic
func main() {
	mstDir, outDir, _, receDir := loadSettings()
	loadMasters(mstDir, outDir)
	start := time.Now()
	//makeSVAndThenExpense() //2.192737秒
	makeExpenseWithSV(mstDir, outDir, receDir) //1.280264秒
	goal := time.Now()
	ckey.UpdatePersonMst()
	fmt.Printf("%f秒\n", (goal.Sub(start)).Seconds())
}

func loadSettings() (string, string, string, string) {
	settings, _ := toml.LoadFile("./settings.toml")
	if settings == nil {
		fmt.Println("cant load settings.toml file")
		return "C:/task/prj/YG01/mst/", "C:/Users/woodside3/go/output/",
			"C:/Users/woodside3/go/settings/", "C:/task/prj/YG01/sample/rece"
	}
	mstDir := settings.Get("MasterPath.MST_DIR").(string)
	outDir := settings.Get("OutputPath.OUT_DIR").(string)
	setDir := settings.Get("SettingsPath.SETTING_DIR").(string)
	receDir := settings.Get("RecePath.RECE_DIR").(string)
	return mstDir, outDir, setDir, receDir
}

func loadMasters(mstDir, outDir string) {
	//load mst B / cd119
	fnmMstB := mstDir + "2020/b/b_20200301.txt"
	fnmMstCd119 := mstDir + "2020/etc/SYB_MIDDLE_CD119_2013ICD10_NM_202004.csv"
	fnmMstHB := mstDir + "2020/h/hb_20200101.txt"
	fnmMstSTopic := mstDir + "2020/D_mst_toseki_go.csv"
	ana.LoadDisB(fnmMstB, fnmMstCd119, fnmMstHB, fnmMstSTopic)
	//Person Mst
	fnmPsnMst := outDir + "person.csv"
	ckey.LoadPersonMst(fnmPsnMst)
	//generic mst
	fnmAnyDrg := mstDir + "2020/etc/generic/A015_01_mst_any_yakka_period.csv"
	fnmGeneFlg := mstDir + "2020/etc/generic/A015_01_mst_genestat_period.csv"
	fnmCheap := mstDir + "2020/etc/generic/A015_01_mst_cheapest_period_202004.csv"
	fnmExpensv := mstDir + "2020/etc/generic/A015_01_mst_expensive_period_202004.csv"
	ana.LoadMstGeneric(fnmAnyDrg, fnmGeneFlg, fnmCheap, fnmExpensv)

}

/* func makeSVAndThenExpense() {

	//make sv
	//fnm := "C:/task/garden/py/smartana/sample01.csv"
	fnm := "C:/task/garden/py/smartana/sample/11_RECODEINFO_MED.CSV"
	ofnm := "C:/Users/woodside3/go/output/disease.csv"
	ofile, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_CREATE, 0666)
	defer ofile.Close()
	rece.Load(fnm, ana.MakeSVMed, []interface{}{ofile})

	//make expense
	ofnmExp := "C:/Users/woodside3/go/output/expense.csv"
	ofileExp, _ := os.OpenFile(ofnmExp, os.O_WRONLY|os.O_CREATE, 0666)
	defer ofileExp.Close()
	rece.Load(fnm, ana.MakeExpenseMed, []interface{}{ofileExp})

} */

func makeExpenseWithSV(mstDir, outDir, receDir string) {
	//fnm := "C:/task/garden/py/smartana/sample/11_RECODEINFO_MED.CSV"
	//fnm := "C:/task/garden/py/smartana/sample01.csv"
	receFnms := common.ListUpRece(receDir, common.Empty)

	//make sv
	ofnm := outDir + "disease.csv"
	ofileSV, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_CREATE, 0666)
	defer ofileSV.Close()
	//make expense
	ofnmExp := outDir + "expense.csv"
	ofileExp, _ := os.OpenFile(ofnmExp, os.O_WRONLY|os.O_CREATE, 0666)
	defer ofileExp.Close()
	//make toseki and s-topics
	ofnmTopic := outDir + "tosekiTopic.csv"
	ofileTopic, _ := os.OpenFile(ofnmTopic, os.O_WRONLY|os.O_CREATE, 0666)
	defer ofileTopic.Close()
	//make PDMData
	ohandlesPDM := ana.PreparePDMData(outDir)
	//make generic
	ofnmGene := outDir + "generic.csv"
	ofileGene, _ := os.OpenFile(ofnmGene, os.O_WRONLY|os.O_CREATE, 0666)

	rece.LoadArr(receFnms[0], ana.MakeBasicsMED, []interface{}{ofileExp, ofileSV, ofileTopic, ohandlesPDM, ofileGene})
	rece.LoadArr(receFnms[1], ana.MakeBasicsDPC, []interface{}{ofileExp, ofileSV, ofileTopic, ohandlesPDM, ofileGene})
	rece.LoadArr(receFnms[2], ana.MakeBasicsDEN, []interface{}{ofileExp, ofileSV, ofileTopic, ohandlesPDM, ofileGene})
	rece.LoadArr(receFnms[3], ana.MakeBasicsPHA, []interface{}{ofileExp, ofileGene})
	//rece.Load(fnm, ana.MakeBasicsMed, []interface{}{ofileExp, ofileSV, ofileTopic, ohandlesPDM})
	//closing
	ana.ClosePDMHandle(ohandlesPDM)
	ana.RunPDM(outDir+"pdm/"+ana.PdmDataMale, mstDir, outDir)
	ana.RunPDM(outDir+"pdm/"+ana.PdmDataFemale, mstDir, outDir)
	ana.OpAfterPDM(outDir)

	//	logic.RunLogic(outDir)
}
