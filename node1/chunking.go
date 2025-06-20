package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func SplitFiles(filePath string, n int) ([]string, error) {
	mainFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer mainFile.Close()

	fileInfo, err := mainFile.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	chunkSize := fileSize / int64(n)

	chunkingFiles := []string{}
	baseName := filepath.Base(filePath)

	for i := 0; i < n; i++ {
		partFileName := fmt.Sprintf("%s.part%d", baseName, i+1)
		tempFile, err := os.Create(partFileName)//this part is saving all three parts on our node also
		if err != nil {
			return nil, err
		}

		if i == n-1 {
			// Last part – copy all remaining bytes
			_, err = io.Copy(tempFile, mainFile)
		} else {
			_, err = io.CopyN(tempFile, mainFile, chunkSize)
		}

		tempFile.Close()

		if err != nil && err != io.EOF {
			return nil, err
		}

		chunkingFiles = append(chunkingFiles, partFileName)
	}

	return chunkingFiles, nil
}

