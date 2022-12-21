package controller

import (
	"lecture/go-mvc/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	md *model.Model
}

func NewCTL(rep *model.Model) (*Controller, error) {
	r := &Controller{md: rep}
	return r, nil
}

func (p *Controller) RespOK(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
}

func (p *Controller) GetOK(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "ok"})
	return
}

func (p *Controller) GetPerson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"res":  "ok",
		"data": p.md.GetPerson(),
	})
}
