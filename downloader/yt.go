package downloader

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func Download(videoId string, format string) error {
	cmd := exec.Command(
		"yt-dlp",
		"-P", "public",
		"-f", format,
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
}

func GetFormats(videoId string) ([]yt_formats, error) {
	cmd := exec.Command("yt-dlp", "-F", "https://youtube.com/watch?v="+videoId)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(\S+)\s+(\S+)\s+(.*?)\s+(\d*)(?:\s+\S*)?\s+\|`)
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
			if err != nil {
				continue
			}
			if fps_int >= 15 {
				formats = append(formats, yt_formats{
					Id:          id,
					Format_name: ext,
					Resolution:  resolution,
					Fps:         fps,
				})
			}
		}
	}

	return formats, nil
}
