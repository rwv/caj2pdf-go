package CAJParser

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"

	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	_TOC_NUMBER_OFFSET   int64 = 0x110
	_TOC_TITLE_OFFSET    int64 = 256
	_TOC_UNKNOWN1_OFFSET int64 = 24
	_TOC_PAGE_OFFSET     int64 = 12
	_TOC_UNKNOWN2_OFFSET int64 = 12
	_TOC_LEVEL_OFFSET    int64 = 16
	tocStructLength      int64 = _TOC_TITLE_OFFSET + _TOC_PAGE_OFFSET + _TOC_UNKNOWN1_OFFSET + _TOC_UNKNOWN2_OFFSET + _TOC_LEVEL_OFFSET
)

func getToc(file io.ReadSeeker) []toc {
	tocList := make([]toc, 0)

	if _TOC_NUMBER_OFFSET == 0 {
		return tocList
	}

	tocNum := getTocNum(file)
	for i := int64(0); i < int64(tocNum); i++ {
		file.Seek(_TOC_NUMBER_OFFSET+4+0x134*i, io.SeekStart)
		var tocSlice []byte = make([]byte, tocStructLength)
		_, err := file.Read(tocSlice)
		if err != nil {
			continue
		}
		tocList = append(tocList, parseToc(tocSlice))
	}

	return tocList
}

func getTocNum(file io.ReadSeeker) int32 {
	file.Seek(_TOC_NUMBER_OFFSET, io.SeekStart)

	var tocSlice []byte = make([]byte, 4)
	_, err := file.Read(tocSlice)
	if err != nil {
		return 0
	}

	tocNum := int32(binary.LittleEndian.Uint32(tocSlice))
	return tocNum
}

func parseToc(byteSlice []byte) toc {
	titlePart := byteSlice[0:256]
	titleEnd := bytes.IndexByte(titlePart, '\x00')
	titleBytes := titlePart[:titleEnd]
	decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(titleBytes)
	titleStr := string(decodeBytes)

	pagePart := byteSlice[_TOC_TITLE_OFFSET+_TOC_UNKNOWN1_OFFSET : _TOC_TITLE_OFFSET+_TOC_UNKNOWN1_OFFSET+_TOC_PAGE_OFFSET]
	pageEnd := bytes.IndexByte(pagePart, '\x00')
	pageStr := string(pagePart[:pageEnd])
	pageInt, _ := strconv.Atoi(pageStr)

	levelPart := byteSlice[_TOC_TITLE_OFFSET+_TOC_UNKNOWN1_OFFSET+_TOC_PAGE_OFFSET+_TOC_UNKNOWN2_OFFSET : _TOC_TITLE_OFFSET+_TOC_UNKNOWN1_OFFSET+_TOC_PAGE_OFFSET+_TOC_UNKNOWN2_OFFSET+_TOC_LEVEL_OFFSET]
	levelInt := int32(binary.LittleEndian.Uint32(levelPart))

	tocItem := toc{
		Title: titleStr,
		Page:  int32(pageInt),
		Level: levelInt,
	}
	return tocItem
}

type toc struct {
	Title string
	Page  int32
	Level int32
}
