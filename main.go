package main

import (
	"fmt"

	"github.com/rwv/caj2pdf-go/parser"
)

func main() {
	parser := parser.NewCAJParser("temp/5.caj")
	parser.Convert("/tmp/5.pdf")
	fmt.Println(parser)
}
