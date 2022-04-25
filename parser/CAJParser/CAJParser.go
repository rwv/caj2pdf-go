package CAJParser

import (
	"bytes"
	"os"
)

const caj_TOC_NUMBER_OFFSET = 0x110

type CAJParser struct {
	filePath string
	pageNum  int32
}

func New(filePath string) CAJParser {
	parser := CAJParser{
		filePath: filePath,
	}

	file, err := os.Open(parser.filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	parser.pageNum = getPageNum(file)
	return parser
}

func (parser CAJParser) Convert(target string) error {
	file, err := os.Open(parser.filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	pdfData, err := extractData(file)

	pdfData, err = handlePages(bytes.NewReader(pdfData), &parser)

	// write pdfData to File
	file, err = os.CreateTemp("", "caj2pdf")
	if err != nil {
		panic(err)
	}

	_, err = file.Write(pdfData)
	if err != nil {
		panic(err)
	}
	file.Close()

	repairXref(file.Name(), target)

	return nil
}
