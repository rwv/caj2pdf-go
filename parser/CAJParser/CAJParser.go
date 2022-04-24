package CAJParser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const caj_PAGE_NUMBER_OFFSET = 0x10
const caj_TOC_NUMBER_OFFSET = 0x110

type CAJParser struct {
	filePath string
}

func New(filePath string) CAJParser {
	return CAJParser{
		filePath: filePath,
	}
}

func (parser CAJParser) Convert(target string) error {
	file, err := os.Open(parser.filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var startOffset int64 = caj_PAGE_NUMBER_OFFSET + 4
	file.Seek(startOffset, io.SeekStart)

	// Seek to PDF start pointer
	var pdfStartPointerSlice []byte = make([]byte, 4)

	_, err = file.Read(pdfStartPointerSlice)
	if err != nil {
		return err
	}

	pdfStartPointer := int64(int32(binary.LittleEndian.Uint32(pdfStartPointerSlice)))
	file.Seek(pdfStartPointer, io.SeekStart)

	var pdf_start []byte = make([]byte, 4)
	_, err = file.Read(pdf_start)
	if err != nil {
		return err
	}

	pdf_start_value := int64(binary.LittleEndian.Uint32(pdf_start))

	endPattern := []byte("endobj")
	occurances := findAllOccurances(file, endPattern)
	pdf_end := occurances[len(occurances)-1] + 6

	pdfLength := pdf_end - pdf_start_value

	file.Seek(pdf_start_value, io.SeekStart)

	pdfHeader := []byte("%PDF-1.3\r\n")
	pdfBody := make([]byte, pdfLength)

	_, err = file.Read(pdfBody)
	if err != nil {
		return err
	}
	pdfFooter := []byte("\r\n")

	// Concat
	pdfData := append(pdfHeader, pdfBody...)
	pdfData = append(pdfData, pdfFooter...)

	// Write to file
	fmt.Println("Writing to file...")
	out, err := os.Create("pdf.tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write(pdfData)

	return nil
}

func findAllOccurances(file *os.File, pattern []byte) []int64 {
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

func find(file *os.File, pattern []byte, start int64) int64 {
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
