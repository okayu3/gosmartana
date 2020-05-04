package ana

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/okayu3/gosmartana/pkg/common"
)

var dicHst = make(map[string]([]float64))
var bankLn [300]float64

var basedir = "C:/works/data/GoProj/APDM/"
var arrWriter [3]*bufio.Writer

const (
	maxIcd10Codes = 150
)

//RunPDM -- PDM の実行
func RunPDM(fnm string) {
	iniPDM()
	loadPDM(fnm)
	calcDivRatio()

}

func iniPDM() {
	for i := 1; i < 300; i++ {
		bankLn[i] = math.Log(float64(i))
	}
}

func loadPDM(fnm string) {
	common.LoadCSV(fnm, opLoadPDM, common.ModeCsvSJIS)
}

func opLoadPDM(arr []string, lineno int) {
	days, _ := strconv.ParseFloat(arr[1], 64)
	points, _ := strconv.ParseFloat(arr[2], 64)
	cntDis := len(arr) - 3
	if cntDis == 0 {
		return
	}
	//vLn := math.Log(float64(cntDis))
	vLn := bankLn[cntDis]
	var cd119 string
	pBydLn := points / (days * (1 + vLn))
	daysP2 := days * days
	//pBydLnP2 := Round(pBydLn,3) * Round(pBydLn,3)
	pBydLnP2 := pBydLn * pBydLn
	for i := 3; i < len(arr); i++ {
		if i-2 > maxIcd10Codes {
			break
		}
		cd119 = arr[i]
		if _, ok := dicHst[cd119]; ok == false {
			dicHst[cd119] = make([]float64, 12)
		}
		dicHst[cd119][0] += 1.0
		dicHst[cd119][1] += days
		dicHst[cd119][2] += points
		dicHst[cd119][3] += points / days
		//dicHst[cd119][4] += points / (days * (1 + vLn))
		dicHst[cd119][4] += common.Round(pBydLn, 3)
		dicHst[cd119][5] += float64(cntDis)
		dicHst[cd119][8] += daysP2
		dicHst[cd119][9] += pBydLnP2
		dicHst["0000"][0] += 1.0
	}
	cd119 = "0000"
	dicHst[cd119][1] += days
	dicHst[cd119][2] += points
	dicHst[cd119][3] += points / days
	dicHst[cd119][4] += pBydLn
	dicHst[cd119][5] += days * (1 + vLn)

}

func calcDivRatio() {
	var candMeanPoint float64
	mValue := dicHst["0000"][2] / dicHst["0000"][5]
	for cd119 := range dicHst {
		meanPByDlnN := common.Round(dicHst[cd119][4]/dicHst[cd119][0], 1)
		meanComb := common.Round(dicHst[cd119][5]/dicHst[cd119][0], 3)
		if meanPByDlnN >= mValue {
			candMeanPoint = meanPByDlnN * math.Pow((meanPByDlnN/mValue), (math.Log(meanComb)))
		} else {
			candMeanPoint = meanPByDlnN * math.Pow((meanPByDlnN/mValue), (meanComb-1))
		}
		//dicHst[cd119][6] = candMeanPoint
		dicHst[cd119][6] = common.Ceil(candMeanPoint, 3)
		dicHst[cd119][7] = common.Round(dicHst[cd119][1]/dicHst[cd119][0], 1)
	}
}

func divideFromData(fnm, tstmp string) {
	_, fileName := filepath.Split(fnm)
	oPoints := basedir + "result/ATensu_" + fileName
	oDates := basedir + "result/ANissu_" + fileName

	fpPoints, err0 := os.Create(oPoints)
	failOnError(err0)
	defer fpPoints.Close()
	arrWriter[0] = bufio.NewWriter(fpPoints)

	fpDates, err0 := os.Create(oDates)
	failOnError(err0)
	defer fpDates.Close()
	arrWriter[1] = bufio.NewWriter(fpDates)

	common.LoadCSV(fnm, opDivideFromData, common.ModeCsvSJIS)

	arrWriter[0].Flush()
	fpPoints.Close()

	arrWriter[1].Flush()
	fpDates.Close()
}

func opDivideFromData(arr []string, lineno int) {
	wp := arrWriter[0]
	wd := arrWriter[1]
	dates, _ := strconv.ParseFloat(arr[1], 64)
	points, _ := strconv.ParseFloat(arr[2], 64)
	cntDis := len(arr) - 3
	if cntDis == 0 {
		return
	}
	var cd119 string
	var sumMeanP float64
	var sumMeanD float64
	ilen := len(arr)
	for i := 3; i < ilen; i++ {
		if i-2 > maxIcd10Codes {
			break
		}
		cd119 = arr[i]
		sumMeanP += dicHst[cd119][6]
		sumMeanD += dicHst[cd119][7]
	}

	var divP [maxIcd10Codes]int
	var divD [maxIcd10Codes]float64
	for i := 3; i < ilen; i++ {
		if i-2 > maxIcd10Codes {
			break
		}
		cd119 = arr[i]
		//4.49999999 などが 4になってしまっていたので、いったん10ケタ目で四捨五入してからにしてみる。
		//divP[i-3] = RoundI(points * dicHst[cd119][6] / sumMeanP)
		divP[i-3] = common.RoundI(common.Round(points*dicHst[cd119][6]/sumMeanP, 10))

		//0.349999999 などが 0.3になってしまっていたので、いったん10ケタ目で四捨五入してからにしてみる。
		//divD[i-3] = Round(dates*dicHst[cd119][7]/sumMeanD, 1)
		divD[i-3] = common.Round(dates*dicHst[cd119][7]/sumMeanD, 10)
		divD[i-3] = common.Round(divD[i-3], 1)
		dicHst[cd119][10] += float64(divP[i-3])
		dicHst[cd119][11] += divD[i-3]
		// 		if arr[0] == "12142604911039286" {
		// 			echo(fmt.Sprintf("%s\t%f\t%f\t%f\t%f\t%.10f\t%f\r\n",
		// 				cd119,
		// 				dates,
		// 				sumMeanD,
		// 				dicHst[cd119][7],
		// 				divD[i-3],
		// 				dates*dicHst[cd119][7]/sumMeanD,
		// 				Round(dates*dicHst[cd119][7]/sumMeanD, 1)))
		// 		}
	}
	var m0 = make([]byte, 0, 512)
	m0 = append(m0, arr[0]...)
	m0 = append(m0, ',')
	m0 = append(m0, arr[1]...)
	m0 = append(m0, ',')
	m0 = append(m0, arr[2]...)
	for _, v := range divP {
		m0 = append(m0, ',')
		s := strconv.Itoa(v)
		m0 = append(m0, s...)
	}
	fmt.Fprintf(wp, "%s\r\n", string(m0))
	m0 = make([]byte, 0, 512)
	m0 = append(m0, arr[0]...)
	m0 = append(m0, ',')
	m0 = append(m0, arr[1]...)
	m0 = append(m0, ',')
	m0 = append(m0, arr[2]...)
	var dsp string
	for _, v := range divD {
		m0 = append(m0, ',')
		if v == 0 {
			m0 = append(m0, '0')
		} else {
			if float64(int(v)) == v {
				dsp = fmt.Sprintf("%.0f", v)
			} else {
				dsp = fmt.Sprintf("%.1f", v)
			}
			m0 = append(m0, dsp...)
		}
	}
	fmt.Fprintf(wd, "%s\r\n", string(m0))
}

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}
