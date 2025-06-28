package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func Uploading(chunks []Chunk) (string,error) {
	for i, chunk := range chunks {
		firstNode := total_Servers[i%len(total_Servers)]
		secondNode := total_Servers[(i+1)%len(total_Servers)]
		
		if err := sendToNode(firstNode, chunk); err != nil {
			return "",err
		}
		if err := sendToNode(secondNode, chunk); err != nil {
			return "",err
		}
		defer chunk.Reader.Close()//why we need this and how this works ...
	}
	
	return nil
}

func sendToNode(nodeURI string, chunkData Chunk) error {
	chunkData.Reader.Seek(0, io.SeekStart)

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()

		part, err := writer.CreateFormFile("file", chunkData.Name)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		_, err = io.Copy(part, chunkData.Reader)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	// Prepare HTTP request
	req, err := http.NewRequest("POST", nodeURI+"/file/upload", pr)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}



func GetFromNode(fileName string) (error, []Chunk) {
	var filesPart []Chunk
	for i := 0; i < len(total_Servers); i++ {
		nodeURI := node_info[i][0]
		err, chunk := getFileFromNode(nodeURI, fileName)

		if err != nil {
			// Try fallback node
			nextNodeURI := node_info[i][1]
			err1, chunk := getFileFromNode(nextNodeURI, fileName)
			if err1 != nil {
				return err1,nil
			}
			filesPart = append(filesPart, chunk)
		} else {
			filesPart = append(filesPart, chunk)
		}
	}
	return nil, filesPart
}


func getFileFromNode(URI string, fileName string) (error, Chunk) {
	// 1. Make request
	req, err := http.NewRequest("GET", URI+"/getFile?name="+fileName, nil)
	if err != nil {
		return err, Chunk{}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, Chunk{}
	}
	defer resp.Body.Close()

	// 2. Create file on disk to store incoming data
	savePath := filepath.Join("downloadedChunks", fileName) // or wherever you want
	outFile, err := os.Create(savePath)
	if err != nil {
		return err, Chunk{}
	}

	// 3. Copy the response body directly into the file (streaming)
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		outFile.Close()
		return err, Chunk{}
	}

	// 4. Reopen file in read mode
	outFile.Close()
	readFile, err := os.Open(savePath)
	if err != nil {
		return err, Chunk{}
	}
    defer readFile.Close()
	// 5. Return Chunk with *os.File reader
	return nil, Chunk{
		Name:   fileName,
		Reader: readFile,
	}
}
