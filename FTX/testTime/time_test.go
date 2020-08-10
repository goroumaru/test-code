package testTime

import (
	"fmt"
	"testing"
	"time"
)

func TestClock(t *testing.T) {
	now := time.Now()
	fmt.Println(now)
	h, m, s := now.Clock()
	fmt.Printf("%v,%v,%v\n", h, m, s)
}

func TestJustTime(t *testing.T) {
	for {
		if _, _, s := time.Now().Clock(); s == 0 {
			fmt.Printf("just time second = %v\n", s)
			break
		}
	}
}
