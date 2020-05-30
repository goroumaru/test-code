package orDoneChannel_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/goroumaru/test-code/orDoneChannel"
)

func TestOrDone(t *testing.T) {

	// ここでやっていること
	// 定時実行した結果をmychannelとして渡す。
	// このとき、orDoneを利用する。
	// 時限処置によりタイムアウトして、すべての処理を終える。

	// doneチャンネルをクローズするために、コンテキストでタイムアウトする
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 定時実行
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	// メインゴルーチンが先に終了してしまうので、子ゴルーチンを待たせる
	wg := sync.WaitGroup{}
	wg.Add(1)
	myChannel := make(chan interface{})
	defer close(myChannel)

	done := make(chan interface{})
	go func() {
		defer close(done) // gorutineから抜けるとき、doneチャンネルも閉じられる
		defer wg.Done()
		for {
			select {
			case <-ctx.Done(): // 時限でcontextがクローズする
				fmt.Println("context is closed!")
				return
			case <-tick.C: // 定時実行
				myChannel <- "my channel!"
			}
		}
	}()

	// orDoneのゴルーチンは、doneチャンネルが送信される(閉じれれる)まで終了しない。
	for val := range orDoneChannel.OrDone(done, myChannel) {
		fmt.Printf("valに対して何かするところ: %v\n", val)
	}

	// 子ゴルーチンを待機する
	wg.Wait()
}

func TestOrDoneCtx(t *testing.T) {

	// ここでやっていること
	// 定時実行した結果をmychannelとして渡す。
	// このとき、orDoneを利用する。
	// 時限処置によりタイムアウトして、すべての処理を終える。

	// doneチャンネルをクローズするために、コンテキストでタイムアウトする
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 定時実行
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	// メインゴルーチンが先に終了してしまうので、子ゴルーチンを待たせる
	wg := sync.WaitGroup{}
	wg.Add(1)
	myChannel := make(chan interface{})
	defer close(myChannel)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done(): // 時限でcontextがクローズする
				fmt.Println("context is closed!")
				return
			case <-tick.C: // 定時実行
				myChannel <- "my channel!"
			}
		}
	}()

	// doneからctxへ変更
	childCtx, childCancel := context.WithCancel(ctx)
	defer childCancel()

	// contextがキャンセルされるまで終了しない。
	for val := range orDoneChannel.OrDoneCtx(childCtx, myChannel) {
		fmt.Printf("valに対して何かするところ: %v\n", val)
	}

	// 子ゴルーチンを待機する
	wg.Wait()
}
