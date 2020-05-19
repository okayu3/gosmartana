package common

import (
	"math"
	"strconv"
)

var bankPow10 = [11]float64{1.0, 10.0, 100.0, 1000.0, 10000.0, 100000.0, 1000000.0, 10000000.0, 100000000.0, 1000000000.0, 10000000000.0}

//Min -- minimum integer from integers
func Min(nums ...int) int {
	if len(nums) == 0 {
		panic("funciton min() requires at least one argument.")
	}
	res := nums[0]
	for i := 0; i < len(nums); i++ {
		res = int(math.Min(float64(res), float64(nums[i])))
	}
	return res
}

//Max -- maximum integer from integers
func Max(nums ...int) int {
	if len(nums) == 0 {
		panic("funciton max() requires at least one argument.")
	}
	res := nums[0]
	for i := 0; i < len(nums); i++ {
		res = int(math.Max(float64(res), float64(nums[i])))
	}
	return res
}

//Atoi -- atoi
func Atoi(s string, defaultvalue int) int {
	ans, err := strconv.Atoi(s)
	if err != nil {
		return defaultvalue
	}
	return ans
}

// RoundI : 小数点以下四捨五入
func RoundI(f float64) int {
	return int(math.Floor(f + 0.5))
}

// Round : 任意小数点位置四捨五入
func Round(f float64, places int) float64 {
	shift := bankPow10[places]
	return math.Floor(f*shift+.5) / shift
}

// Ceil : 任意小数点位置切り上げ
func Ceil(f float64, places int) float64 {
	shift0 := bankPow10[places+1]
	ff := math.Floor(f*shift0) / shift0
	shift := bankPow10[places]
	wk := ff * shift
	if ff > 0 {
		if wk == math.Floor(wk) {
			return wk / shift
		}
		return math.Floor(wk+1) / shift
	}
	return math.Floor(wk) / shift
}

// Round5sha 五捨五超入
func Round5sha(f float64) float64 {
	ans := math.Trunc(f)
	if f-ans > 0.5 {
		ans = ans + 1.0
	}
	return ans
}
