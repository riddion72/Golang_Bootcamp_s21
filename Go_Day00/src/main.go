// package main

// import (
// 	"fmt"
// 	"bufio"
// 	"os"
// 	"sort"
// 	"strconv"
// 	"math"
// )

// func mean(data []int) float64 {
// 	if len(data) == 0 {
// 		return 0
// 	}
// 	var sum float64
// 	for _, d := range data {
// 		sum += float64(d)
// 	}
// 	return sum / float64(len(data))
// }

// func median(data []int) (median float64) {
//     dataLen := len(data)
//     if dataLen == 0 {
//         return 0
//     } else if dataLen%2 == 0 {
//         median = float64(data[dataLen / 2 - 1] + data[dataLen / 2]) / 2
//     } else {
//         median = float64(data[dataLen / 2])
//     }
//     return
// }

// func getMode(data []int) (mode int) {
// 	countMap := make(map[int]int)
// 	for _, value := range data {
// 		countMap[value] += 1
// 	}
// 	max := 0
// 	for _, key := range data {
// 		freq := countMap[key]
// 		if freq > max {
// 			mode = key
// 			max = freq
// 		}
// 	}
// 	return
// }

// func standardDeviation(data []int) (sd float64) {
// 	localMean := mean(data)
// 	for  _, value := range data {
// 		sd += float64(math.Pow(float64(value) - localMean, 2))
// 	}
// 	sd = float64(math.Sqrt(sd / float64(len(data))))
// 	return
// }

// func main() {
// 	var inputArr []int
// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println("Unknown panic happend: ", err)
// 		}
// 	}()
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for scanner.Scan() {
// 		input := scanner.Text()
// 		barInt, err := strconv.Atoi(input)
// 		if err != nil {
// 			fmt.Printf("Erorr: %v\n", err)
// 			os.Exit(3)
// 	    } else if barInt > 100000 || barInt < -100000 {
// 			fmt.Printf("Number %v out of range (from -100000 to 100000)\n", barInt)
// 			os.Exit(4)
// 		}
// 		inputArr = append(inputArr, barInt)
// 	}
// 	// fmt.Print(inputArr)
// 	inputArgs := os.Args[1:]
// 	sort.Ints(inputArr)
// 	if len(inputArgs) > 0 {
// 		for _, arg := range inputArgs {
// 			switch {
// 			case arg == "Mean":
// 				fmt.Printf("Mean: %.3g\n", mean(inputArr))
// 			case arg == "Median":
// 				fmt.Printf("Median: %.3g\n", median(inputArr))
// 			case arg == "Mode":
// 				fmt.Printf("Mode: %d\n", getMode(inputArr))
// 			case arg == "SD":
// 				fmt.Printf("Median: %.3g\n", standardDeviation(inputArr))
// 			default:
// 				fmt.Printf("I can't calculate %v\n", arg)
// 				os.Exit(5)
// 			}
// 		}
// 	} else {
// 		fmt.Printf("Mean: %.3g\n", mean(inputArr))
// 		fmt.Printf("Median: %.3g\n", median(inputArr))
// 		fmt.Printf("Mode: %d\n", getMode(inputArr))
// 		fmt.Printf("Median: %.3g\n", standardDeviation(inputArr))
// 	}
//     // fmt.Println(inputArgs)
// 	return
// }

package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var archiveDir string

func init() {
	flag.StringVar(&archiveDir, "a", "", "Directory to store archived log files")
	flag.Parse()
}

func main() {
	files := parseFiles()
	createDirectoryForArchives()
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			rotateLog(file)
		}(file)
	}
	wg.Wait()
}

func parseFiles() []string {
	if flag.NArg() == 0 {
		fmt.Println("Usage: ./myRotate [-a archiveDir] file1 [file2 ...]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	return flag.Args()
}

func rotateLog(file string) {
	info, err := os.Stat(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	timestamp := info.ModTime().Unix()
	outputFile := fmt.Sprintf("%s_%d.tar.gz", filepath.Base(file), timestamp)

	if archiveDir != "" {
		outputFile = filepath.Join(archiveDir, outputFile)
	}
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = out.Close() }()

	gw := gzip.NewWriter(out)
	defer func() { _ = gw.Close() }()

	tw := tar.NewWriter(gw)
	defer func() { _ = tw.Close() }()

	if err := addFileToArchive(file, tw); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Log file %s rotated and archived as %s\n", file, outputFile)
}

func createDirectoryForArchives() {
	if archiveDir != "" {
		if err := os.MkdirAll(archiveDir, 0755); err != nil {
			log.Fatal("ERROR: can't create directory", err)
		}
	}
}

func addFileToArchive(file string, tw *tar.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}

	header.Name = filepath.Base(file)
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := io.Copy(tw, f); err != nil {
		return err
	}

	return nil
}
