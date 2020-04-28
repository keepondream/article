package controller

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keepondream/article/common"
	"github.com/rs/xid"
)

func TransformFile(c *gin.Context) {
	// 获取HTML 连接
	url := c.PostForm("url")
	oldFileExt := "html"
	oldFileName := "html"
	shellCommand := ""
	transformFileName := ""
	// 3.获取需要转换的格式
	transformExt := c.PostForm("ext")
	if transformExt == "" {
		common.Failed(c, common.WithMsg("请求参数有误!"))
		return
	}
	transformExt = strings.ToLower(transformExt)

	// 4.获取项目路劲,并检测创建文件上传目录
	basePath := common.GetBasePath()
	filePath := basePath + "/libreoffice/"

	if !common.Exists(filePath) {
		common.Mkdir(filePath)
	}

	// 5.将文件挪入libreoffice,并给一个唯一的新文件名称
	newFileId := xid.New()
	newFileName := ""
	newFileFullName := ""

	if url == "" {
		// 1.获取文件
		// FormFile方法会读取参数“file”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			common.Failed(c, common.WithMsg("文件获取失败"))
			return
		}

		// 2.获取文件名称,文件地址,文件格式
		// header调用Filename方法，就可以得到文件名
		oldFullFilename := header.Filename
		oldFullFilename = strings.Replace(oldFullFilename, " ", "", -1)
		oldFullFilename = strings.Replace(oldFullFilename, "　", "", -1)
		oldFullFilename = strings.Replace(oldFullFilename, "\\", "", -1)
		oldFullFilename = strings.Replace(oldFullFilename, "(", "", -1)
		oldFullFilename = strings.Replace(oldFullFilename, ")", "", -1)
		oldFullFilename = strings.Replace(oldFullFilename, "&", "", -1)
		oldFullFilename = strings.Replace(oldFullFilename, "<", "", -1)
		oldFullFilename = strings.Replace(oldFullFilename, ">", "", -1)
		oldFullFilename = strings.Replace(oldFullFilename, "?", "", -1)
		oldFileNameArr := strings.Split(oldFullFilename, ".")
		oldFileExt = oldFileNameArr[len(oldFileNameArr)-1]
		oldFileExt = strings.ToLower(oldFileExt)
		oldFileName = strings.Join(oldFileNameArr[0:len(oldFileNameArr)-1], ".")
		newFileName = oldFileName + "_00_" + newFileId.String()
		newFileFullName = newFileName + "." + oldFileExt
		fmt.Println(oldFileExt, oldFileName)
		fmt.Println(file, err)

		fmt.Println(url, transformExt)

		out, err := os.Create(filePath + newFileFullName)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		oldFileExt = "html"
		newFileName = oldFileName + "_00_" + newFileId.String()
		newFileFullName = newFileName + "." + oldFileExt
		oldFileName = newFileName
	}
	fmt.Println(transformExt, url)

	// 被转换文件的名称
	transformFileName = newFileName + "." + transformExt
	// 下载用的文件名
	downloadFileName := ""
	// 根据需要转换的类型,进行文件docker命令生成
	// 1. 	pdf -> 图片					  	 √
	//    	pdf <- 图片                    	 √
	// 2. 	pdf -> word				    	×
	//    	pdf <- word                   	√
	// 3. 	pdf -> excel 					×
	//    	pdf <- excel                  	×
	// 4. 	pdf -> html                   	√
	//    	pdf <- html                   	√
	// 5. 	图片 -> word					 ×
	//	  	图片 <- word                   	 ×
	// 6. 	图片 -> excel                  	 ×
	//    	图片 <- excel                  	 ×
	// 7. 	图片 -> html                   	 ×
	//    	图片 <- html                   	 ×
	// 8. 	word -> excel                 	×
	//    	word <- excel 					×
	// 9. 	word -> html 					×
	//	  	word <- html 					×
	// 10.	excel -> html					×
	//      excel <- html 					×
	switch transformExt {
	case "pdf":
		// 转换成PDF
		fmt.Println("pdf")
		// 图片 -> pdf
		if (oldFileExt == "png") || (oldFileExt == "jpg") || (oldFileExt == "jpeg") {
			shellCommand = "docker exec imagemagick_1 /bin/bash -c 'convert " + newFileFullName + " " + transformFileName + "'"
			downloadFileName = transformFileName
		}
		// word -> pdf
		if (oldFileExt == "docx") || (oldFileExt == "doc") {
			shellCommand = "docker exec libreoffice_1 /bin/bash -c 'soffice --headless --invisible --convert-to pdf " + newFileFullName + "'"
			downloadFileName = transformFileName
		}
		// html -> pdf
		if (oldFileExt == "html") && (url != "") {
			// shellCommand = "docker run -i --rm -v " + filePath + ":/root icalialabs/wkhtmltopdf " + url + " " + filePath + transformFileName
			shellCommand = "docker exec wkhtmltopdf_1 /bin/bash -c 'wkhtmltopdf -s A4  " + url + " " + transformFileName + "'"
			downloadFileName = transformFileName
		}
	case "jpg", "jpeg", "png":
		// 转换成图片
		fmt.Println("jjp")
		// pdf -> 图片
		if oldFileExt == "pdf" {
			shellCommand = "docker exec imagemagick_1 /bin/bash -c 'convert -density 300 -background white  -alpha remove -append " + newFileFullName + " " + transformFileName + "'"
			downloadFileName = transformFileName
		}
	case "docx", "doc":
		// 转换成word
		fmt.Println("to word")
		// pdf
		if oldFileExt == "pdf" {
			if transformExt == "doc" {
				shellCommand = "docker exec libreoffice_1 /bin/bash -c 'soffice --headless --convert-to " + transformExt + ":\"MS Word 2007 XML\" " + newFileFullName + "'"
			}
			if transformExt == "docx" {
				shellCommand = "docker exec libreoffice_1 /bin/bash -c 'soffice --headless --convert-to " + transformExt + ":\"Microsoft Word 2007/2010/2013 XML\" " + newFileFullName + "'"
			}
			downloadFileName = transformFileName
			shellCommand = ""
			// pdf 转 doc 和 docx 不行
		}
	case "xlsx":
		// 转成Excel
		fmt.Println("to excel")
	case "html":
		// 转成html
		// pdf -> html
		if oldFileExt == "pdf" {
			// shellCommand = "docker run -i --rm -v " + filePath + ":/pdf bwits/pdf2htmlex-alpine pdf2htmlEX --zoom 1.3 " + newFileFullName
			shellCommand = "docker run -i --rm -v " + filePath + ":/pdf bwits/pdf2htmlex-alpine pdf2htmlEX --zoom 1.3 " + newFileFullName
			downloadFileName = transformFileName
		}
	default:
		common.Failed(c, common.WithMsg("暂时不支持该类型转换,请等待后续开发"))
		return
	}

	if shellCommand == "" {
		common.Failed(c, common.WithMsg("暂时不支持该类型转换,请等待后续开发"))
		return
	}
	fmt.Println(shellCommand)

	// Golang 执行 shell脚本，并实时打印 shell 脚本输出日志信息
	// 实际业务比如：异步任务调度系统、自动发布系统等都有可能需要 shell 脚本的配合来完成，就需要实时打印 shell 脚本的中每条命令的输出日志信息，便于查看任务进度等
	// 这里遇到一个坑就是  错误:the input device is not a TTY  去掉 docker 命令中 -it  的 t
	// logfile := " >> " + filePath + "file.log"
	logfile := " "
	cmd := exec.Command("sh", "-c", shellCommand+logfile)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var outfile bytes.Buffer
	cmd.Stdout = &outfile

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()

	fmt.Println(outfile.String(), err)

	if err != nil {
		common.Failed(c, common.WithMsg("文件转换失败,请稍后重试"))
		return
	}

	res := make(map[string]interface{})
	res["filefullname"] = downloadFileName
	res["filename"] = oldFileName

	common.Success(c, common.WithData(res))
	return
}

func asyncLog(reader io.ReadCloser) error {
	cache := "" //缓存不足一行的日志信息
	buf := make([]byte, 1024)
	for {
		num, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if num > 0 {
			b := buf[:num]
			s := strings.Split(string(b), "\n")
			line := strings.Join(s[:len(s)-1], "\n") //取出整行的日志
			fmt.Printf("%s%s\n", cache, line)
			cache = s[len(s)-1]
		}
	}
	return nil
}

func DownloadFile(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		common.Failed(c, common.WithMsg("请求参数有误!"))
	}
	filePath := path.Join(common.GetBasePath(), "libreoffice")
	fmt.Println(fileName, filePath)
	// 打开文件
	fileFullPath := path.Join(filePath, fileName)
	file, err := os.Open(fileFullPath)
	if err != nil {
		common.Failed(c, common.WithMsg("文件不存在"))
		return
	}
	// 结束后关闭文件
	defer file.Close()

	// 设置响应的header头

	c.Writer.Header().Add("Content-type", "application/octet-stream")
	c.Writer.Header().Add("content-disposition", "attachment; filename="+fileName)
	c.File(fileFullPath)
}
