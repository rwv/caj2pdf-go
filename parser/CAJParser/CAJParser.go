package CAJParser

import (
	"os"
)

const caj_PAGE_NUMBER_OFFSET = 0x10
const caj_TOC_NUMBER_OFFSET = 0x110

type CAJParser struct {
	filePath string
}

func New(filePath string) CAJParser {
	return CAJParser{
		filePath: filePath,
	}
}

func (parser CAJParser) Convert(target string) error {
	file, err := os.Open(parser.filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	extractedData, err := extractData(file)

	handlePages(extractedData)

	return nil
}
