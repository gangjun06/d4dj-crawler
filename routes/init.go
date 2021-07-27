package routes

import (
	"strconv"

	"github.com/gangjun06/d4dj-info-server/conf"
	"github.com/gin-gonic/gin"
)

func InitServer() {
	gin.SetMode(gin.ReleaseMode)
	if conf.Get().Debug {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	r.LoadHTMLGlob("templates/**/*.html")

	initViewRoutes(r.Group("/"))
	initFileRoutes(r.Group("/api/file"))

	port := conf.Get().Port
	r.Run(":" + strconv.Itoa(port))
}
