package utils

import (
	"math"
	"strconv"
)

// Float64ToString : convert float64 to string.
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 10, 64)
}

// ToFloat64 : string to float64
func ToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// Round : 四捨五入
func Round(num, places float64) float64 {
	shift := math.Pow(10, places)
	return roundInt(num*shift) / shift
}

// Ceil : 切り上げ
func Ceil(num, places float64) float64 {
	shift := math.Pow(10, places)
	return roundUpInt(num*shift) / shift
}

// Floor : 切り捨て
func Floor(num, places float64) float64 {
	shift := math.Pow(10, places)
	return math.Trunc(num*shift) / shift
}

// Trunc : 切り捨て(小数部)
func Trunc(num float64) float64 {
	return math.Trunc(num)
}

// roundInt : 四捨五入(整数)
func roundInt(num float64) float64 {
	t := math.Trunc(num)
	if math.Abs(num-t) >= 0.5 {
		return t + math.Copysign(1, num)
	}
	return t
}

// roundInt : 切り上げ(整数)
func roundUpInt(num float64) float64 {
	t := math.Trunc(num)
	return t + math.Copysign(1, num)
}
