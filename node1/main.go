package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"path/filepath"
	
	"github.com/gin-gonic/gin"
	"github.com/mirdha8846/c-hdfs.git/node1/encryption"
	"github.com/mirdha8846/c-hdfs.git/node1/types"
)

type ChunkInfo struct {
    Name         string `json:"name"`
}

func main() {
	fmt.Println("Hello from Node 1!")

	r := gin.Default()

	// Upload Route
	r.POST("/api/fileUpload", func(c *gin.Context) {
		name := c.PostForm("userID")
		userFile, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"message": "failed to read file"})
			return
		}

		// Save file to temp folder
		tempPath := filepath.Join("temp", fmt.Sprintf("%s_%d", userFile.Filename, rand.Int63()))
		err = c.SaveUploadedFile(userFile, tempPath)
		if err != nil {
			c.JSON(500, gin.H{"message": "failed to save file to temp"})
			return
		}

		// Generate AES key
		keyStr, err := encryption.GenrateKey()
		if err != nil {
			c.JSON(500, gin.H{"message": "failed to generate key"})
			return
		}
		key, _ := hex.DecodeString(keyStr) // Convert string key to []byte

		// Encrypt the file and store to encrypted folder
		encPath := filepath.Base(tempPath)
		encryptedPath := filepath.Join("encryptedFiles", encPath+".enc")
		err = encryption.EncryptFile(key, tempPath, encryptedPath)
		if err != nil {
			c.JSON(500, gin.H{"message": "encryption failed"})
			return
		}

		// Split the encrypted file
		chunkFiles, err := SplitFiles(encryptedPath, 3)
		if err != nil {
			c.JSON(500, gin.H{"message": "file splitting failed"})
			return
		}

		// Convert chunks to include base64 content
		var chunks []ChunkInfo
		for _, chunkFile := range chunkFiles {
			// Read file content
		
			
			chunks = append(chunks, ChunkInfo{
				Name:    chunkFile.Name,
				
			})
		}

		// Store in memory (example use of your types.FileStore)
		fs := types.NewFileStore()
		fs.AddFile(name, userFile.Filename)

		c.JSON(200, gin.H{
			"message": "file uploaded, encrypted and split successfully",
			"key":     keyStr,
			"chunks":  chunks,
		})
	})

	r.Run(":8082")
}