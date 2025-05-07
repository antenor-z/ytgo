package downloader

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Download(videoId string, format string) error {
	delOlder()
	cmd := exec.Command(
		"yt-dlp",
		"-P", "public",
		"-f", format+"+bestaudio",
		"https://youtube.com/watch?v="+videoId,
	)
	err := cmd.Start()
	if err != nil {
		return errors.New("yt-dlp failed on download")
	}
	return nil
}

type yt_formats struct {
	Id          string `json:"id"`
	Format_name string `json:"formatName"`
	Resolution  string `json:"resolution"`
	Fps         string `json:"fps"`
	Size        string `json:"size"`
}

func GetFormats(videoId string) ([]yt_formats, error) {
	cmd := exec.Command("yt-dlp", "-F", "https://youtube.com/watch?v="+videoId)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(\d+) *(\w+) *([\dx]+) *(\d+) *[│| ]*([\d.GMKiB]+) *(\w+) *https`)
	lines := strings.Split(string(output), "\n")
	var formats []yt_formats

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
				formats = append([]yt_formats{{
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
