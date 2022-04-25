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

	writer := new(bytes.Buffer)
	err = extractData(file, writer)

	extractedReader := bytes.NewReader(writer.Bytes())
	writer = new(bytes.Buffer)
	handlePages(extractedReader, writer, &parser)

	reader := bytes.NewReader(writer.Bytes())
	writer = new(bytes.Buffer)
	repairXref(reader, writer)

	return nil
}
