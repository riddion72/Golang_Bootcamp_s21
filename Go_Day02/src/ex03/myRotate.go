package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
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
	flag.StringVar(&archiveDir, "a", "", "`Directory` to store archived log files")
	flag.Parse()
}

func parseFiles() (_ []string, err error) {
	if flag.NArg() == 0 {
		fmt.Println("Usage: ./myRotate [-a archiveDir] file1 [file2 ...]")
		flag.PrintDefaults()
		err = errors.New("Something wrong with flag")
	}
	return flag.Args(), err
}

func createDirectory() (err error) {
	if archiveDir != "" {
		if err := os.MkdirAll(archiveDir, 0755); err != nil {
			err = errors.New("Can't create directory")
		}
	}
	return
}

func archivateLogs(path string) (err error) {
	fileInfo, lErr1 := os.Stat(path)
	if lErr1 != nil {
		return lErr1
	}
	timestamp := fileInfo.ModTime().Unix()
	outputFileName := fmt.Sprintf("%s_%d.tar.gz", filepath.Base(path), timestamp)

	outputFileName = filepath.Join(archiveDir, outputFileName)

	archive, lErr2 := os.Create(outputFileName)
	if lErr2 != nil {
		return lErr2
	}
	defer archive.Close()

	zipWriter := gzip.NewWriter(archive)
	defer zipWriter.Close()

	tarWriter := tar.NewWriter(zipWriter)
	defer tarWriter.Close()

	return archivator(path, tarWriter)
}

func archivator(file string, tW *tar.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:    filepath.Base(file),
		Size:    fileInfo.Size(),
		Mode:    int64(fileInfo.Mode()),
		ModTime: fileInfo.ModTime(),
	}

	if err := tW.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tW, f); err != nil {
		return err
	}

	return nil
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Unknown panic happend: ", err)
		}
	}()

	paths, err := parseFiles()
	if err != nil {
		log.Fatal(err)
	}

	err = createDirectory()
	if err != nil {
		log.Fatal(err)
	}

	var waitGroup sync.WaitGroup

	for _, filePath := range paths {
		waitGroup.Add(1)
		go func(path string) {
			defer waitGroup.Done()
			if err := archivateLogs(path); err != nil {
				log.Fatal(err)
			}
		}(filePath)
	}

	waitGroup.Wait()
}
