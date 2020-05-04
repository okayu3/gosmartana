package main

import (
	"os"

	"github.com/okayu3/gosmartana/internal/ana"
	"github.com/okayu3/gosmartana/pkg/rece"
)

//Main -- smartana main logic
func main() {
	//load mst B / cd119
	fnmMstB := "C:/task/prj/YG01/mst/2020/b/b_20200301.txt"
	fnmMstCd119 := "C:/task/prj/YG01/mst/2020/etc/SYB_MIDDLE_CD119_2013ICD10_NM_202004.csv"
	fnmMstHB := "C:/task/prj/YG01/mst/2020/h/hb_20200101.txt"
	ana.LoadDisB(fnmMstB, fnmMstCd119, fnmMstHB)

	//make sv
	//fnm := "C:/task/garden/py/smartana/sample01.csv"
	fnm := "C:/task/garden/py/smartana/sample/11_RECODEINFO_MED.CSV"
	ofnm := "C:/Users/woodside3/go/output/sv001.csv"
	ofile, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer ofile.Close()
	rece.Load(fnm, ana.MakeSVMed, []interface{}{ofile})
}
