package objectOrientedLike

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	println("Note...")
	// 0.
	// OO言語の継承は、「コード再利用」「派生型による共用・置換、多様性」の2つを実現するときに利用する。
	// goでは、「コード再利用」はembededで実現し、「派生型による共用・置換、多様性」はinterfaceで実現する。
	// よって、目的と実現方法が1対1となっており、使用した意図が明確に伝わる。
	// 1.
	// goには、クラス概念はないが、オブジェクト指向言語と比較できるように、クラスという名称を使う
	// 2.
	// 基本的な動作を踏まえてながら、Embeded fieldとInterfaceを利用して、構造的部分型（なのか？）を確認していく

	code := m.Run()

	println("Summarized...")
	// 1. goのembedded fieldは、"has a"関係である(単独では、リスコフの置換原則を満たさないと考えて良い？)
	// 2. derivedメソッドとしてアクセスするとき、
	//	  同名メソッドであればDerivedメソッドが呼ばれ、Baseだけに存在する異名メソッドは、Baseメソッドとして呼ばれる。
	//    (委譲されている)
	os.Exit(code)
}

func Test1(t *testing.T) {
	// OK
	// コンストラクタを利用し、初期値のままインスタンスすると・・・

	depth := newDepth()
	fmt.Printf("depth: %+v\n", depth)
	// Result
	// depth: &{timestamp:0 bid:{price:0 amount:0}}
}

func Test2(t *testing.T) {
	// OK
	// コンストラクタを使わずに、インスタンスする。
	// そのとき、ベースクラスのメンバー"だけ"書き換えようとすると・・・

	depth := &Depth{
		timestamp: 1,
	}
	fmt.Printf("depth: %+v\n", depth)
	// Result
	// depth: &{timestamp:1 bid:{price:0 amount:0}} // baseのメンバーが書き換えられている
}

func Test3(t *testing.T) {
	// NG（派生クラスにおけるベースクラスは、"継承ではなく委譲"であるため）
	// 派生クラスのメンバーも書き換えようとすると・・・

	depth := &Depth{
		timestamp: 1,
		// price:     2, // ←　これがあるとコンパイルエラー
		// amount:    3, // ←　これがあるとコンパイルエラー
	}
	fmt.Printf("depth: %+v\n", depth)
	// Compile Error
	// cannot use promoted field bid.price in struct literal of type depth
}

func Test4(t *testing.T) {
	// OK
	// 派生クラスのメンバーを書き換えるには・・・

	depth := &Depth{
		timestamp: 1,
		// 注意点は、"Bid: "となっていて、普通のメンバーとしてではなく、
		// embeded fieldとして記述すること
		Bid: Bid{
			price:  2, // ベースクラスのメンバーとして、書き換えればOK
			amount: 3,
		},
	}
	fmt.Printf("depth: %+v\n", depth)
	// Result
	// depth: &{timestamp:1 bid:{price:2 amount:3}}
}

func Test5(t *testing.T) {
	// OK（オーバライドされているように見えるだけで、実際は委譲されてる）
	// 派生クラスのメソッドを利用すると・・・

	depth := &Depth{
		timestamp: 1,
		Bid: Bid{
			price:  2,
			amount: 3,
		},
	}
	fmt.Printf("depth: %+v\n", depth)
	// Result
	// depth: &{timestamp:1 bid:{price:2 amount:3}}

	// derivedメソッドのようにアクセス
	fmt.Printf("base member: %+v\n", depth.timestamp)
	fmt.Printf("derived member: %+v, %+v\n", depth.price, depth.amount)
	fmt.Printf("derived method: %+v\n", depth.getExecutedPrice())
	// Result
	// base member: 1
	// derived member: 2, 3
	// This method is a derived structure. // ← 派生クラスのメソッドなので、オーバライドされてるように見えるが・・・
	// derived method: 6				   //　　あくまで委譲なので、継承されていない。(次のテストで確認する)

	// baseメソッドのようにアクセス
	fmt.Printf("base member: %+v\n", depth.timestamp)
	fmt.Printf("derived member: %+v, %+v\n", depth.Bid.price, depth.Bid.amount)
	fmt.Printf("derived method: %+v\n", depth.Bid.getExecutedPrice())
	// Result
	// base member: 1
	// derived member: 2, 3
	// This method is a base structure. // ← embeddedされたベースクラスのメソッドを利用している。
	// derived method: 6
}

func Test6(t *testing.T) {
	// OK（オーバライドされているように見えるだけで、実際は委譲されてる）
	// ベースクラスのメソッドを追加すると・・・
	// (ここで追加できないため、メソッド以外は同じクラスを利用する)

	depth := &DepthAsk{
		timestamp: 1,
		Ask: Ask{
			price:  2,
			amount: 3,
		},
	}
	fmt.Printf("depth: %+v\n", depth)
	// Result
	// depth: &{timestamp:1 bid:{price:2 amount:3}}

	// derivedメソッドのようにアクセス
	fmt.Printf("base member: %+v\n", depth.timestamp)
	fmt.Printf("derived member: %+v, %+v\n", depth.price, depth.amount)
	fmt.Printf("derived method: %+v\n", depth.getExecutedPrice())
	// Result
	// base member: 1
	// derived member: 2, 3
	// This method is a derived structure. // ← 派生クラスのメソッドだが、オーバライドされてるわけではない。
	// derived method: 6				   //　　派生クラスにあるメソッドは、そのメソッドが実行される。派生クラスにないメソッドは、ベースクラスのメソッドが実行される。

	fmt.Printf("derived method: %+v\n", depth.getAppliedFeePrice()) // 追加したベースメソッド
	// Result
	// This method is a base structure. // ← ベースクラスのメソッドを実行している。つまり、オーバーライドされていない。
	// derived method: -4

	// これは、コンパイラが以下と解釈して実行している
	fmt.Printf("derived method: %+v\n", depth.Ask.getAppliedFeePrice())
	// Result
	// This method is a base structure.
	// derived method: -4
}

func Test7(t *testing.T) {
	// NG（goのembeded fieldは、"is a"ではなく、"has a"関係であるため）
	// ベースクラスに派生クラスを代入してみる（SOLID：リスコフの置換原則）と・・・

	var bid *Bid
	bid = &Bid{price: 1}
	fmt.Printf("bid: %+v\n", bid) // これは、もちろんOK

	var bid2 *Bid
	// bid2 = &Depth{} // ←　これがあるとコンパイルエラー
	fmt.Printf("depth(bid2): %+v\n", bid2)
	// Compile Error
	// cannot use &Depth literal (type *Depth) as type *Bid in assignment
}

func Test8(t *testing.T) {
	// NG（goのembeded fieldは、"is a"ではなく、"has a"関係であるため）
	// ベースクラスに派生クラスを代入してみる（SOLID：リスコフの置換原則）と・・・

	// base
	var ask *Ask
	ask = &Ask{price: 2, amount: 3}
	fmt.Printf("ask: %+v\n", ask) // これは、もちろんOK
	fmt.Println(ask.getExecutedPrice())

	// derived
	var depth *Depth = &Depth{timestamp: 1, Bid: Bid{price: 4, amount: 5}}
	fmt.Printf("depth: %+v\n", depth) // これも、もちろんOK
	fmt.Println(depth.getExecutedPrice())

	// derived
	var depth2 *DepthAsk = &DepthAsk{timestamp: 1, Ask: Ask{price: 6, amount: 7}}
	fmt.Printf("depth2: %+v\n", depth2) // これも、もちろんOK
	fmt.Println(depth2.getExecutedPrice())

	// baseでderivedを受けられるか？
	// 型が異なるので、受けられない。
	// var baseAsk *Ask = &DepthAsk{timestamp: 1, Ask: Ask{price: 4, amount: 5}}
	// Compile Error
	// cannot use &DepthAsk literal (type *DepthAsk) as type *Ask in assignment
}

func Test9(t *testing.T) {
	// OK
	// interfaceを利用してみると・・・

	// base
	var ask *Ask
	ask = &Ask{price: 2, amount: 3}
	fmt.Printf("ask: %+v\n", ask) // これは、もちろんOK
	fmt.Println(ask.getExecutedPrice())

	// derived
	var depth *Depth = &Depth{timestamp: 1, Bid: Bid{price: 4, amount: 5}}
	fmt.Printf("depth: %+v\n", depth) // これも、もちろんOK
	fmt.Println(depth.getExecutedPrice())

	// derived
	var depth2 *DepthAsk = &DepthAsk{timestamp: 1, Ask: Ask{price: 6, amount: 7}}
	fmt.Printf("depth2: %+v\n", depth2) // これも、もちろんOK
	fmt.Println(depth2.getExecutedPrice())

	depths := []DepthInfo{
		ask,
		depth,
		depth2,
	}
	for _, depth := range depths {
		fmt.Println(getAppliedFeePrice(depth))
	}
	// Result
	// This method is a base structure.
	// -44
	// This method is a derived structure.
	// -30
	// This method is a derived structure.
	// -8
}
