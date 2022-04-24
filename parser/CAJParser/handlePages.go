package CAJParser

import (
	"bytes"
	"fmt"
	"io"
)

// get pages_obj_no list containing distinct elements
// & find missing pages object(s) -- top pages object(s) in pages_obj_no
func handlePages(pdfData []byte) ([]byte, error) {
	reader := bytes.NewReader(pdfData)

	obj_no := dealDisordered(reader)
	inds := addCatalog(reader)

	pages_obj_no := []int64{}
	top_pages_obj_no := []int64{}

	for _, ind := range inds {
		if !contains(pages_obj_no, ind) && !contains(top_pages_obj_no, ind) {
			if find(reader, []byte(fmt.Sprintf("\r%d 0 obj", ind)), 0) != -1 {
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

	return pdfData, nil
}
