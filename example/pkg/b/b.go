package b

import (
	"go-example/pkg/bc"
)

type B struct {
	bc.BCBridge
}

func (b B) Print() {
	b.Printer.Print()
}
