package parser

type ParserType int

const (
	Unknown ParserType = -1
	C8      ParserType = iota
	HN
	CAJ
	KDH
	PDF
	TEB
)
