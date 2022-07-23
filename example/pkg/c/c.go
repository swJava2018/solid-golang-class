package c

import "go-example/pkg/ac"

type C struct {
	ac.ACBridge
}

func (c C) Print() {
	c.ACBridge.Print()
}
