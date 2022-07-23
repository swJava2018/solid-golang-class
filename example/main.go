package main

import (
	"fmt"
	aPkg "go-example/pkg/a"
	bPkg "go-example/pkg/b"
	bcPkg "go-example/pkg/bc"
	cPkg "go-example/pkg/c"
)

func main() {
	fmt.Println("hello world")
	c := &cPkg.C{}
	//A Print
	c.Print()

	a := &aPkg.A{}
	bc := &bcPkg.BCBridge{Printer: a}
	b := &bPkg.B{BCBridge: *bc}
	//A Print
	b.Print()
}
