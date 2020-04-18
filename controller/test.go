package controller

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keepondream/article/common"
	"github.com/rs/xid"
	"github.com/thecodingmachine/gotenberg-go-client/v7"
)

func Test(c *gin.Context) {

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
	// req := gotenberg.NewURLRequest("https://da.hire66.com/pdf/52")
	// req.Margins(gotenberg.NoMargins)

	dest := filePath + newFileName + ".pdf"
	gc.Store(req, dest)

	// req.SetFlags(req.LstdFlags | req.Lcost) // 输出格式显示请求耗时
	// r, _ := req.Get("https://api.df5g.com/api/basicinfo")
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
