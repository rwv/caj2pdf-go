package parser

type KDHParser struct {
	filePath string
}

func NewKDHParser(filePath string) KDHParser {
	return KDHParser{
		filePath: filePath,
	}
}
