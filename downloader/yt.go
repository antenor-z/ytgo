package downloader

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Download(format string, token string) error {
	delOlder()
	fileBytes, err := os.ReadFile(filepath.Join("./public", token, "url.txt"))
	if err != nil {
		return err
	}
	os.Remove(filepath.Join("./public", token, "url.txt"))
	cmd := exec.Command(
		"yt-dlp",
		"-P", path.Join("./public", token),
		"-f", format+"+bestaudio",
		string(fileBytes),
	)
	err = cmd.Start()
	if err != nil {
		return errors.New("yt-dlp failed on download")
	}
	return nil
}

type YtFormatItem struct {
	Id          string `json:"id"`
	Format_name string `json:"formatName"`
	Resolution  string `json:"resolution"`
	Fps         string `json:"fps"`
	Size        string `json:"size"`
}

func CreateDestinationDir(videoId string) (string, error) {
	url := "https://youtube.com/watch?v=" + videoId
	dirName := tokenGenerator()
	err := os.MkdirAll(filepath.Join("./public", dirName), os.ModePerm)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(filepath.Join("./public", dirName, "url.txt"), []byte(url), os.ModePerm)
	if err != nil {
		return "", err
	}

	return dirName, nil
}

func GetFormats(videoId string) ([]YtFormatItem, error) {
	url := "https://youtube.com/watch?v=" + videoId
	cmd := exec.Command("yt-dlp", "-F", url)
	output, err := cmd.Output()
	if err != nil {
		return []YtFormatItem{}, err
	}

	re := regexp.MustCompile(`^(\d+) *(\w+) *([\dx]+) *(\d+) *[│| ]*([\d.GMKiB]+) *(\w+) *https`)
	lines := strings.Split(string(output), "\n")
	var formats []YtFormatItem

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			id := matches[1]
			ext := matches[2]
			resolution := strings.TrimSpace(matches[3])
			fps := matches[4]
			fps_int, err := strconv.Atoi(fps)
			size := matches[5]
			// vbr := matches[6]
			if err != nil {
				continue
			}
			if fps_int >= 15 {
				formats = append([]YtFormatItem{{
					Id:          id,
					Format_name: ext,
					Resolution:  resolution,
					Fps:         fps,
					Size:        size,
				}}, formats...)
			}
		}
	}

	return formats, nil
}

func tokenGenerator() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
func delOlder() {
	entries, err := os.ReadDir("./public")
	if err != nil {
		return
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if time.Since(info.ModTime()) > 15*time.Minute {
			os.Remove(entry.Name())
		}
	}
}
