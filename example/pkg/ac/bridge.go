package ac

import (
	pkga "go-example/pkg/a"
)

type ACBridge struct {
	pkga.A
}

func (ac *ACBridge) Print() {
	// Fowarding
	ac.A.Print()
}
