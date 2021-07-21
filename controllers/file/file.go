package file

import (
	"net/http"
	"path/filepath"

	"github.com/gangjun06/d4dj-info-server/service/file"
	"github.com/gin-gonic/gin"
)

func FileList(c *gin.Context) {

	targetPath := c.Query("path")
	s := file.NewExplorer(targetPath)

	if s.NotFound || !s.IsDir {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": s.FileList(),
	})
}

func DownloadFile(c *gin.Context) {
	targetPath := c.Query("path")
	s := file.NewExplorer(targetPath)
	if s.NotFound || s.IsDir {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(s.Path))
	c.Header("Content-Type", "application/octet-stream")
	c.File(s.Path)
}
