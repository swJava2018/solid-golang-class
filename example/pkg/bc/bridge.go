package bc

import "go-example/pkg/printer"

type BCBridge struct {
	printer.Printer
}

func (bc BCBridge) Print() {
	bc.Printer.Print()
}
