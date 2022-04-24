package CAJParser

import (
	"io"
	"strconv"
)

func addCatalog(reader io.ReadSeeker) []int64 {
	indsParent := findAllOccurances(reader, []byte("/Parent"))
	indsAddr := make([]int64, len(indsParent))
	for i, addr := range indsParent {
		indsAddr[i] = addr + 8
	}

	inds := make([]int64, len(indsAddr))
	for i, addr := range indsAddr {
		length := find(reader, []byte(" "), addr) - addr
		reader.Seek(addr, io.SeekStart)
		var numberStringBytes []byte = make([]byte, length)
		reader.Read(numberStringBytes)
		number_, _ := strconv.Atoi(string(numberStringBytes))
		number := int64(number_)
		inds[i] = number
	}
	return inds
}
