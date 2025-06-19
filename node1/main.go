package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mirdha8846/c-hdfs.git/node1/types"
)

func main() {
	fmt.Println("Hello from Node 1!")
	r:=gin.Default()
	r.POST("/api/fileUpload",func(c *gin.Context) {
		name:=c.PostForm("userID")
		userfile,err:=c.FormFile("file")
		if err!=nil{
			c.JSON(400,gin.H{
				"message":"internal server error",
			})
			return
		}
		fs:=types.NewFileStore()
		fs.AddFile(name,userfile.Filename)
		//now first upload at temp folder and then access from temp and then
		filePath:=filepath.Join("temp",userfile.Filename)
		err=c.SaveUploadedFile(userfile,filePath)
		if err!=nil{
				c.JSON(400,gin.H{
				"message":"internal server error",
			})
			return
		}
		file,err:=os.Open(filePath)
		if err!=nil{
				c.JSON(400,gin.H{
				"message":"internal server error",
			})
			return
		}
		
		//encrypt it and then chunking part and then uploading on all servers
		




	})

}