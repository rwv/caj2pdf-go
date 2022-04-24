package parser

type C8Parser struct {
	filePath string
}

func NewC8Parser(filePath string) C8Parser {
	return C8Parser{
		filePath: filePath,
	}
}
