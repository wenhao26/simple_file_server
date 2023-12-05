package router

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"simple_file_server/config"
)

func InitRouter(e *gin.Engine) {
	e.LoadHTMLGlob("templates/*.html")
	e.Static("/assets", "./assets")

	e.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "一个简单的文件服务应用 - 首页"})
	})

	e.POST("/upload", func(c *gin.Context) {
		file, handler, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": fmt.Sprintf("上传图片出错:", err),
			})
		}

		if _, err := os.Stat(config.FileStorage); os.IsNotExist(err) {
			_ = os.MkdirAll(config.FileStorage, os.ModePerm)
		}

		f, err := os.OpenFile(config.FileStorage+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer f.Close()

		_, _ = io.Copy(f, file)

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "上传成功",
			"data": map[string]interface{}{
				"path":     config.FileStorage,
				"filename": handler.Filename,
			},
		})
	})

}
