package CAJParser

import (
	"bytes"
	"io"
)

func contains[K comparable](s []K, e K) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

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

func findReverse(file io.ReadSeeker, pattern []byte, end int64) int64 {
	patternLen := len(pattern)
	fsize, _ := file.Seek(0, io.SeekEnd)
	const bsize int64 = 4096
	if int64(patternLen) > end {
		panic("Too large string size for search.")
	}

	file.Seek(fsize-bsize, io.SeekStart)

	size := bsize
	if bsize <= end && end < fsize {
		file.Seek(end-bsize, io.SeekStart)
	} else if 0 < end && end < bsize {
		size = end
		file.Seek(0, io.SeekStart)
	}

	overlap := int64(patternLen - 1)

	patternReverse := make([]byte, len(pattern))

	for i := 0; i < len(pattern); i++ {
		patternReverse[i] = pattern[len(pattern)-1-i]
	}

	for {
		buf := make([]byte, size)
		n, _ := file.Read(buf)
		currentOffset, _ := file.Seek(0, io.SeekCurrent)
		if n > 0 {
			bufReverse := make([]byte, n)
			for i := 0; i < n; i++ {
				bufReverse[i] = buf[n-1-i]
			}

			// Find Pos
			var pos int64 = -1

			for i := 0; i < n-patternLen; i++ {
				if bytes.Compare(bufReverse[i:i+patternLen], patternReverse) == 0 {
					pos = int64(i)
					break
				}
			}

			if pos >= 0 {
				return currentOffset - pos
			}
		}

		if (2*bsize - overlap) < currentOffset {
			file.Seek(currentOffset-(2*bsize-overlap), io.SeekStart)
			size = bsize
		} else if (bsize - overlap) < currentOffset {
			size = currentOffset - (bsize - overlap)
			file.Seek(0, io.SeekStart)
		} else {
			return -1
		}
	}

}

func findUnusedNo(numberList []int64, numberList2 []int64) int64 {
	var unuse_no int64 = -1

	const maxNo = 99999
	for i := int64(0); i < maxNo; i++ {
		temp := maxNo - i
		if !contains(numberList, temp) && !contains(numberList2, temp) {
			unuse_no = temp
			break
		}
	}
	if unuse_no == -1 {
		panic("Error on PDF objects numbering.")
	}
	return unuse_no
}
