package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/keepondream/article/common"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/thecodingmachine/gotenberg-go-client/v7"
)

func Test(c *gin.Context) {
	fmt.Println()
	fmt.Println("start--------------")
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("ctx.Request.body: %v", string(data))
	fmt.Println("end---------------")
	common.Failed(c)
	return

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose("hello.pdf")

	// pdf := gofpdf.New("P", "mm", "A4", "")
	// pdf.AddPage()
	// pdf.SetFont("Arial", "", 11)
	// pdf.Image("test.png", 10, 10, 30, 0, false, "", 0, "")
	// pdf.Text(50, 20, "test.png")
	// pdf.Image("test.gif", 10, 40, 30, 0, false, "", 0, "")
	// pdf.Text(50, 50, "test.gif")
	// pdf.Image("test.jpg", 10, 130, 30, 0, false, "", 0, "")
	// pdf.Text(50, 140, "test.jpg")

	// err := pdf.OutputFileAndClose("write_pdf_with_image.pdf")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	common.Success(c)
	return

	//Info级别的日志
	common.Logger().WithFields(logrus.Fields{
		"name": "hanyun",
	}).Info("记录一下日志", "Info")
	//Error级别的日志
	common.Logger().WithFields(logrus.Fields{
		"name": "hanyun",
	}).Error("记录一下日志", "Error")
	//Warn级别的日志
	common.Logger().WithFields(logrus.Fields{
		"name": "hanyun",
	}).Warn("记录一下日志", "Warn")
	//Debug级别的日志
	common.Logger().WithFields(logrus.Fields{
		"name": "hanyun",
	}).Debug("记录一下日志", "Debug")

	common.Success(c)

	return
	basePath := common.GetBasePath()
	filePath := basePath + "/libreoffice/"
	newFileId := xid.New()

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
	oldFileExt := oldFileNameArr[len(oldFileNameArr)-1]
	oldFileExt = strings.ToLower(oldFileExt)
	oldFileName := strings.Join(oldFileNameArr[0:len(oldFileNameArr)-1], ".")
	newFileName := oldFileName + "_00_" + newFileId.String()
	newFileFullName := newFileName + "." + oldFileExt
	fmt.Println(oldFileExt, oldFileName)
	fmt.Println(file, err)

	out, err := os.Create(filePath + newFileFullName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(newFileFullName, filePath)

	gc := &gotenberg.Client{Hostname: "http://localhost:3000"}

	doc, _ := gotenberg.NewDocumentFromPath("document.docx", filePath+newFileFullName)
	req := gotenberg.NewOfficeRequest(doc)
	// doc, _ := gotenberg.NewDocumentFromPath(newFileFullName, filePath+newFileFullName)
	// // doc2, _ := gotenberg.NewDocumentFromPath("document2.docx", "/path/to/file")
	// req := gotenberg.NewOfficeRequest(doc)

	//# html -> pdf
	// req := gotenberg.NewURLRequest("https://www.baidu.com")
	// req.Margins(gotenberg.NoMargins)

	dest := filePath + newFileName + ".pdf"
	gc.Store(req, dest)

	// req.SetFlags(req.LstdFlags | req.Lcost) // 输出格式显示请求耗时
	// r, _ := req.Get("https://www.baidu.com")
	// log.Println(r)
	// if r.Cost() > 3*time.Second {
	// 	log.Println("WARN: slow request:", r)
	// }

	// n := 11132
	// var wg sync.WaitGroup
	// wg.Add(n)
	// for i := 0; i < n; i++ {
	// 	go func() {
	// 		req.SetFlags(req.LstdFlags | req.Lcost) // 输出格式显示请求耗时
	// 		r, _ := req.Get("https://www.baidu.com")
	// 		log.Println(r)
	// 		// if r.Cost() > 3*time.Second {
	// 		// log.Println("WARN: slow request:", r)
	// 		// }
	// 	}()
	// }

	// common.Failed(c)
	// common.Failed(c)
	common.Success(c, common.WithCode(403), common.WithMsg("坤哥牛逼"))
}
