package controller

import (
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
	oldFileNameArr := strings.Split(oldFullFilename, ".")
	oldFileExt := oldFileNameArr[len(oldFileNameArr)-1]
	oldFileName := strings.Join(oldFileNameArr[0:len(oldFileNameArr)-1], ".")
	fmt.Println(oldFileExt, oldFileName)
	fmt.Println(file, err)
	// 3.获取需要转换的格式
	transformExt := c.PostForm("ext")
	if transformExt == "" {
		common.Failed(c, common.WithMsg("请求参数有误!"))
		return
	}

	// 4.获取项目路劲,并检测创建文件上传目录
	basePath := common.GetBasePath()
	filePath := basePath + "/libreoffice/"

	if !common.Exists(filePath) {
		common.Mkdir(filePath)
	}

	// 5.将文件挪入libreoffice,并给一个唯一的新文件名称
	newFileId := xid.New()
	newFileName := oldFileName + "_00_" + newFileId.String()
	newFileFullName := newFileName + "." + oldFileExt
	fmt.Println(newFileFullName)
	out, err := os.Create(filePath + newFileFullName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	shellCommand := ""
	transformFileName := ""
	// 根据需要转换的类型,进行文件docker命令生成
	switch transformExt {
	case "pdf":
		// 转换成PDF
		fmt.Println("pdf")
		// 被转换文件的名称
		transformFileName = newFileName + ".pdf"
		if (oldFileExt == "png") || (oldFileExt == "jpg") || (oldFileExt == "jpeg") {
			shellCommand = "docker exec imagemagick_1 /bin/bash -c 'convert " + newFileFullName + " " + transformFileName + "'"
			fmt.Println(shellCommand)
		}
	case "jpg", "jpeg", "png":
		// 转换成图片
		fmt.Println("jjp")
	case "doc", "docx":
		// 转换成word
		fmt.Println("to word")
	default:
		common.Failed(c, common.WithMsg("不支持该类型"))
		return
	}

	if shellCommand == "" {
		common.Failed(c, common.WithMsg("不支持该类型"))
		return
	}

	logfile := " >> " + filePath + "file.log"
	cmd := exec.Command("sh", "-c", shellCommand+logfile)

	transformOut, err := cmd.Output()

	fmt.Println(string(transformOut), err)

	common.Success(c)
	return

	// 调用docker容器,根据需要转换的格式进行处理

	// 将处理好的文件挪入服务器下载目录,并还原新的名称

	// 返出组装好的下载链接

}
