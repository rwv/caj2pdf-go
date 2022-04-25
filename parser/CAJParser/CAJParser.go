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

	writer := bytes.NewBuffer([]byte{})

	err = extractData(file, writer)

	extractedReader := bytes.NewReader(writer.Bytes())

	handledWriter := bytes.NewBuffer([]byte{})

	_ = handlePages(extractedReader, handledWriter, &parser)

	// pdfData := handledWriter.Bytes()

	// write pdfData to File
	file, err = os.CreateTemp("", "caj2pdf")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(file, handledWriter)

	file.Close()

	repairXref(file.Name(), target)

	return nil
}
