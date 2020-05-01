package common

import (
	"math"
	"strconv"
)

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
