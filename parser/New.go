package parser

func New(filePath string) Parser {
	return NewCajParser(filePath)
}
