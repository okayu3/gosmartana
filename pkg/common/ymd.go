package common

import (
	"strconv"
	"strings"
)

//YmdW2g -- YMD Wareki to Gregorius reki
func YmdW2g(s string) string {
	s = strings.TrimSpace(s)
	if s == Empty {
		return s
	}
	var i int
	var err error
	i, err = strconv.Atoi(s)
	if (err != nil) || (i < 1010101) {
		return Empty
	}
	if i > 9999999 {
		//初めから西暦のYMDで 来ている
		return s
	} else if i > 5000000 {
		return strconv.Itoa(i - 5000000 + 20180000)
	} else if i > 4000000 {
		return strconv.Itoa(i - 4000000 + 19880000)
	} else if i > 3000000 {
		return strconv.Itoa(i - 3000000 + 19250000)
	} else if i > 2000000 {
		return strconv.Itoa(i - 2000000 + 19110000)
	} else {
		return strconv.Itoa(i - 1000000 + 18670000)
	}
}

//YmW2g -- YM Wareki to Gregorius reki
func YmW2g(s string) string {
	s = strings.TrimSpace(s)
	if s == Empty {
		return s
	}
	var i int
	var err error
	i, err = strconv.Atoi(s)
	if (err != nil) || (i < 10101) {
		return Empty
	}
	if i > 99999 {
		//初めから西暦のYMで 来ている
		return s
	} else if i > 50000 {
		return strconv.Itoa(i - 50000 + 201800)
	} else if i > 40000 {
		return strconv.Itoa(i - 40000 + 198800)
	} else if i > 30000 {
		return strconv.Itoa(i - 30000 + 192500)
	} else if i > 20000 {
		return strconv.Itoa(i - 20000 + 191100)
	} else {
		return strconv.Itoa(i - 10000 + 186700)
	}
}

//AgeAt -- atymdの段階での生年月日ymd_b の年齢
func AgeAt(ymdB, atYmd string) int {
	var i, j int
	var err0, err1 error
	i, err0 = strconv.Atoi(ymdB)
	j, err1 = strconv.Atoi(atYmd)
	if (err0 != nil) || (err1 != nil) {
		return -1
	}
	age := (j - i) / 10000
	return age
}

//AnnualAtYmd -- atymd の年度。ただし年度開始を 4/1とする
func AnnualAtYmd(atYmd string) int {
	var i int
	var err0 error
	i, err0 = strconv.Atoi(atYmd)
	if err0 != nil {
		return -1
	}
	//year
	y := i / 10000
	//month
	m := (((i/100)%100 - 1) % 12) + 1
	if m < 4 {
		y--
	}
	return y
}

//AnnualAtYm -- atym の年度。ただし年度開始を 4/1とする
func AnnualAtYm(atYm string) int {
	var i int
	var err0 error
	i, err0 = strconv.Atoi(atYm)
	if err0 != nil {
		return -1
	}
	//year
	y := i / 100
	//month
	m := ((i%100 - 1) % 12) + 1
	if m < 4 {
		y--
	}
	return y
}
