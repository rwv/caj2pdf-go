package parser

type HNParser struct {
	filePath string
}

func NewHNParser(filePath string) HNParser {
	return HNParser{
		filePath: filePath,
	}
}
