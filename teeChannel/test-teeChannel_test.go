package teeChannel_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/goroumaru/test-code/teeChannel"
)

func TestTeeChannel(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 定時実行
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	inData := make(chan interface{})
	go func() {
		defer close(inData)
		var cnt int
		for {
			select {
			case <-tick.C:
				cnt++
				inData <- cnt
			case <-ctx.Done():
				return
			}
		}
	}()

	// teeチャンネルを利用してチャンネルを分岐する
	out1, out2 := teeChannel.Tee(ctx, inData)

	for v1 := range out1 {
		v2 := <-out2
		fmt.Printf("２つに分割:\n(out1, ok) = (%v)\n(out2, ok) = (%v)\n", v1, v2)
	}
}
