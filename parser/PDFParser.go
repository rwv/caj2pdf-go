package parser

import (
	"io"
	"os"
)

type PDFParser struct {
	filePath string
}

func NewPDFParser(filePath string) PDFParser {
	return PDFParser{
		filePath: filePath,
	}
}

func (parser PDFParser) Convert(target string) error {
	// Just change the file extension to .pdf is enough
	source, err := os.Open(parser.filePath)
	if err != nil {
		return err
	}

	destination, err := os.Create(target)
	if err != nil {
		return err
	}

	_, err = io.Copy(destination, source)
	return err
}
