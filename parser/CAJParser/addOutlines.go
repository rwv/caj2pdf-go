package CAJParser

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
)

func addOutlines(source io.ReadSeeker, output io.Writer, toc []toc) error {
	// write pdfData to File
	tempFile, err := os.CreateTemp("", "caj2pdf")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tempFile.Name())

	// Copy source to temp file
	_, err = io.Copy(tempFile, source)

	tempFile.Close()

	tempFolder, err := os.MkdirTemp("", "caj2pdfDir")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempFolder)

	outputFilename := path.Join(tempFolder, "temp.pdf")

	outlineFile, err := os.CreateTemp("", "caj2pdf")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tempFile.Name())

	_generateOutlineFile(toc, outlineFile)

	_addOutlines(tempFile.Name(), outputFilename, outlineFile.Name())

	// Copy repaired file to target
	outputFileBuffer, err := os.Open(outputFilename)
	if err != nil {
		panic(err)
	}
	defer outputFileBuffer.Close()

	_, err = io.Copy(output, outputFileBuffer)

	return nil
}

func _addOutlines(source string, target string, outline string) {
	cmd := exec.Command("pdfoutline", source, outline, target)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func _generateOutlineFile(toc []toc, file io.Writer) {
	outline := ""
	for _, item := range toc {
		outline += fmt.Sprintf("%d %d %s\n", item.Level, item.Page, item.Title)
	}
	file.Write([]byte(outline))
}
