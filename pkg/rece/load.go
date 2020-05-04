package rece

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	//	c "golang.org/okayu3/gosmartana/pkg/common"

	"github.com/okayu3/gosmartana/pkg/common"
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

//Callback -- コールバック関数の型
type Callback func([][]string, string, []interface{}) int

func receSorting(arr [][]string, flgJy bool) [][]string {
	if !flgJy {
		return arr
	}
	sort.Slice(arr, func(i, j int) bool {
		if len(arr[i]) < 3 {
			return true
		} else if len(arr[j]) < 3 {
			return true
		} else if arr[i][1] < arr[j][1] {
			return true
		} else if (arr[i][1] == arr[j][1]) && (arr[i][2] < arr[j][2]) {
			return true
		}
		return false
	})
	return arr
}

//loadMain -- レセプトデータの処理のメイン
func loadMain(in *os.File, r Callback, fnm string, a []interface{}) {
	reader := csv.NewReader(transform.NewReader(in, japanese.ShiftJIS.NewDecoder()))
	//行ごとのフィールド数が違ってもOKにする。
	reader.FieldsPerRecord = -1
	var sgn, sort, lineno, edano, lnoeda string
	var rece [][]string
	linenum := 0
	receIdx := make(map[string]int)
	eda := make(map[string]int)
	hosei := make(map[string][]string)
	flgJy := false
	for {
		arr, err := reader.Read()
		if err == io.EOF {
			break
		}
		linenum++
		if len(arr) < 3 {
			continue
		}
		if linenum%common.TickLineNum == 0 {
			fmt.Printf("[loading %d lines...]\n", linenum)
		}
		sgn = arr[locSgn]
		sort = arr[locSort]
		lineno = arr[locLineno]
		edano = arr[locEda]
		lnoeda = strings.Join([]string{lineno, edano}, common.Collon)

		if sort == "MN" {
			r(receSorting(rece, flgJy), fnm, a)
			rece = nil
			receIdx = make(map[string]int)
			eda = make(map[string]int)
			hosei = make(map[string][]string)
			flgJy = false
		}
		if (sgn == "1") || (sort == "MN") {
			rece = append(rece, arr)
			receIdx[lnoeda] = len(rece) - 1
			eda[lnoeda] = common.Max(eda[lnoeda], common.Atoi(edano, 0))
		} else if (sgn == "2") || (sort != "JY") {
			hosei[lnoeda] = arr
		} else if (sgn == "2") || (sort == "JY") {
			flgJy = true
			jyOp := arr[locJYOp]
			oblno := arr[locJYObLno]
			obeda := arr[locJYObEda]
			oblnoeda := strings.Join([]string{oblno, obeda}, common.Collon)
			nwlnoeda := strings.Join([]string{arr[locJYNwLno], arr[locJYNwEda]}, common.Collon)
			if jyOp == "1" { //add
				nwArr := hosei[nwlnoeda]
				nwArr[locSgn] = "1"
				nwArr[locLineno] = oblno
				eda[oblnoeda] = eda[oblnoeda] + 1
				nwArr[locEda] = strconv.Itoa(eda[oblnoeda])
				rece = append(rece, nwArr)
			} else if jyOp == "2" { //replace
				nwArr := hosei[nwlnoeda]
				nwArr[locSgn] = "1"
				nwArr[locLineno] = oblno
				nwArr[locEda] = obeda
				rece = append(rece, nwArr)
				rece[receIdx[oblnoeda]] = []string{common.Empty}
			} else if jyOp == "3" { //delete
				rece[receIdx[oblnoeda]] = []string{common.Empty}
			}
		}
	}
	if len(rece) > 0 {
		r(receSorting(rece, flgJy), fnm, a)
	}
	fmt.Printf("[loaded %d lines.]\n", linenum)
}

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

//Load -- レセファイル１つのロード
func Load(fnm string, r Callback, a []interface{}) {
	fmt.Printf("[target:%s start..]\n", fnm)
	f, err := os.Open(fnm)
	failOnError(err)
	defer f.Close()
	loadMain(f, r, fnm, a)
}

//LoadArr -- レセファイル群を 一気にロード
func LoadArr(fnms []string, r Callback, a []interface{}) {
	sort.Sort(sort.Reverse(sort.StringSlice(fnms)))
	for _, fnm := range fnms {
		Load(fnm, r, a)
	}
}
