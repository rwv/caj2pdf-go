package CAJParser

import (
	"bytes"
	"encoding/binary"
	"io"
)

func extractData(file io.ReadSeeker) (*bytes.Reader, error) {
	var err error

	var startOffset int64 = caj_PAGE_NUMBER_OFFSET + 4
	file.Seek(startOffset, io.SeekStart)

	// Seek to PDF start pointer
	var pdfStartPointerSlice []byte = make([]byte, 4)

	_, err = file.Read(pdfStartPointerSlice)
	if err != nil {
		return nil, err
	}

	pdfStartPointer := int64(int32(binary.LittleEndian.Uint32(pdfStartPointerSlice)))
	file.Seek(pdfStartPointer, io.SeekStart)

	var pdf_start []byte = make([]byte, 4)
	_, err = file.Read(pdf_start)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	pdfFooter := []byte("\r\n")

	// Concat
	pdfData := append(pdfHeader, pdfBody...)
	pdfData = append(pdfData, pdfFooter...)

	return bytes.NewReader(pdfData), nil
}
