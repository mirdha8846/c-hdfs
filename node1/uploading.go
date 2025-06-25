package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	// "os"
)



var total_Servers = []string{
	"http://localhost:4005",
	"http://localhost:4006",
	"http://localhost:4007",
}

func Uploading(chunks []Chunk) error {
	for i, chunk := range chunks {
		firstNode := total_Servers[i%len(total_Servers)]
		secondNode := total_Servers[(i+1)%len(total_Servers)]

		if err := sendToNode(firstNode, chunk); err != nil {
			return err
		}
		if err := sendToNode(secondNode, chunk); err != nil {
			return err
		}
	}
	return nil
}

func sendToNode(nodeURI string, chunkData Chunk) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Correct method: for files
	part, err := writer.CreateFormFile("file", chunkData.Name)
	if err != nil {
		return err
	}

	// Copy file content into part
	if _, err := io.Copy(part, chunkData.Reader); err != nil {
		return err
	}

	// Close the writer to finalize the multipart body
	err = writer.Close()
	if err != nil {
		return err
	}

	// Prepare HTTP request
	req, err := http.NewRequest("POST", nodeURI+"/file/upload", body)
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
