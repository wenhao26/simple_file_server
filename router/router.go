package router

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"simple_file_server/config"
)

func InitRouter(e *gin.Engine) {
	e.LoadHTMLGlob("templates/*.html")
	e.Static("/assets", "./assets")
	e.Static("/files", "./files")

	e.GET("/", func(c *gin.Context) {
		files, err := filepath.Glob(filepath.Join(config.FileStorage, "*"))
		if err != nil {
			panic(err)
		}

		var links []string
		for _, file := range files {
			//links = append(links, strings.Replace(strings.TrimPrefix(file, config.FileStorage+"/"), "\\", "/", -1))
			links = append(links, strings.Replace(file, "\\", "/", -1))
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "一个简单的文件服务应用 - 首页",
			"files": links,
		})
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

	e.GET("/download", func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			c.String(http.StatusOK, "无下载资源~")
			return
		}

		file, err := os.Open(filename)
		if err != nil {
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Writer.WriteHeader(http.StatusOK)
		c.Header("Content-Disposition", "attachment; filename="+fi.Name())
		c.Header("Content-Type", "application/octet-stream")
		_, _ = io.Copy(c.Writer, file)
	})

	e.POST("/delete", func(c *gin.Context) {
		filename := c.PostForm("filename")

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusNotFound,
				"message": "删除资源不存在",
			})
			return
		}

		err := os.Remove(filename)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "删除失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "删除成功",
		})
	})

}
