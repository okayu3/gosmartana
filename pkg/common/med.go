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
