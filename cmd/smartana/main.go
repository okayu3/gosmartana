package main

import (
	"fmt"
	"os"
	"time"

	"github.com/okayu3/gosmartana/internal/ana"
	"github.com/okayu3/gosmartana/pkg/ckey"
	"github.com/okayu3/gosmartana/pkg/rece"
)

//Main -- smartana main logic
func main() {
	loadMasters()
	start := time.Now()
	//makeSVAndThenExpense() //2.192737秒
	makeExpenseWithSV() //1.280264秒
	goal := time.Now()
	ckey.UpdatePersonMst()
	fmt.Printf("%f秒\n", (goal.Sub(start)).Seconds())
}

func loadMasters() {
	//load mst B / cd119
	fnmMstB := "C:/task/prj/YG01/mst/2020/b/b_20200301.txt"
	fnmMstCd119 := "C:/task/prj/YG01/mst/2020/etc/SYB_MIDDLE_CD119_2013ICD10_NM_202004.csv"
	fnmMstHB := "C:/task/prj/YG01/mst/2020/h/hb_20200101.txt"
	fnmMstSTopic := "C:/task/prj/YG01/mst/2020/D_mst_toseki_go.csv"
	ana.LoadDisB(fnmMstB, fnmMstCd119, fnmMstHB, fnmMstSTopic)
	fnmPsnMst := "C:/Users/woodside3/go/output/person.csv"
	ckey.LoadPersonMst(fnmPsnMst)
}

func makeSVAndThenExpense() {

	//make sv
	//fnm := "C:/task/garden/py/smartana/sample01.csv"
	fnm := "C:/task/garden/py/smartana/sample/11_RECODEINFO_MED.CSV"
	ofnm := "C:/Users/woodside3/go/output/disease.csv"
	ofile, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer ofile.Close()
	rece.Load(fnm, ana.MakeSVMed, []interface{}{ofile})

	//make expense
	ofnmExp := "C:/Users/woodside3/go/output/expense.csv"
	ofileExp, _ := os.OpenFile(ofnmExp, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer ofileExp.Close()
	rece.Load(fnm, ana.MakeExpenseMed, []interface{}{ofileExp})

}

func makeExpenseWithSV() {
	//make sv
	//fnm := "C:/task/garden/py/smartana/sample01.csv"
	fnm := "C:/task/garden/py/smartana/sample/11_RECODEINFO_MED.CSV"
	ofnm := "C:/Users/woodside3/go/output/disease.csv"
	ofileSV, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer ofileSV.Close()
	//make expense
	ofnmExp := "C:/Users/woodside3/go/output/expense.csv"
	ofileExp, _ := os.OpenFile(ofnmExp, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer ofileExp.Close()
	//make toseki and s-topics
	ofnmTopic := "C:/Users/woodside3/go/output/tosekiTopic.csv"
	ofileTopic, _ := os.OpenFile(ofnmTopic, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer ofileTopic.Close()

	rece.Load(fnm, ana.MakeBasicsMed, []interface{}{ofileExp, ofileSV, ofileTopic})
}
