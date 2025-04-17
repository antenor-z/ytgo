package main

import (
	"fmt"
	"os"
	"regexp"
	"time"
	"ytgo/downloader"
	"ytgo/noteConfig"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	noteConfig.ConfigInit()

	if noteConfig.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.MaxMultipartMemory = 256 << 20 // 256MB file max
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{noteConfig.GetDomain()},
		AllowMethods:     []string{"PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api/formats", getFormats)
	r.GET("/api/download/request", requestDownload)
	r.GET("/api/download", Download)

	r.Run(":5100")
}

func getFormats(c *gin.Context) {
	re := regexp.MustCompile(`[A-Za-z0-9_\-]{11}`)
	video_id := re.FindString(c.Query("v"))
	formats, err := downloader.GetFormats(video_id)
	if err != nil || video_id == "" {
		c.JSON(400, gin.H{"error": "invalid video id"})
		return
	}
	fmt.Print(formats)
	c.JSON(200, formats)
}

func requestDownload(c *gin.Context) {
	re := regexp.MustCompile(`[A-Za-z0-9_\-]{11}`)
	video_id := re.FindString(c.Query("v"))
	format := re.FindString(c.Query("f"))

	err := downloader.Download(video_id, format)
	if err != nil || video_id == "" {
		c.JSON(400, gin.H{"error": "invalid video id or format"})
		return
	}
	c.JSON(200, gin.H{"video_id": video_id})
}

func Download(c *gin.Context) {
	// re := regexp.MustCompile(`[A-Za-z0-9_\-]{11}`)
	// video_id := re.FindString(c.Query("v"))
	entries, err := os.ReadDir("./public")
	if err != nil {
		panic(err)
	}

	for _, e := range entries {
		fmt.Println(e.Name())
	}
}
