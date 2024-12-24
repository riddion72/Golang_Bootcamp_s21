package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

type Flags struct {
	l, m, w *bool
}

func flagChecker(l, m, w bool) (err error) {
	if (l && m) || (l && w) || (m && w) {
		return errors.New("You can use only one flag")
	}
	return nil
}

func worldCounter(filePath string, flag int, wg *sync.WaitGroup) {
	var cnt int
	defer wg.Done()
	contFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Can't open new file")
		os.Exit(4)
	}
	line := bufio.NewScanner(contFile)
	for line.Scan() {
		curLine := line.Text()
		if flag <= 2 {
			cnt++
			if flag == 2 {
				cnt += utf8.RuneCountInString(curLine)
			}
		} else if flag == 3 {
			for _, word := range strings.Split(curLine, " ") {
				if word != "" {
					cnt++
				}
			}
		}
	}
	fmt.Printf("%8.v %s\n", cnt, filePath)
	return
}

func main() {
	wg := new(sync.WaitGroup)
	var option Flags
	var fl int
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Unknown panic happend: ", err)
		}
	}()
	option.l = flag.Bool("l", false, "counting lines")
	option.m = flag.Bool("m", false, "counting characters")
	option.w = flag.Bool("w", false, "counting words")
	flag.Parse()

	if err := flagChecker(*option.l, *option.m, *option.w); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(2)
	}

	for _, filePath := range flag.Args() {
		wg.Add(1)
		if *option.l {
			fl = 1
		} else if *option.m {
			fl = 2
		} else if *option.w {
			fl = 3
		}
		go worldCounter(filePath, fl, wg)
	}
	wg.Wait()
}
