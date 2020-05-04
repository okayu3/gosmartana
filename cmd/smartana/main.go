package main

import (
	"os"

	"github.com/okayu3/gosmartana/internal/ana"
	"github.com/okayu3/gosmartana/pkg/rece"
)

//Main -- smartana main logic
func main() {
	fnm := "C:/task/garden/py/smartana/sample01.csv"
	ofnm := "C:/Users/woodside3/go/output/sv001.csv"
	ofile, _ := os.OpenFile(ofnm, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer ofile.Close()
	rece.Load(fnm, ana.MakeSVMed, []interface{}{ofile})
}
