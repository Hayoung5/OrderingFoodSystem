package controller

import (
	"fmt"
	model "lecture/go-project/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	md *model.Model
}

type UpdateMenuStruct struct {
	MenuName string            `bson:"menuName"`
	Update   map[string]string `bson:"update"`
}

func NewCTL(rep *model.Model) (*Controller, error) {
	r := &Controller{md: rep}
	return r, nil
}

func (p *Controller) GetOK(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "ok"})
	return
}

func (p *Controller) GetMenu(c *gin.Context) {
	// sort?q=latest(최신),reorder(재주문),highest(별점높은순),reco(금일추천메뉴)
	r, _ := model.NewModel()
	query := c.Request.URL.Query().Get("q")
	if query == "latest" || query == "reorder" || query == "highest" || query == "reco" || query == "" {
		c.JSON(200, r.GetMenu(query))
	} else {
		c.JSON(400, "Check you query parameter")
	}
	return
}

func (p *Controller) PostMenu(c *gin.Context) {
	r, _ := model.NewModel()

	var newMenu model.Menu
	if err := c.BindJSON(&newMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	r.PostMenu(newMenu)
	c.JSON(200, gin.H{"respoense": newMenu})
	return
}

func (p *Controller) UpdateMenu(c *gin.Context) {
	r, _ := model.NewModel()

	var request UpdateMenuStruct

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println(request)

	r.UpdateMenu(request.MenuName, request.Update)
	c.JSON(200, gin.H{"respoense": `the menu is changed`})
	return
}

func (p *Controller) DeleteMenu(c *gin.Context) {
	r, _ := model.NewModel()

	slice := strings.Split(c.Request.URL.Path, "/")

	r.DeleteMenu(slice[len(slice)-1])
	c.JSON(200, gin.H{"respoense": "the menu is deleted", "menuName": slice[len(slice)-1]})
	return
}

func (p *Controller) GetScoreByMenuName(c *gin.Context) {
	r, _ := model.NewModel()

	slice := strings.Split(c.Request.URL.Path, "/")
	c.JSON(200, r.GetScoreByMenuName(slice[len(slice)-1]))
	return
}

func (p *Controller) GetReviewByMenuName(c *gin.Context) {
	r, _ := model.NewModel()

	slice := strings.Split(c.Request.URL.Path, "/")
	c.JSON(200, r.GetReviewByMenuName(slice[len(slice)-1]))
	return
}

func (p *Controller) PostReview(c *gin.Context) {
	r, _ := model.NewModel()

	var newReview model.Review
	if err := c.BindJSON(&newReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	r.PostReview(newReview)
	c.JSON(200, gin.H{"respoense": newReview})
	return
}

func (p *Controller) GetOrder(c *gin.Context) {
	r, _ := model.NewModel()

	slice := strings.Split(c.Request.URL.Path, "/")
	c.JSON(200, r.GetOrderByOrderer(slice[len(slice)-1]))
	return
}

func (p *Controller) PostOrder(c *gin.Context) {
	r, _ := model.NewModel()

	var newOrder model.Order
	if err := c.BindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	r.PostOrder(newOrder)
	c.JSON(200, gin.H{"respoense": newOrder})
	return
}
