package main

import (
	"fmt"

	"github.com/rwv/caj2pdf-go/parser"
)

func main() {
	parser, _ := parser.New("temp/1.caj")
	fmt.Println(parser)
}
