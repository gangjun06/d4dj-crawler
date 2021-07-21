package views

import (
	"net/http"
	"path/filepath"

	"github.com/gangjun06/d4dj-info-server/service/file"
	"github.com/gin-gonic/gin"
)

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "views/index.html", gin.H{})
}

func ExplorePage(c *gin.Context) {
	targetPath := c.Query("path")
	s := file.NewExplorer(targetPath)
	if s.NotFound {
		c.HTML(http.StatusNotFound, "views/404.html", gin.H{})
		return
	}

	if !s.IsDir {
		ext := filepath.Ext(s.Path)
		data, err := s.FileData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		dir, file := filepath.Split(targetPath)
		if dir == "" {
			dir = "Root/"
		}
		c.HTML(http.StatusNotFound, "views/preview.html", gin.H{
			"type":     ext,
			"data":     string(data),
			"name":     file,
			"dir":      dir,
			"src":      "/api/file/download?path=" + targetPath,
			"download": "/api/file/download?path=" + targetPath,
		})
		return
	}

	list := s.FileList()

	prevPage := filepath.Dir(targetPath)
	if prevPage == "." {
		prevPage = "/explore"
	} else {
		prevPage = "/explore?path=" + prevPage
	}
	formatedPath := targetPath
	if targetPath == "" {
		formatedPath = "Root"
		targetPath = "."
	}
	c.HTML(http.StatusOK, "views/explore.html", gin.H{
		"list":         list,
		"path":         targetPath,
		"formatedPath": formatedPath,
		"prevPage":     prevPage,
	})
}
