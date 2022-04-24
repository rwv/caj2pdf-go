package parser

type CAJParser struct {
	filePath string
}

func NewCAJParser(filePath string) CAJParser {
	return CAJParser{
		filePath: filePath,
	}
}
