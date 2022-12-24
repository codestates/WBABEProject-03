package router

import (
	"fmt"

	"github.com/gin-gonic/gin"

	ctl "lecture/WBABEProject-03/controller"
)

type Router struct {
	ct *ctl.Controller
}

func NewRouter(ct *ctl.Controller) (*Router, error) {
	r := &Router{
		ct: ct,
	}

	return r, nil
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Forwarded-For, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func liteAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c == nil {
			c.Abort()
			return
		}
		auth := c.GetHeader("Authorization")
		fmt.Println("Authorization-word ", auth)

		c.Next()
	}
}

func (p *Router) Idx() *gin.Engine {
	e := gin.Default()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Use(CORS())

	e.GET("/health")

	papi := e.Group("person/v01", liteAuth())
	{
		//조회 - name
		papi.GET("/nser/:name", p.ct.GetPersonWithName)

		//조회 - pnum
		papi.GET("/pser/:pnum", p.ct.GetPersonWithPnum)
		fmt.Println("post")
		//신규인입 - name , age, pnum
		papi.POST("/ins", p.ct.NewPersonInsert)

		//삭제 - pnum
		papi.DELETE("/del/:pnum", p.ct.DelPerson)

		//수정
		papi.PUT("/upd", p.ct.UpdatePerson)
	}

	return e
}
