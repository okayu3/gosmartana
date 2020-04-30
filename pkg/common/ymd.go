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
	if (err != nil) || (i < 1010101) || (i > 9999999) {
		return Empty
	}
	if i > 5000000 {
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
	if (err != nil) || (i < 10101) || (i > 99999) {
		return Empty
	}
	if i > 50000 {
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
