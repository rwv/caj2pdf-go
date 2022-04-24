package CAJParser

import (
	"bytes"
	"io"
)

func findAllOccurances(file io.ReadSeeker, pattern []byte) []int64 {
	var results []int64
	var last_address int64 = int64(-len(pattern))

	for {
		addr := find(file, pattern, last_address+int64(len(pattern)))
		if addr != -1 {
			results = append(results, addr)
			last_address = addr
		} else {
			return results
		}
	}
}

func find(file io.ReadSeeker, pattern []byte, start int64) int64 {
	patternLen := len(pattern)
	fsize, _ := file.Seek(0, io.SeekEnd)
	file.Seek(0, io.SeekStart)

	const bsize int64 = 4096

	file.Seek(start, io.SeekStart)

	if start > 0 {
		file.Seek(start, io.SeekStart)
	}
	overlap := int64(len(pattern) - 1)

	for {
		currentOffset, _ := file.Seek(0, io.SeekCurrent)
		if overlap <= currentOffset && currentOffset < fsize {
			currentOffset = currentOffset - overlap
			file.Seek(currentOffset, io.SeekStart)
		}

		buf := make([]byte, bsize)
		n, _ := file.Read(buf)
		if n == 0 {
			return -1
		} else {
			for i := 0; i < n-patternLen; i++ {
				if bytes.Compare(buf[i:i+patternLen], pattern) == 0 {
					return currentOffset + int64(i)
				}
			}
		}
	}
}
