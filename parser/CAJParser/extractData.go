package CAJParser

import (
	"encoding/binary"
	"io"
)

const _PAGE_NUMBER_OFFSET int64 = 0x10

func extractData(file io.ReadSeeker, output io.Writer) error {
	var err error

	startOffset := _PAGE_NUMBER_OFFSET + 4
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

	output.Write(pdfHeader)
	output.Write(pdfBody)
	output.Write(pdfFooter)

	return nil
}
