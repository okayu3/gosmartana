package common

import "strings"

//IsDoubtDisease -- 疑い病名を判定
func IsDoubtDisease(s string) bool {
	if s == Empty {
		return false
	}
	len := len(s)
	w := Empty
	i := 0
	for i = 0; i < len-3; i += 4 {
		w += s[i:i+4] + ","
	}
	judge := strings.Contains(w, "8002")
	return judge
}

//DivAffix -- 修飾語を 接頭語と 接尾語に分ける
func DivAffix(affix string) (string, string) {
	if affix == Empty {
		return Empty, Empty
	}
	len := len(affix)
	var prefix, suffix, one string
	for i := 0; i < len-3; i += 4 {
		one = affix[i : i+4]
		if one[0:1] == "8" {
			suffix += one
		} else {
			prefix += one
		}
	}
	return prefix, suffix
}

//IsLongCareRece -- 長期レセ判定
func IsLongCareRece(s string) bool {
	if s == Empty {
		return false
	}
	len := len(s)
	w := Empty
	i := 0
	for i = 0; i < len-1; i += 2 {
		w += s[i:i+2] + ","
	}
	judge := strings.Contains(w, "02")
	judge = judge || strings.Contains(w, "06")
	judge = judge || strings.Contains(w, "16")
	return judge
}
