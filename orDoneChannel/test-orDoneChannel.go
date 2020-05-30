package orDoneChannel

import "context"

// OrDone :
func OrDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})

	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-done:
					// <-doneがないとdoneチャネルが送信されてきても、
					// valStream <- vが送信されてくるまでブロックされ続けてしまう。
					// ここで<-doneとなると、1つ上のネストにおける<-doneで抜けられる。
				}
			}
		}
	}()
	return valStream
}

// OrDoneCtx : doneの代わりにcontextを使う
func OrDoneCtx(ctx context.Context, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})

	go func() {
		defer close(valStream)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-ctx.Done():
					// <-ctx.Done()がないとdoneチャネルが送信されてきても、
					// valStream <- vが送信されてくるまでブロックされ続けてしまう。
					// ここで<-ctx.Done()となると、1つ上のネストにおける<-ctx.Done()で抜けられる。
				}
			}
		}
	}()
	return valStream
}
