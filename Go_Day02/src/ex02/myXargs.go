package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var out []string

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Unknown panic happend: ", err)
		}
	}()

	out = ParseStdout()
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	cmd.Args = append(cmd.Args, out...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("myError: %v", err)
		os.Exit(2)
	}

}

func ParseStdout() []string {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		slice := strings.Fields(scanner.Text())
		lines = append(lines, slice...)
	}

	return lines
}
