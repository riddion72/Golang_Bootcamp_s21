package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func fileFormat(path string) {
	if !strings.HasSuffix(path, ".txt") {
		fmt.Println("Wrong file format (necessary txt format)")
		os.Exit(4)
	}
}

func serchAdded(oldPath string, newPath string) {
	oldFile, oldErr := os.ReadFile(oldPath)
	if oldErr != nil {
		fmt.Println("No such old file")
		os.Exit(3)
	}
	content := string(oldFile)
	oldStrings := strings.Split(content, "\n")
	newFile, newErr := os.Open(newPath)
	if newErr != nil {
		fmt.Println("Can't open new file")
		os.Exit(4)
	}
	newLine := bufio.NewScanner(newFile)
	for newLine.Scan() {
		var count int
		var line string
		for _, line = range oldStrings {
			if line == newLine.Text() {
				count++
			}
		}
		if count == 0 {
			fmt.Println("ADDED ", newLine.Text())
		}
	}
	newFile.Close()
}

func serchRemuved(oldPath string, newPath string) {
	newFile, newErr := os.ReadFile(newPath)
	if newErr != nil {
		fmt.Println("No such new file")
		os.Exit(5)
	}
	content := string(newFile)
	newStrings := strings.Split(content, "\n")
	oldFile, oldErr := os.Open(oldPath)
	if oldErr != nil {
		fmt.Println("Can't open old file")
		os.Exit(6)
	}
	oldLine := bufio.NewScanner(oldFile)
	for oldLine.Scan() {
		var count int
		var line string
		for _, line = range newStrings {
			if line == oldLine.Text() {
				count++
			}
		}
		if count == 0 {
			fmt.Println("REMOVED ", oldLine.Text())
		}
	}
	oldFile.Close()
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Unknown panic happend: ", err)
		}
	}()
	if len(os.Args) != 5 {
		fmt.Println("Use patern --old <old_file_name>.txt --new <new_file_name>.txt")
		os.Exit(7)
	}
	inputArgs := os.Args[1:]
	if inputArgs[0] == "--old" && inputArgs[2] == "--new" {
		fileFormat(inputArgs[1])
		fileFormat(inputArgs[3])

	} else {
		fmt.Println("Use --old --new flags.")
		os.Exit(8)
	}
	serchAdded(inputArgs[1], inputArgs[3])
	serchRemuved(inputArgs[1], inputArgs[3])
}
