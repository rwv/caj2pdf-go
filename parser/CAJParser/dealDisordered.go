package CAJParser

import (
	"io"
	"strconv"
)

func dealDisordered(reader io.ReadSeeker) []int64 {
	endobjSlice := findAllOccurances(reader, []byte("endobj"))
	obj_no := []int64{}

	for _, endobjAddress := range endobjSlice {
		startObj := findReverse(reader, []byte(" 0 obj"), endobjAddress)
		startObj1 := findReverse(reader, []byte("\r"), startObj)
		startObj2 := findReverse(reader, []byte("\n"), startObj)
		// startObj = Max(startObj1, startObj2)
		if startObj1 > startObj2 {
			startObj = startObj1
		} else {
			startObj = startObj2
		}
		length := find(reader, []byte(" "), startObj) - startObj

		reader.Seek(startObj, io.SeekStart)
		var numberStringBytes []byte = make([]byte, length)
		reader.Read(numberStringBytes)
		number_, _ := strconv.Atoi(string(numberStringBytes))
		number := int64(number_)
		if !contains(obj_no, number) {
			obj_no = append(obj_no, number)
			// obj_len = addr - startobj + 6 // from original caj2pdf
			reader.Seek(startObj, io.SeekStart)
			// [obj] = struct.unpack(str(obj_len) + "s", pdf.read(obj_len)) // from original caj2pdf
		}
	}

	return obj_no
}
