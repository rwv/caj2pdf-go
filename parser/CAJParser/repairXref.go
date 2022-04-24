package CAJParser

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func repairXref(source string, target string) {
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
