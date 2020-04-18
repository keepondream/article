package controller

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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
	}
	fmt.Println(transformExt, url)

	// 被转换文件的名称
	transformFileName = newFileName + "." + transformExt
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
		}
		// word -> pdf
		if (oldFileExt == "docx") || (oldFileExt == "doc") {
			shellCommand = "docker exec libreoffice_1 /bin/bash -c 'soffice --headless --invisible --convert-to pdf " + newFileFullName + "'"
			shellCommand = ""
		}
		// html -> pdf
		if (oldFileExt == "html") && (url != "") {
			// shellCommand = "docker run -i --rm -v " + filePath + ":/root icalialabs/wkhtmltopdf " + url + " " + filePath + transformFileName
			shellCommand = "docker exec wkhtmltopdf_1 /bin/bash -c 'wkhtmltopdf -s A4  " + url + " " + transformFileName + "'"
		}
	case "jpg", "jpeg", "png":
		// 转换成图片
		fmt.Println("jjp")
		// pdf -> 图片
		if oldFileExt == "pdf" {
			shellCommand = "docker exec imagemagick_1 /bin/bash -c 'convert -density 300 -background white  -alpha remove -append " + newFileFullName + " " + transformFileName + "'"
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
		}
	default:
		common.Failed(c, common.WithMsg("不支持该类型"))
		return
	}

	if shellCommand == "" {
		common.Failed(c, common.WithMsg("不支持该类型"))
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

	// stdout, _ := cmd.StdoutPipe()
	// stderr, _ := cmd.StderrPipe()

	// if err := cmd.Start(); err != nil {
	// 	log.Printf("Error starting command: %s......", err.Error())
	// 	return
	// }

	// go asyncLog(stdout)
	// go asyncLog(stderr)

	// if err := cmd.Wait(); err != nil {
	// 	log.Printf("Error waiting for command execution: %s......", err.Error())
	// 	return
	// }

	// transformOut, err := cmd.Output()

	// fmt.Println(string(transformOut), err)

	common.Success(c)
	return

	// 调用docker容器,根据需要转换的格式进行处理

	// 将处理好的文件挪入服务器下载目录,并还原新的名称

	// 返出组装好的下载链接

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
