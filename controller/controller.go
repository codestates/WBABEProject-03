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

// NewPersonInsert godoc
// @Summary call NewPersonInsert, return ok by json.
// @Description 모든 메뉴 조회
// @name NewPersonInsert
// @Accept json
// @Produce json
// @Param : body Controller true "new Person"
// @Router /api/v1/persons [post]
// @Success 200 {object} Controller
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

// GetPersonAll godoc
// @Summary call GetPersonAll, return ok by json.
// @Description 모든 고객 조회 기능
// @name GetPersonAll
// @Accept json
// @Produce json
// @Router /api/v1/persons [get]
// @Success 200 {object} Controller
func (p *Controller) GetPersonAll(c *gin.Context) {
	p.md.GetPerson()
	c.JSON(http.StatusOK, gin.H{
		"res":  "ok",
		"body": "pers",
	})
	c.Next()

}

// GetMenuAll godoc
// @Summary call GetMenuAll, return ok by json.
// @Description 모든 메뉴 조회
// @name GetMenuAll
// @Accept json
// @Produce json
// @Router /api/v1/menus [get]
// @Success 200 {object} Controller
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

type menuDto struct {
	Name    string `json:"name"`
	IsToday bool   `json:"isToday"`
}

func (p *Controller) UpdateMenuOnTodayMenu(c *gin.Context) {
	var menuDto menuDto
	fmt.Println(menuDto)
	if err := c.ShouldBindJSON(&menuDto); err == nil {
		fmt.Printf("menu dto - %+v \n", menuDto)
	} else {
		fmt.Printf("error - %+v \n", err)
	}

	if err := p.md.UpdateMenuOnTodayMenu(menuDto.Name, menuDto.IsToday); err != nil {
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

// GetOrdersWithoutDone godoc
// @Summary call GetOrdersWithoutDone, return ok by json.
// @Description 모든 주문 조회 (완료 제외)
// @name GetOrdersWithoutDone
// @Accept json
// @Produce json
// @Router /api/v1/order [get]
// @Success 200 {object} Controller
func (p *Controller) GetOrdersWithoutDone(c *gin.Context) {
	if orders, err := p.md.FindAllOrderWithoutDone(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
		c.Abort()
	} else {
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": orders,
		})
		c.Next()
	}
}

type orderDto struct {
	ClientName string   `json:"clientName" bson:"clientName"`
	MenuList   []string `json:"menuList" bson:"menuList"`
}

// InsertOrder godoc
// @Summary call InsertOrder, return ok by json.
// @Description 모든 메뉴 조회
// @name InsertOrder
// @Accept json
// @Produce json
// @Param : body orderDto true "new order"
// @Router /api/v1/order [post]
// @Success 200 {object} Controller
func (p *Controller) InsertOrder(c *gin.Context) {
	var orderDto orderDto

	if err := c.ShouldBindJSON(&orderDto); err == nil {
		fmt.Printf("order dto - %+v \n", orderDto)
	} else {
		fmt.Printf("error - %+v \n", err)
	}

	if len(orderDto.MenuList) <= 0 {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", nil)
		return
	}
	var totalPrice int
	for _, menu := range orderDto.MenuList {
		tMenu, _ := p.md.GetOneMenu(menu)
		if tMenu.Name == "" {
			p.RespError(c, nil, http.StatusUnprocessableEntity, "Can Not Find Menu", nil)
			return
		}
		totalPrice += tMenu.Price
	}

	req := model.Order{
		ClientName: orderDto.ClientName,
		MenuList:   orderDto.MenuList,
		Status:     "접수중",
		TotalPrice: totalPrice,
	}

	if err := p.md.CreateOrder(req); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}

type orderUpdateMenuDto struct {
	Pnum           string   `json:"pnum" bson:"pnum"`
	UpdateMenuList []string `json:"updateMenuList" bson:"updateMenuList"`
}

func (p *Controller) UpdateOrderMenu(c *gin.Context) {
	var orderUpdateMenuDto orderUpdateMenuDto
	if err := c.ShouldBindJSON(&orderUpdateMenuDto); err == nil {
		fmt.Printf("orderUpdate dto - %+v \n", orderUpdateMenuDto)
	} else {
		fmt.Printf("error - %+v \n", err)
	}

	var totalPrice int
	for _, menu := range orderUpdateMenuDto.UpdateMenuList {
		tMenu, _ := p.md.GetOneMenu(menu)
		if tMenu.Name == "" {
			p.RespError(c, nil, http.StatusUnprocessableEntity, "Can Not Find Menu", nil)
			return
		}
		totalPrice += tMenu.Price
	}

	nPnum, err := strconv.Atoi(orderUpdateMenuDto.Pnum)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "pnum error", err)
		return
	}

	if err := p.md.UpdateOrderMenu(nPnum, orderUpdateMenuDto.UpdateMenuList, totalPrice); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}

type orderUpdateStatusDto struct {
	Pnum         string `json:"pnum" bson:"pnum"`
	UpdateStatus string `json:"updateStatus" bson:"updateStatus"`
}

func (p *Controller) UpdateOrderStatus(c *gin.Context) {
	var orderUpdateStatusDto orderUpdateStatusDto
	if err := c.ShouldBindJSON(&orderUpdateStatusDto); err == nil {
		fmt.Printf("orderUpdate dto - %+v \n", orderUpdateStatusDto)
	} else {
		fmt.Printf("error - %+v \n", err)
	}

	nPnum, err := strconv.Atoi(orderUpdateStatusDto.Pnum)
	if err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "pnum error", err)
		return
	}

	if err := p.md.UpdateOrderStatus(nPnum, orderUpdateStatusDto.UpdateStatus); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}

type reviewDto struct {
	OrderId   string         `json:"orderId" bson:"orderId"`
	MenuScore map[string]int `json:"menuScore" bson:"menuScore"`
}

func (p *Controller) InsertReview(c *gin.Context) {
	var reviewDto reviewDto
	if err := c.ShouldBindJSON(&reviewDto); err == nil {
		fmt.Printf("review dto - %+v \n", reviewDto)
	} else {
		fmt.Printf("error - %+v \n", err)
	}

	for key := range reviewDto.MenuScore {
		menu, _ := p.md.GetOneMenu(key)
		if menu.Name == "" {
			p.RespError(c, nil, http.StatusUnprocessableEntity, "Can not find Menu", nil)
			return
		}
	}

	if err := p.md.UpdateReviewInOrder(reviewDto.OrderId, reviewDto.MenuScore); err != nil {
		p.RespError(c, nil, http.StatusUnprocessableEntity, "parameter not found", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
	c.Next()
}
