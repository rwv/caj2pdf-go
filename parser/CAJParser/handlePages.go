package CAJParser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// get pages_obj_no list containing distinct elements
// & find missing pages object(s) -- top pages object(s) in pages_obj_no
func handlePages(reader io.ReadSeeker, parser *CAJParser) ([]byte, error) {
	pdfData, _ := io.ReadAll(reader)

	obj_no := dealDisordered(reader)
	inds := addCatalog(reader)

	pages_obj_no := []int64{}
	top_pages_obj_no := []int64{}

	for _, ind := range inds {
		if !contains(pages_obj_no, ind) && !contains(top_pages_obj_no, ind) {
			if find(reader, []byte(fmt.Sprintf("\r%d 0 obj", ind)), 0) == -1 {
				top_pages_obj_no = append(top_pages_obj_no, ind)
			} else {
				pages_obj_no = append(pages_obj_no, ind)
			}
		}
	}

	single_pages_obj_missed := len(top_pages_obj_no) == 1
	multi_pages_obj_missed := len(top_pages_obj_no) > 1

	// generate catalog object
	catalog_obj_no := findUnusedNo(obj_no, top_pages_obj_no)
	obj_no = append(obj_no, catalog_obj_no)

	var root_pages_obj_no int64

	if multi_pages_obj_missed {
		root_pages_obj_no = findUnusedNo(obj_no, top_pages_obj_no)
	} else if single_pages_obj_missed {
		root_pages_obj_no = top_pages_obj_no[0]
		top_pages_obj_no = pages_obj_no
	} else { // root pages object exists, then find the root pages object
		found := false
		for _, pon := range pages_obj_no {
			tmp_addr := find(reader, []byte(fmt.Sprintf("\r%d 0 obj", pon)), 0)
			for {
				reader.Seek(tmp_addr, io.SeekStart)
				_str := make([]byte, 6)
				reader.Read(_str)
				str := string(_str)
				if str == "Parent" {
					break
				} else if str == "endobj" {
					root_pages_obj_no = pon
					found = true
					break
				}
				tmp_addr = tmp_addr + 1
			}
			if found {
				break
			}
		}
	}

	catalog := []byte(fmt.Sprintf("%d 0 obj\r<</Type /Catalog\r/Pages %d 0 R\r>>\rendobj\r", catalog_obj_no, root_pages_obj_no))

	pdfData = append(pdfData, catalog...)

	pdfData = addPagesObjAndEOFMark(pdfData, single_pages_obj_missed, multi_pages_obj_missed, top_pages_obj_no, root_pages_obj_no, parser.pageNum)

	return pdfData, nil
}

func addPagesObjAndEOFMark(
	pdfData []byte,
	single_pages_obj_missed bool,
	multi_pages_obj_missed bool,
	top_pages_obj_no []int64,
	root_pages_obj_no int64,
	pageNum int32,
) []byte {
	// Add Pages obj and EOF mark
	// if root pages object exist, pass
	// deal with single missing pages object
	if single_pages_obj_missed || multi_pages_obj_missed {
		inds_str := make([]string, len(top_pages_obj_no))
		for idx, i := range top_pages_obj_no {
			inds_str[idx] = fmt.Sprintf("%d 0 R", i)
		}

		kids_str := "[" + strings.Join(inds_str, " ") + "]"

		pages_str := fmt.Sprintf("%d 0 obj\r<<\r/Type /Pages\r/Kids %s\r/Count %d\r>>\rendobj\r", root_pages_obj_no, kids_str, pageNum)

		pdfData = append(pdfData, []byte(pages_str)...)
	}

	// deal with multiple missing pages objects
	if multi_pages_obj_missed {
		kids_dict := make(map[int64][]int64)
		for _, i := range top_pages_obj_no {
			kids_dict[i] = []int64{}
		}

		count_dict := make(map[int64]int64)
		for _, i := range top_pages_obj_no {
			count_dict[i] = 0
		}

		reader := bytes.NewReader(pdfData)

		for _, tpon := range top_pages_obj_no {
			kids_addr := findAllOccurances(reader, []byte(fmt.Sprintf("/Parent %d 0 R", tpon)))
			for _, kid := range kids_addr {
				ind := findReverse(reader, []byte("obj"), kid) - 4
				addr := findReverse(reader, []byte("\r"), ind)
				length := find(reader, []byte(" "), addr) - addr
				reader.Seek(addr, io.SeekStart)
				indStrSlice := make([]byte, length)
				reader.Read(indStrSlice)
				indStr := string(indStrSlice)
				ind_, _ := strconv.Atoi(indStr)
				ind = int64(ind_)
				kids_dict[tpon] = append(kids_dict[tpon], ind)

				type_addr := find(reader, []byte("/Type"), addr) + 5
				tmp_addr := find(reader, []byte("/"), type_addr) + 1

				reader.Seek(tmp_addr, io.SeekStart)

				typeStrSlice := make([]byte, 5)
				reader.Read(typeStrSlice)
				typeStr := string(typeStrSlice)
				if typeStr == "Pages" {
					cnt_addr := find(reader, []byte("/Count "), addr) + 7
					reader.Seek(cnt_addr, io.SeekStart)

					_strSlice := make([]byte, 1)
					reader.Read(_strSlice)
					var cnt_len int64 = 0

					for !contains([]byte{' ', '\r', '/'}, _strSlice[0]) {
						cnt_len += 1
						reader.Seek(cnt_addr+cnt_len, io.SeekStart)
						reader.Read(_strSlice)
					}

					reader.Seek(cnt_addr, io.SeekStart)
					cntSlice := make([]byte, cnt_len)
					reader.Read(cntSlice)
					ind_, _ := strconv.Atoi(string(cntSlice))
					ind = int64(ind_)
					count_dict[tpon] += ind
				} else { // _type == b"Page"
					count_dict[tpon] += 1
				}
			}
			kids_no_str := make([]string, len(kids_dict[tpon]))

			for idx, i := range kids_dict[tpon] {
				kids_no_str[idx] = fmt.Sprintf("%d 0 R", i)
			}

			kids_str := "[" + strings.Join(kids_no_str, " ") + "]"
			pages_str := fmt.Sprintf("%d 0 obj\r<<\r/Type /Pages\r/Kids %s\r/Count %d\r>>\rendobj\r", tpon, kids_str, count_dict[tpon])
			pdfData = append(pdfData, []byte(pages_str)...)
		}
	}

	// add EOF mark
	pdfData = append(pdfData, []byte("\n%%EOF\r")...)

	return pdfData
}
