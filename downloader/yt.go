package downloader

import (
	"os/exec"
	"regexp"
	"strings"
)

func Download(videoId string, format string) error {
	err := exec.Command("yt-dlp", "-P", "public", "-f", format, "https://youtube.com/watch?v="+videoId).Start()
	if err != nil {
		return err
	}
	return nil
}

type yt_formats struct {
	id          string
	format_name string
	resolution  string
	fps         string
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
			formats = append(formats, yt_formats{
				id:          id,
				format_name: ext,
				resolution:  resolution,
				fps:         fps,
			})
		}
	}

	return formats, nil
}
