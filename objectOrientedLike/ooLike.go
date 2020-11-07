package objectOrientedLike

import "fmt"

// Bid is base
type Bid struct {
	price  int
	amount int
}

func (b *Bid) getExecutedPrice() (executedPrice int) {
	fmt.Println("This method is a base structure.")
	return b.price * b.amount
}

// Ask is base
type Ask struct {
	price  int
	amount int
}

func (a *Ask) getExecutedPrice() (executedPrice int) {
	fmt.Println("This method is a base structure.")
	return a.price * a.amount
}

func (a *Ask) getAppliedFeePrice() (appliedFeePrice int) {
	return a.getExecutedPrice() - 10 // fee is 10
}

// Depth is derived
type Depth struct {
	timestamp int
	Bid       // embedded field
}

func newDepth() *Depth {
	return &Depth{}
}

func (d *Depth) getExecutedPrice() (executedPrice int) {
	fmt.Println("This method is a derived structure.")
	return d.price * d.amount
}

// DepthAsk is derived
type DepthAsk struct {
	timestamp int
	Ask       // embedded field
}

func (d *DepthAsk) getExecutedPrice() (executedPrice int) {
	fmt.Println("This method is a derived structure.")
	return d.price * d.amount
}

// DepthInfo is interface
type DepthInfo interface {
	getExecutedPrice() (executedPrice int)
}

func getAppliedFeePrice(depth DepthInfo) (appliedFeePrice int) {
	return depth.getExecutedPrice() - 50 // fee is 50
}
