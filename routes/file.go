package routes

import (
	c "github.com/gangjun06/d4dj-info-server/controllers/file"
	"github.com/gin-gonic/gin"
)

func initFileRoutes(r *gin.RouterGroup) {
	r.GET("/list", c.FileList)
	r.GET("/download", c.DownloadFile)
}
