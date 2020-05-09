package ana

import (
	"fmt"
	"os"
	"strings"

	"github.com/okayu3/gosmartana/pkg/ckey"
	"github.com/okayu3/gosmartana/pkg/common"
	"github.com/okayu3/gosmartana/pkg/rece"
)

//MakeExpenseMed -- 医科レセプトからExpense作成
func MakeExpenseMed(one [][]string, fnm string, args []interface{}) int {
	var arr []string
	var sort, mnKensaku, ymdB, gend, name, kananame string
	var prf, pnt, ircd, irnm, seikyuYm string
	var sinryoYm, innDate, shubetu string
	var jitsuDates, ten string
	var ck, iKi, iBan string

	if one == nil {
		return -1
	}
	outHandle := args[0].(*os.File)
	flgHO := 0

	kbn := "医科"

	for _, arr = range one {
		if (arr == nil) || (len(arr) <= 0) {
			continue
		}
		sort = arr[rece.RSort]
		if sort == "MN" {
			mnKensaku = arr[rece.RMNkensaku]
		} else if sort == "IR" {
			prf = arr[rece.RIRprf]
			pnt = arr[rece.RIRpnttbl]
			ircd = arr[rece.RIRircode]
			irnm = arr[rece.RIRirname]
			seikyuYm = common.YmW2g(arr[rece.RIRseikyuYm])
		} else if sort == "YK" {
			prf = arr[rece.RYKprf]
			pnt = arr[rece.RYKpnttbl]
			ircd = arr[rece.RYKykcode]
			irnm = arr[rece.RYKykname]
			seikyuYm = common.YmW2g(arr[rece.RYKseikyuYm])
		} else if sort == "RE" {
			//共通(のちにrefactoringすること)
			if mnKensaku == common.Empty {
				mnKensaku = arr[rece.RREkensaku]
			}
			ymdB = common.YmdW2g(arr[rece.RREbirth])
			gend = arr[rece.RREgender]
			name = arr[rece.RREname]
			if len(arr) > rece.RREkananame {
				kananame = arr[rece.RREkananame]
			}
			//独自
			sinryoYm = common.YmW2g(arr[rece.RREsinryoYm])
			innDate = common.YmdW2g(arr[rece.RREinnDate])
			shubetu = arr[rece.RREshubetu]
		} else if sort == "HO" {
			flgHO = 1
			//共通(のちにrefactoringすること)
			iKi = arr[rece.RHOinsKigo]
			iBan = arr[rece.RHOinsBango]
			ck = ckey.Get(iKi, iBan, ymdB, gend, name, kananame)
			//独自
			jitsuDates = arr[rece.RHOjitudate]
			ten = arr[rece.RHOten]
		} else if sort == "HO" {
			if flgHO == 0 {
				//共通(のちにrefactoringすること)
				iKi = arr[rece.RKOftn]
				iBan = arr[rece.RKOrcv]
				ck = ckey.Get(iKi, iBan, ymdB, gend, name, kananame)
				//独自
				jitsuDates = arr[rece.RKOjitudate]
				ten = arr[rece.RKOten]
			}
		}
	}
	opSaveExpense(outHandle, ck, mnKensaku, gend, ymdB, shubetu,
		jitsuDates, ten, kbn, sinryoYm, innDate,
		prf, pnt, ircd, irnm, seikyuYm)
	return 1
}

func opSaveExpense(outHandle *os.File, ck, mnKensaku,
	gend, ymdB, shubetu,
	jitsuDates, ten, kbn, sinryoYm, innDate,
	prf, pnt, ircd, irnm, seikyuYm string) {

	nyugai := "2"
	honnin := "1"
	if len(shubetu) == 0 {
		fmt.Println("shubetu not set on" + mnKensaku)
	}
	wk := int(shubetu[3] - '0')
	if wk%2 == 1 {
		nyugai = "1"
	}
	if (wk != 1) && (wk != 2) {
		honnin = "2"
	}
	oneExpense := strings.Join([]string{ck, mnKensaku,
		gend, ymdB, shubetu, nyugai, honnin,
		jitsuDates, ten, kbn, sinryoYm, innDate,
		prf, pnt, ircd, irnm, seikyuYm},
		common.Comma)
	outHandle.WriteString(oneExpense + "\n")

}
