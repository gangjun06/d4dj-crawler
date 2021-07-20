package routes

import (
	"github.com/gangjun06/d4dj-info-server/env"
	"github.com/gin-gonic/gin"
)

func InitServer() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/**/*.html")

	initViewRoutes(r.Group("/"))

	port := env.Get(env.KeyServerPort)
	r.Run(":" + port)
}
