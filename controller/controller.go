package controller

import (
	"encoding/json"
	"fmt"
	"lecture/WBABEProject-03/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	md *model.Model
}

func NewCTL(rep *model.Model) (*Controller, error) {
	r := &Controller{
		md: rep,
	}

	return r, nil
}

func (p *Controller) Check(c *gin.Context) {
	p.RespOK(c, 0)
}

func (p *Controller) RespOK(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
}

func (p *Controller) RespError(c *gin.Context, body interface{}, status int, err ...interface{}) {
	bytes, _ := json.Marshal(body)

	fmt.Println("Request error", "path", c.FullPath(), "body", bytes, "status", status, "error", err)

	c.JSON(status, gin.H{
		"Error":  "Request Error",
		"path":   c.FullPath(),
		"body":   bytes,
		"status": status,
		"error":  err,
	})
	c.Abort()
}

func (p *Controller) GetOK(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
}

func (p *Controller) GetPersonWithName(c *gin.Context) {
	fmt.Println(c.ClientIP())
	// sName := c.Param("name")
	sName := c.Query("name")
	fmt.Println(sName)
	if len(sName) <= 0 {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}

	if per, err := p.md.GetOnePerson("name", sName); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": per,
		})
	}
}

func (p *Controller) GetPersonWithPnum(c *gin.Context) {
	fmt.Println(c.ClientIP())

	sPnum := c.Param("pnum")
	if len(sPnum) <= 0 {
		p.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}

	if per, err := p.md.GetOnePerson("pnum", sPnum); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
		c.Abort()
	} else {
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": per,
		})
		c.Next()
	}
}

func (p *Controller) NewPersonInsert(c *gin.Context) {
	name := c.PostForm("name")
	sAge := c.PostForm("age")
	spnum := c.PostForm("pnum")

	if len(name) <= 0 || len(spnum) <= 0 {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
		return
	}

	per, _ := p.md.GetOnePerson("pnum", spnum)
	if per != (model.Person{}) {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "already resistery person", nil)
		return
	}

	nAge, err := strconv.Atoi(sAge)
	if err != nil {
		nAge = 1
	}

	req := model.Person{Name: name, Age: nAge, Pnum: spnum}
	if err := p.md.CreatePerson(req); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}

func (p *Controller) DelPerson(c *gin.Context) {
	spnum := c.Param("pnum")

	if len(spnum) <= 0 {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
		return
	}

	_, err := p.md.GetOnePerson("pnum", spnum)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "exist resistery person", nil)
		return
	}

	if err := p.md.DeletePerson(spnum); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail delete db", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}

func (p *Controller) UpdatePerson(c *gin.Context) {
	sAge := c.PostForm("age")
	spnum := c.PostForm("pnum")

	if len(spnum) <= 0 {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
		return
	}

	per, _ := p.md.GetOnePerson("pnum", spnum)
	fmt.Println("res ", per)
	if per == (model.Person{}) {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "could not found person", nil)
		return
	}

	nAge, err := strconv.Atoi(sAge)
	if err != nil {
		nAge = 1
	}

	if err := p.md.UpdatePerson(spnum, nAge); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}
func (p *Controller) GetPersonAll(c *gin.Context) {
	p.md.GetPerson()
	c.JSON(http.StatusOK, gin.H{
		"res":  "ok",
		"body": "pers",
	})
	c.Next()

}
func (p *Controller) GetMenuAll(c *gin.Context) {
	if menus, err := p.md.FindAllMenu(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
		c.Abort()
	} else {
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": menus,
		})
		c.Next()
	}
}

func (p *Controller) InsertMenu(c *gin.Context) {
	name := c.PostForm("name")
	sPrice := c.PostForm("price")

	if len(name) <= 0 {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
		return
	}

	menu, _ := p.md.GetOneMenu(name)
	if menu != (model.Menu{}) {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "already resistery Menu", nil)
		return
	}

	nPrice, err := strconv.Atoi(sPrice)
	if err != nil {
		nPrice = 0
	}

	req := model.Menu{Name: name, Price: nPrice}
	if err := p.md.CreateMenu(req); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}
func (p *Controller) UpdateMenu(c *gin.Context) {
	name := c.PostForm("name")
	sPrice := c.PostForm("price")

	if len(sPrice) <= 0 {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
		return
	}

	menu, _ := p.md.GetOneMenu(name)
	fmt.Println("res ", menu)
	if menu == (model.Menu{}) {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "could not found menu", nil)
		return
	}

	nPrice, err := strconv.Atoi(sPrice)
	if err != nil {
		nPrice = 0
	}

	if err := p.md.UpdateMenu(name, nPrice); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}
func (p *Controller) DelMenu(c *gin.Context) {
	name := c.Param("name")

	_, err := p.md.GetOneMenu(name)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "exist resistery menu", nil)
		return
	}

	if err := p.md.DeleteMenu(name); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "fail delete db", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}
