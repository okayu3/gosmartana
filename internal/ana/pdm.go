package ana

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/okayu3/gosmartana/pkg/common"
	"github.com/okayu3/gosmartana/pkg/rece"
)

var dicHst = make(map[string]([]float64))
var bankLn [300]float64

var arrWriter [3]*bufio.Writer

//const -- ファイル名など
const (
	maxIcd10Codes = 150
	PdmDataMale   = "pdm_icd10s_male.csv"
	PdmDataFemale = "pdm_icd10s_female.csv"
)

//RunPDM -- PDM の実行
func RunPDM(fnm, mstdir, outDir string) {
	iniPDM()
	loadPDM(fnm)
	calcDivRatio()
	tstmp := time.Now().Format("20060102150405")
	divideFromData(fnm, tstmp, outDir)
	opSummary(fnm, tstmp, mstdir, outDir)

}

func iniPDM() {
	for i := 1; i < 300; i++ {
		bankLn[i] = math.Log(float64(i))
	}
	dicHst["0000"] = make([]float64, 12)
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
	if days == 0 {
		days = 1
	}
	//vLn := math.Log(float64(cntDis))
	vLn := bankLn[cntDis]
	var cd119 string
	pBydLn := points / (days * (1 + vLn))
	daysP2 := days * days
	//pBydLnP2 := common.Round(pBydLn,3) * common.Round(pBydLn,3)
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

func divideFromData(fnm, tstmp, outDir string) {
	_, fileName := filepath.Split(fnm)
	oPoints := outDir + "pdm/result/ATensu_" + fileName
	oDates := outDir + "pdm/result/ANissu_" + fileName

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
		//divP[i-3] = common.RoundI(points * dicHst[cd119][6] / sumMeanP)
		divP[i-3] = common.RoundI(common.Round(points*dicHst[cd119][6]/sumMeanP, 10))

		//0.349999999 などが 0.3になってしまっていたので、いったん10ケタ目で四捨五入してからにしてみる。
		//divD[i-3] = common.Round(dates*dicHst[cd119][7]/sumMeanD, 1)
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
		// 				common.Round(dates*dicHst[cd119][7]/sumMeanD, 1)))
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

func opSummary(dataFnm, tstmp, mstdir, outDir string) {
	_, fileName := filepath.Split(dataFnm)
	//	ofnm := basedir + "result/APDM_RESULT_" + tstmp + ".csv"
	ofnm := outDir + "pdm/result/AResult_" + fileName
	ofp, err := os.Create(ofnm)
	failOnError(err)
	defer ofp.Close()
	w := bufio.NewWriter(ofp)
	mValue := dicHst["0000"][2] / dicHst["0000"][5]
	w.WriteString(fmt.Sprintf("ファイル,%s,点数合計,%.0f,重み計算,4,平均値Ｍ,%.11f,補正次数,自動\r\n", fileName, dicHst["0000"][2], mValue))
	header := "傷病コード,傷病名,傷病数,Σ日数,ΣＰ／Ｄ(１＋ＬＮ(Ｎ)),AV日数,SD日数,AV点数,SD点数,Σ配分日数,Σ配分点数,日数割合,点数割合,(重み)AV日数,(重み)１日当り点数,(重み)Σ配分日数,(重み)Σ配分点数,(重み)日数割合,(重み)点数割合"
	w.WriteString(fmt.Sprintf("%s\r\n", header))

	for k, v := range dicHst {
		if k != "0000" {
			dicHst["0000"][11] += v[11]
			dicHst["0000"][10] += v[10]
		}
	}

	arrWriter[2] = w
	common.LoadCSV(mstdir+"/cd119.txt", opSummaryByLine, common.ModeCsvUTF8)
	w.Flush()
	ofp.Close()
}

func opSummaryByLine(arr []string, lineno int) {
	cd119 := arr[0]
	catName := arr[1]
	w := arrWriter[2]
	if _, ok := dicHst[cd119]; ok == false {
		w.WriteString(fmt.Sprintf("%s,%s,%s\r\n", cd119, catName, "0,0,0,0,0.00,0,0.00,0,0,0,0,0,0.000,0.000,0.000,0.000,0.000"))
		return
	}
	a := dicHst[cd119]
	cnt := a[0]
	sumDays := a[1]
	//sumPoints := a[2]
	//sumPByD:= a[3]
	sumPByDVLN := a[4]
	//sumCDisNum := a[5]
	meanP := a[6]
	meanD := a[7]
	sumP2Days := a[8]
	sumP2PByDVLN := a[9]
	//	sdDays := math.Pow(sumP2Days/cnt-math.Pow(sumDays/cnt, 2), 0.5)
	//	sdPbyDlnN := math.Pow(sumP2PByDVLN/cnt-math.Pow(sumPByDVLN/cnt, 2), 0.5)
	var sdPbyDlnN, sdDays float64
	if cnt != 1 {
		//標本標準偏差で行っている模様。
		// 通常の分散：Sum(X^2)/N  - (Sum(X)/N)^2
		// 標本分散：  Sum(X^2/N-1 - (Sum(X))^2/N/N-1
		sdDays = math.Pow((sumP2Days/(cnt-1) - math.Pow(sumDays, 2)/cnt/(cnt-1)), 0.5)
		sdPbyDlnN = math.Pow((sumP2PByDVLN/(cnt-1) - math.Pow(sumPByDVLN, 2)/cnt/(cnt-1)), 0.5)
		//sdDays = math.Pow((cnt/(cnt-1))*(sumP2Days/cnt-math.Pow(sumDays/cnt, 2)), 0.5)
		//sdPbyDlnN = math.Pow((cnt/(cnt-1))*(sumP2PByDVLN/cnt-math.Pow(sumPByDVLN/cnt, 2)), 0.5)
	} else {
		sdDays = 0.0
		sdPbyDlnN = 0.0
	}
	if math.IsNaN(sdDays) {
		sdDays = 0.0
	}
	if math.IsNaN(sdPbyDlnN) {
		sdPbyDlnN = 0.0
	}

	// 	if (cd119 == "0101") || (cd119 == "1403") {
	// 		echo(fmt.Sprintf("%s\t%.3f\t%f\t%f\t%f\r\n", cd119, sdPbyDlnN, sumP2PByDVLN, cnt, a[5]))
	// 	}

	sumDivDays := a[11]
	sumDivPoints := a[10]
	// 	ratioDivDays := common.Round(100*a[11]/dicHst["0000"][1], 3)
	// 	ratioDivPoints := common.Round(100*a[10]/dicHst["0000"][2], 3)
	ratioDivDays := common.Round(100*a[11]/dicHst["0000"][11], 3)
	ratioDivPoints := common.Round(100*a[10]/dicHst["0000"][10], 3)
	w.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\r\n",
		cd119,
		catName,
		f2a(cnt, 0),
		f2a(sumDays, 0),
		f2a(sumPByDVLN, 3),
		f2a(meanD, 1),
		fmt.Sprintf("%.2f", common.Round(sdDays, 2)),
		f2a(meanP, 3),
		fmt.Sprintf("%.2f", common.Round(sdPbyDlnN, 2)),
		f2a(sumDivDays, 1),
		f2a(sumDivPoints, 0),
		f2a(ratioDivDays, 3),
		f2a(ratioDivPoints, 3),
		"0,0.000,0.000,0.000,0.000,0.000"))

}

func f2a(v float64, digit int) string {
	//	if (math.IsNaN(v)) || (v == 0.0) {
	if math.IsNaN(v) {
		wk := fmt.Sprintf("%."+strconv.Itoa(digit)+"f", 0.0)
		//		wk := "0.00"
		return wk
	}
	vv := common.Round(v, digit)
	wk := fmt.Sprintf("%."+strconv.Itoa(digit)+"f", vv)
	if digit > 0 {
		wk = strings.TrimRight(wk, "0")
		wk = strings.TrimRight(wk, ".")
	}
	return wk
}

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}

//PreparePDMData -- PDM用データの作成ファイルハンドルの作成
func PreparePDMData(outDir string) [](*os.File) {
	path := outDir + "pdm"
	common.MakeDir(path)
	common.MakeDir(path + "/result")
	fnms := []string{path + "/pdm_disseq.csv",
		path + "/" + PdmDataMale,
		path + "/" + PdmDataFemale,
	}
	var handles [](*os.File)
	for _, fnm := range fnms {
		ofile, _ := os.OpenFile(fnm, os.O_WRONLY|os.O_CREATE, 0666)
		handles = append(handles, ofile)
	}
	return handles
}

//ClosePDMHandle -- PDMData 作成用ファイルハンドルのclosing
func ClosePDMHandle(a [](*os.File)) {
	for _, oHandle := range a {
		oHandle.Close()
	}
}

//OpSavePDMData -- PDM用データの保存 (aOutHandlesPDM, mnKensaku, jitsuDates, ten, gend, aCnt, aIcd10)
func OpSavePDMData(aOutHandlesPDM [](*os.File), mnKensaku, jitsuDates, ten, gend string, aCnt []string, aIcd10 []string) {
	oneSeq := strings.Join(append([]string{mnKensaku, jitsuDates, ten}, aCnt...), common.Comma)
	oneIcd10 := strings.Join(append([]string{mnKensaku, jitsuDates, ten}, aIcd10...), common.Comma)
	outHandleSeq := aOutHandlesPDM[0]
	outHandleSeq.WriteString(oneSeq + "\n")
	var outHandleIcd10 *os.File
	if gend == "2" {
		outHandleIcd10 = aOutHandlesPDM[2]
	} else {
		outHandleIcd10 = aOutHandlesPDM[1]
	}
	outHandleIcd10.WriteString(oneIcd10 + "\n")
}

//OpAfterPDM -- operation after PDM
func OpAfterPDM(outdir string) {
	fnmDisseq := outdir + "pdm/pdm_disseq.csv"
	dicSeq := make(map[string]string)
	common.LoadByLine(fnmDisseq, func(s string, lineno int) {
		a := strings.SplitN(s, ",", 4)
		if len(a) == 4 {
			mnKensaku := a[0]
			dicSeq[mnKensaku] = a[3]
		}
	}, common.ModeSJIS)
	fnmsTensu := common.ListUpFiles(outdir+"pdm/result/", "ATensu_", ".csv")
	dicGakuSV := make(map[string]string)
	common.LoadCSVArr(fnmsTensu, func(arr []string, lineno int) {
		var ten string
		mnKensaku := arr[0]
		seqs := strings.Split(dicSeq[mnKensaku], ",")
		for i := 0; i < len(seqs); i++ {
			ten = arr[3+i]
			dicGakuSV[mnKensaku+common.Collon+seqs[i]] =
				fmt.Sprintf("%d", common.Atoi(ten, 0)*rece.GakuForTen)
		}
	}, common.ModeCsvSJIS)
	fnmDisease := outdir + "disease.csv"
	fnmDiseaseOut := outdir + "diseasePdm.csv"
	ofileDiseaseOut, _ := os.OpenFile(fnmDiseaseOut, os.O_WRONLY|os.O_CREATE, 0666)
	defer ofileDiseaseOut.Close()
	common.LoadCSV(fnmDisease, func(arr []string, lineno int) {
		mnKensaku := arr[1]
		seq := arr[2]
		kk := mnKensaku + common.Collon + seq
		arr = append(arr, dicGakuSV[kk])
		oneline := strings.Join(arr, common.Comma)
		ofileDiseaseOut.WriteString(oneline + "\n")
	}, common.ModeCsvUTF8)
}
