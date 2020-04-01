package controller

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

func Test(c *gin.Context) {

	req.SetFlags(req.LstdFlags | req.Lcost) // 输出格式显示请求耗时
	r, _ := req.Get("https://api.df5g.com/api/basicinfo")
	log.Println(r)
	if r.Cost() > 3*time.Second {
		log.Println("WARN: slow request:", r)
	}

	n := 32
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			req.SetFlags(req.LstdFlags | req.Lcost) // 输出格式显示请求耗时
			r, _ := req.Get("https://api.df5g.com/api/basicinfo")
			log.Println(r)
			if r.Cost() > 3*time.Second {
				log.Println("WARN: slow request:", r)
			}
		}()
	}

	// common.Failed(c)
	// common.Failed(c)
	// common.Success(c, common.WithCode(403), common.WithMsg("坤哥牛逼"))
}
