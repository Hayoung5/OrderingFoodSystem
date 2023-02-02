package router

import (
	"fmt"
	ctl "lecture/go-project/controller"

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

// 임의 인증을 위한 함수
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

// 실제 라우팅
func (p *Router) Index() *gin.Engine {
	e := gin.Default()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Use(CORS())

	account := e.Group("acc/v01", liteAuth())
	{
		account.GET("/ok", p.ct.GetOK)
	}

	menu := e.Group("menu", liteAuth())
	{
		menu.GET("/getMenu", p.ct.GetMenu)
		menu.GET("/getScore", p.ct.GetScoreByMenuName)
		menu.POST("/postMenu", p.ct.PostMenu)
		menu.PATCH("/updateMenu", p.ct.UpdateMenu)
		menu.DELETE("/updateMenu/:menuName", p.ct.DeleteMenu)
	}

	review := e.Group("review", liteAuth())
	{
		review.GET("/getReview/:menuName", p.ct.GetReviewByMenuName)
		review.POST("/postReview", p.ct.PostReview)
	}

	order := e.Group("order", liteAuth())
	{
		order.GET("/getOrder/:orderer", p.ct.GetOrder)
		order.POST("/postOrder", p.ct.PostOrder)
		// order.PATCH()
		// order.PATCH()
	}

	return e
}
