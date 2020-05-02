package common

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

//const -- CSV load mode
const (
	ModeCsvSJIS = 0
	ModeCsvUTF8 = 1
	ModeCsvEUC  = 2
	ModeTsvSJIS = 3
	ModeTsvUTF8 = 4
	ModeTsvEUC  = 5
	TickLineNum = 100000
)

//CsvCallback -- type of callback for CSV operation
type CsvCallback func([]string, int)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func loadCsvMain(in *os.File, r CsvCallback, mode int) {
	var reader *csv.Reader
	var lineno = 0
	if (mode == ModeCsvSJIS) || (mode == ModeTsvSJIS) {
		reader = csv.NewReader(transform.NewReader(in, japanese.ShiftJIS.NewDecoder()))
	} else if (mode == ModeCsvEUC) || (mode == ModeTsvEUC) {
		reader = csv.NewReader(transform.NewReader(in, japanese.EUCJP.NewDecoder()))
	} else { //UTF8 : default
		reader = csv.NewReader(in)
	}
	if (mode == ModeTsvEUC) || (mode == ModeTsvSJIS) || (mode == ModeTsvUTF8) {
		reader.Comma = '\t'
	}
	for {
		arr, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			failOnError(err)
		}
		lineno++
		if lineno%TickLineNum == 0 {
			fmt.Printf("[loading %d lines...]\n", lineno)
		}
		r(arr, lineno)
	}
	fmt.Printf("[loading %d lines...]\n", lineno)
}

//LoadCSV -- reading one csv file
func LoadCSV(fnm string, r CsvCallback, mode int) {
	f, err := os.Open(fnm)
	failOnError(err)
	defer f.Close()
	loadCsvMain(f, r, mode)
}

//LoadCSVArr -- reading csv files
func LoadCSVArr(fnms []string, r CsvCallback, mode int) {
	for _, fnm := range fnms {
		LoadCSV(fnm, r, mode)
	}
}
