package routes

import (
	c "github.com/gangjun06/d4dj-info-server/controllers/view"
	"github.com/gin-gonic/gin"
)

func initViewRoutes(r *gin.RouterGroup) {
	r.GET("/", c.IndexPage)
	r.GET("/explore", c.ExplorePage)
}
