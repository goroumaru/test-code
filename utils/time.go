package utils

import (
	"strconv"
	"time"
)

const (
	one   string = "01"
	two   string = "02"
	three string = "03"
	four  string = "04"
	five  string = "05"
	six   string = "06"
	seven string = "07"
	eight string = "08"
	nine  string = "09"
)

// ConvertFormatDate convers date to strings.
func ConvertFormatDate(now time.Time) (year, month, day string) {
	y, m, d := now.Date()
	year = strconv.Itoa(y)
	month = convertToString(int(m))
	day = convertToString(d)
	return year, month, day
}

func convertToString(value int) (str string) {
	switch value {
	case 1:
		str = one
	case 2:
		str = two
	case 3:
		str = three
	case 4:
		str = four
	case 5:
		str = five
	case 6:
		str = six
	case 7:
		str = seven
	case 8:
		str = eight
	case 9:
		str = nine
	default:
		str = strconv.Itoa(value)
	}
	return str
}
