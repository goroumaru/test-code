package teeChannel

import (
	"context"

	"github.com/goroumaru/test-code/orDoneChannel"
)

// Tee :
func Tee(ctx context.Context, in <-chan interface{}) (_, _ <-chan interface{}) { // 戻り値ってこんな書き方もできるんだ！
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)
		for val := range orDoneChannel.OrDoneCtx(ctx, in) {
			var out1, out2 = out1, out2 // コピーしてローカル変数用意
			for i := 0; i < 2; i++ {    // out1とout2を確実に選択するため。
				select {
				case out1 <- val:
					out1 = nil // コピー側へnil代入し、out1チャンネルをブロックさせる。(=out2を選択させる)
				case out2 <- val:
					out2 = nil // コピー側へnil代入し、out2チャンネルをブロックさせる。(=out1を選択させる)
				}
			}
			// out1とout2の書き込みが終わると、inチャンネルが読み込み可能となる。
			// for ~ rangeのイテレーションがひとつ進むから。
		}
	}()
	return out1, out2
}
