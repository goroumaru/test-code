package pipeline

import (
	"context"
	"fmt"
)

// Generator : 入力データをストレームへ変換する
func Generator(ctx context.Context, integers ...int) <-chan int { // ...int : int型可変スライス
	intStream := make(chan int, len(integers)) // パイプライン入口であり、バッファチャンネルとする
	go func() {
		defer close(intStream)
		for idx, i := range integers {
			select {
			case intStream <- i: // intStreamがFullだとブロックする。つまり、イテレーション停止。
				fmt.Printf("Generator[%v]: %v\n", idx, i)
			case <-ctx.Done(): // ゴルーチンリークを防ぐ
				return
			}
		}
	}()
	return intStream
}

// Multiply : 単純な掛け算
func Multiply(ctx context.Context, intStream <-chan int, multiplier int) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer close(multipliedStream)
		var idx int
		for i := range intStream { // chan型なので戻り値ひとつ。idxはない。
			select {
			case multipliedStream <- i * multiplier: // multipliedStreamがFullだとブロックする。つまり、イテレーション停止。
				fmt.Printf("Multiply[%v]: %v\n", idx, i*multiplier)
				idx++
			case <-ctx.Done(): // ゴルーチンリークを防ぐ
				return
			}
		}
	}()
	return multipliedStream
}

// Add : 単純な加算
func Add(ctx context.Context, intStream <-chan int, additive int) <-chan int {
	addedStream := make(chan int)
	go func() {
		defer close(addedStream)
		var idx int
		for i := range intStream {
			select {
			case addedStream <- i + additive: // addedStreamがFullだとブロックする。つまり、イテレーション停止。
				fmt.Printf("Multiply[%v]: %v\n", idx, i+additive)
				idx++
			case <-ctx.Done(): // ゴルーチンリークを防ぐ
				return
			}
		}
	}()
	return addedStream
}
