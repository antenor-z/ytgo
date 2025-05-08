package main

import (
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
	"ytgo/config"
	"ytgo/downloader"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	internal.GET("/api/download/file", download)
	internal.GET("/api/download/ready", isVideoReady)
	r.GET("/", getMainPage)
	r.GET("/static/:fileName", getStaticFile)

	r.Run(":5100")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		password := c.Query("p")
		if password != config.GetPassword() {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
		}
	}
}

type yt_formats struct {
	Token   string                    `json:"token"`
	Formats []downloader.YtFormatItem `json:"formats"`
}

func getFormats(c *gin.Context) {
	re := regexp.MustCompile(`[A-Za-z0-9_\-]{11}`)
	video_id := re.FindString(c.Query("v"))
	formats, err := downloader.GetFormats(video_id)
	if err != nil || video_id == "" {
		c.JSON(400, gin.H{"error": "invalid video id"})
		return
	}
	token, err := downloader.CreateDestinationDir(video_id)
	if err != nil || video_id == "" {
		c.JSON(400, gin.H{"error": "invalid video id"})
		return
	}
	c.JSON(200, yt_formats{token, formats})
}

func requestDownload(c *gin.Context) {
	reFormat := regexp.MustCompile(`[0-9]{3}`)
	token := c.Query("t")
	format := reFormat.FindString(c.Query("f"))

	err := downloader.Download(format, token)
	if err != nil || token == "" {
		c.JSON(400, gin.H{"error": "invalid video id or format"})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func download(c *gin.Context) {
	token := c.Query("t")

	dirPath := path.Join("./public", token)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read download directory"})
		return
	}

	if len(entries) != 0 {
		c.FileAttachment(path.Join(dirPath, entries[0].Name()), entries[0].Name())
		return
	}

	c.JSON(404, gin.H{"error": "file not found for video id"})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func isVideoReady(c *gin.Context) {
	token := c.Query("t")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to ws"})
		return
	}

	defer conn.Close()
	for {
		time.Sleep(time.Second * 2)
		entries, err := os.ReadDir(path.Join("./public", token))
		if err != nil {
			conn.WriteJSON(gin.H{"error": err.Error()})
			return
		}

		if len(entries) == 1 {
			if !entries[0].IsDir() && videoReady(entries[0].Name()) {
				conn.WriteJSON(gin.H{"status": "ok"})
				return
			}
		}

		conn.WriteJSON(gin.H{"status": "waiting"})
	}
}

func videoReady(filename string) bool {
	return !strings.Contains(filename, ".part") &&
		!strings.Contains(filename, ".ytdl") &&
		!strings.Contains(filename, ".txt")
}

func getMainPage(c *gin.Context) {
	c.File("front/index.html")
}

func getStaticFile(c *gin.Context) {
	path := "front/static/" + c.Param("fileName")
	_, err := os.Stat(path)
	if err != nil {
		c.Abort()
		return
	}

	c.File(path)
}
