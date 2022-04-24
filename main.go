package main

import (
	"fmt"

	"github.com/rwv/caj2pdf-go/parser/CAJParser"
)

func main() {
	parser := CAJParser.New("temp/5.caj")
	parser.Convert("/tmp/5.pdf")
	fmt.Println(parser)
}
