package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"ytgo/config"
	"ytgo/downloader"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConfigInit()

	if config.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(cors.Default())

	internal := r.Group("/")
	internal.Use(AuthMiddleware())

	internal.GET("/api/formats", getFormats)
	internal.GET("/api/download/request", requestDownload)
	internal.GET("/api/download", download)
	internal.GET("/api/download/ready", isVideoReady)
	r.GET("/", getMainPage)

	r.Run(":5100")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		password := c.Query("p")
		if password != config.GetPassword() {
			c.JSON(401, "unauthorized")
			c.Abort()
		}
	}
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
	format := c.Query("f")

	err := downloader.Download(video_id, format)
	if err != nil || video_id == "" {
		c.JSON(400, gin.H{"error": "invalid video id or format"})
		return
	}
	c.JSON(200, gin.H{"video_id": video_id})
}

func download(c *gin.Context) {
	videoID := c.Query("v")
	if videoID == "" {
		c.JSON(400, gin.H{"error": "missing video id"})
		return
	}

	entries, err := os.ReadDir("./public")
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read download directory"})
		return
	}

	// Look for a file that includes the video ID
	for _, entry := range entries {
		if !entry.IsDir() && containsVideoID(entry.Name(), videoID) {
			filePath := "./public/" + entry.Name()
			// Serve it as a file download
			c.FileAttachment(filePath, entry.Name())
			return
		}
	}

	c.JSON(404, gin.H{"error": "file not found for video id"})
}

func isVideoReady(c *gin.Context) {
	videoID := c.Query("v")
	if videoID == "" {
		c.JSON(400, gin.H{"error": "missing video id"})
		return
	}

	entries, err := os.ReadDir("./public")
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read download directory"})
		return
	}

	// Look for a file that includes the video ID
	for _, entry := range entries {
		if !entry.IsDir() && containsVideoID(entry.Name(), videoID) {
			c.JSON(200, gin.H{"status": "ok"})
			return
		}
	}

	c.JSON(404, gin.H{"error": "file not found for video id"})
}

func containsVideoID(filename, videoID string) bool {
	return strings.Contains(filename, videoID) &&
		!strings.Contains(filename, ".part") &&
		!strings.Contains(filename, ".ytdl")
}

func getMainPage(c *gin.Context) {
	c.File("front/index.html")
}
