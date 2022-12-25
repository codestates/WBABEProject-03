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

	papi := e.Group("api/v1", liteAuth())
	{
		//조회 - 모든고객
		papi.GET("/persons", p.ct.GetPersonAll)

		//조회 - name
		// papi.GET("/persons/:name", p.ct.GetPersonWithName)
		papi.GET("/persons/", p.ct.GetPersonWithName)

		//조회 - pnum
		papi.GET("/persons/:pnum", p.ct.GetPersonWithPnum)

		//신규인입 - name , age, pnum
		papi.POST("/persons", p.ct.NewPersonInsert)

		//삭제 - pnum
		papi.DELETE("/persons/:pnum", p.ct.DelPerson)

		//수정
		papi.PUT("/persons", p.ct.UpdatePerson)
	}

	menu := e.Group("api/v1")
	{
		// 전체메뉴조회
		menu.GET("/menus", p.ct.GetMenuAll)
		// 메뉴 등록
		menu.POST("/menus", p.ct.InsertMenu)
		// 메뉴 수정
		menu.PUT("/menus", p.ct.UpdateMenu)
		// 메뉴 삭제
		menu.DELETE("/menus/:name", p.ct.DelMenu)
		// 오늘의 메뉴 등록
		menu.POST("/menus/todays", p.ct.UpdateMenuOnTodayMenu)
	}
	order := e.Group("api/v1")
	{
		// 완료되지 않은 전체 주문 조회
		order.GET("/order", p.ct.GetOrdersWithoutDone)
		// 주문 등록
		order.POST("/order", p.ct.InsertOrder)

	}
	return e
}
