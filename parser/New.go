package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func New(filePath string) (Parser, error) {
	fileType, _ := detectType(filePath)
	switch fileType {
	case C8:
		return NewC8Parser(filePath), nil
	case HN:
		return NewHNParser(filePath), nil
	case KDH:
		return NewKDHParser(filePath), nil
	case PDF:
		return NewPDFParser(filePath), nil
	case CAJ:
		return NewCAJParser(filePath), nil
	}
	return nil, fmt.Errorf("Unknown file type")
}

func detectType(filePath string) (ParserType, error) {
	fi, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	r := bufio.NewReader(fi)

	var head []byte
	head, err = r.Peek(4)
	if err != nil {
		return Unknown, err
	}

	// Check C8
	if head[1] == '\xc8' {
		return C8, nil
	}

	// Check HN
	if string(head[:2]) == "HN" {
		var head2 []byte
		head2, err = r.Peek(2)
		if err != nil {
			return Unknown, err
		}
		if head2[0] == '\xc8' && head2[1] == '\x00' {
			return HN, nil
		}
	}

	format := strings.Replace(string(head), "\x00", "", -1)

	fmt.Println(format)

	switch format {
	case "CAJ":
		return CAJ, nil
	case "HN":
		return HN, nil
	case "KDH":
		return KDH, nil
	case "%PDF":
		return PDF, nil
	case "TEB":
		return TEB, nil
	}

	return Unknown, fmt.Errorf("Unknown file type")
}
