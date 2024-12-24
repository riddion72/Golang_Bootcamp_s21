package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type Flags struct {
	f, d, sl *bool
	ext      *string
}

func main() {
	var option Flags
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Unknown panic happend: ", err)
		}
	}()
	option.f = flag.Bool("f", false, "print files")
	option.d = flag.Bool("d", false, "print  directories")
	option.sl = flag.Bool("sl", false, "print symlinks")
	option.ext = flag.String("ext", "", "files with a certain extension")
	flag.Parse()

	if !*option.f && *option.ext != "" {
		fmt.Println("Requires the flag -f to use the flag -ext")
		os.Exit(2)
	}
	if !*option.f && !*option.d && !*option.sl {
		*option.f, *option.d, *option.sl = true, true, true
	}
	filepath.Walk(flag.Arg(0), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if *option.sl && info.Mode()&fs.ModeSymlink != 0 {
			fmt.Print(path + "->")
			sLink, err := filepath.EvalSymlinks(path)
			if err != nil {
				fmt.Println("[broken]")
			} else {
				fmt.Println(sLink)
			}
		} else if *option.d && info.IsDir() {
			fmt.Println(path)
		} else if *option.f && info.Mode()>>9 == 0 {
			if *option.ext != "" {
				if filepath.Ext(path) == "."+*option.ext {
					fmt.Println(path)
				}
			} else {
				fmt.Println(path)
			}
		}
		return nil
	})
}
