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

func repairXref(source io.Reader, target io.Writer) {
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

	repairedFilename := path.Join(tempFolder, "temp.pdf")

	_repairXref(tempFile.Name(), repairedFilename)

	// Copy repaired file to target
	repairedFile, err := os.Open(repairedFilename)
	if err != nil {
		panic(err)
	}
	defer repairedFile.Close()

	_, err = io.Copy(target, repairedFile)
}

func _repairXref(source string, target string) {
	cmd := exec.Command("mutool", "clean", source, target)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
