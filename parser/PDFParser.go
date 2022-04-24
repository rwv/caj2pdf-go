package parser

type PDFParser struct {
	filePath string
}

func NewPDFParser(filePath string) PDFParser {
	return PDFParser{
		filePath: filePath,
	}
}
