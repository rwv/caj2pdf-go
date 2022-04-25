package CAJParser

import (
	"bytes"
	"io"
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

	toc := getToc(file)

	writer := new(bytes.Buffer)
	err = extractData(file, writer)

	extractedReader := bytes.NewReader(writer.Bytes())
	writer = new(bytes.Buffer)
	handlePages(extractedReader, writer, &parser)

	reader := bytes.NewReader(writer.Bytes())
	writer = new(bytes.Buffer)
	repairXref(reader, writer)

	reader = bytes.NewReader(writer.Bytes())
	writer = new(bytes.Buffer)
	addOutlines(reader, writer, toc)

	reader = bytes.NewReader(writer.Bytes())
	file1, err := os.Create(target)
	if err != nil {
		panic(err)
	}
	defer file1.Close()
	io.Copy(file1, reader)

	return nil
}
