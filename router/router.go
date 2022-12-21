package router

import (
	"fmt"
	ctl "lecture/go-mvc/controller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	ct *ctl.Controller
}

func NewRouter(ctl *ctl.Controller) (*Router, error) {
	r := &Router{ct: ctl} //controller 포인터를 ct로 복사, 할당

	return r, nil
}

// cross domain을 위해 사용
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		//~ 생략
		c.Next()
	}
}

// 임의 인증을 위한 함수
func liteAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//~ 생략
		c.Next()
	}
}

// 실제 라우팅
func (p *Router) Idx() *gin.Engine {
	e := gin.Default()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Use(CORS())

	account := e.Group("acc/v01", liteAuth())
	{
		fmt.Println(account)
		account.GET("/ok", p.ct.GetOK) // controller 패키지의 실제 처리 함수
		account.GET("/persons", p.ct.GetPerson)
	}

	return e
}
