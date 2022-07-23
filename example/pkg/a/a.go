package a

import (
	"fmt"
	pkgb "go-example/pkg/b"
)

type A struct {
	pkgb.B
}

func (a A) Print() {
	fmt.Println("A Print")
}
