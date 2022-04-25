package CAJParser

import (
	"encoding/binary"
	"io"
)

func getPageNum(file io.ReadSeeker) int32 {
	file.Seek(_PAGE_NUMBER_OFFSET, io.SeekStart)

	var pageSlice []byte = make([]byte, 4)
	_, err := file.Read(pageSlice)
	if err != nil {
		return -1
	}

	pageNum := int32(binary.LittleEndian.Uint32(pageSlice))
	return pageNum
}
