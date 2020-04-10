package controller

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/keepondream/article/common"
)

func TransformFile(c *gin.Context) {
	err := c.Request.ParseMultipartForm(200000)
	if err != nil {
		log.Fatal(err)
	}
	// 获取表单
	form := c.Request.MultipartForm
	// 获取参数upload后面的多个文件名,存放到数组files里面
	files := form.File["upload"]
	// 遍历数组,每去除一个file就拷贝一次
	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		fileName := files[i].Filename
		fmt.Println(fileName)
		// fileNames := strings.Split('.', fileName)
		// fmt.Println(fileNames)

		out, err := os.Create(fileName)
		defer out.Close()
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println(files)
	common.Success(c)
}
