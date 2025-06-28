package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mirdha8846/c-hdfs.git/node1/encryption"
	"github.com/mirdha8846/c-hdfs.git/node1/types"
	
)


//Todo-uploading file name ..so update uploading function and think how to store their name so we access these files
type ChunkInfo struct {
    Name         string `json:"name"`
}
var fs = types.NewFileStore()
func main() {
	fmt.Println("Hello from Node 1!")
	Init()

	r := gin.Default()
	keyStr, err := encryption.GenrateKey()
		if err != nil{
			return
		}
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
		
		key, err := hex.DecodeString(keyStr) // Convert string key to []byte
		if err!=nil{
			c.JSON(400,gin.H{
				"message":"Internal Server Error",
			})
		}

		// Encrypt the file and store to encrypted folder
		encPath := filepath.Base(tempPath)
		encryptedPath := filepath.Join("encryptedFiles", encPath+".enc")
	
		go func(){

		err = encryption.EncryptFile(key, tempPath, encryptedPath)
		if err != nil {
			c.JSON(500, gin.H{"message": "encryption failed"})
			return
		}
		defer os.Remove(tempPath)
		
		
		chunkFiles, err := SplitFiles(encryptedPath, tempPath,3)
		if err != nil {
			c.JSON(500, gin.H{"message": "file splitting failed"})
			return
		}
		Uploading(chunkFiles)
		
		
		//remove chunking files
		defer func() {
			for _, chunkFile := range chunkFiles {
				os.Remove(chunkFile.Name)
			}
		}()
           
		
		defer os.Remove(encryptedPath)
		
		fs.AddFile(name, userFile.Filename)
}()
		c.JSON(200, gin.H{
			"message": "file uploaded, encrypted and split successfully",
			"key":     keyStr,
		
		})
	})

	r.POST("/api/getFiles",func(c *gin.Context) {
       userID:=c.PostForm("userID")
	   fileName:=c.PostForm("fileName")
       isExist := fs.GetFile(userID, fileName)
	    if !isExist {
			c.JSON(400,gin.H{
				"message":"no file found!!!",
			})
		}
       
		err,files:=GetFromNode(fileName)
		if err!=nil{
			c.JSON(400,gin.H{
				"message":"Internal server error!!!",
			})
		}
		tempFileName:=fmt.Sprintf("%s_%s",userID,fileName)
		completFile,err:=AddFiles(files,tempFileName)
		if err!=nil{
			c.JSON(400,gin.H{
				"message":err,
			})
		}

		key, err := hex.DecodeString(keyStr) // Convert string key to []byte
		if err!=nil{
			c.JSON(400,gin.H{
				"message":"Internal Server Error",
			})
		}
		finalFile,err:=encryption.DecryptFile(key,completFile.Reader)
		if err!=nil{
			c.JSON(400,gin.H{
				"message":"internal server error",
			})
		}

	
		c.JSON(200,gin.H{
			"File":finalFile,
		})
	})

	r.Run(":8082")
}