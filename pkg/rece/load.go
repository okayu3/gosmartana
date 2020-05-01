package rece

import (
	"encoding/csv"
	"io"
	"os"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

//constants -- location for array of receipt
const (
	locSgn    = 0
	locLineno = 1
	locEda    = 2
	locSort   = 2 + 1

	locJYOp    = 4
	locJYObLno = 5
	locJYObEda = 6
	locJYNwLno = 9
	locJYNwEda = 10
)

//Args -- 引数
type Args struct {
	Out *os.File
	Prm []string
}

//Callback -- コールバック関数の型
type Callback func([][]string, Args, string) int

//Load -- レセプトデータの処理
func Load(in *os.File, r Callback, a Args, fnm string) {
	reader := csv.NewReader(transform.NewReader(in, japanese.ShiftJIS.NewDecoder()))
	var sgn, sort, lineno, eda string
	var rece [][]string
	flgJy := false
	for {
		arr, err := reader.Read()
		if err == io.EOF {
			break
		}
		if len(arr) < 3 {
			continue
		}
		sgn = arr[locSgn]
		sort = arr[locSort]
		lineno = arr[locLineno]
		eda = arr[locEda]
		if sort == "MN" {
			r(rece, a, fnm)
			flgJy = false
		}
		if (sgn == "1") || (sort == "MN") {

		}
	}
}
