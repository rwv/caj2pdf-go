package parser

type CAJParser struct {
	filePath string
}

func NewCajParser(filePath string) CAJParser {
	return CAJParser{
		filePath: filePath,
	}
}
