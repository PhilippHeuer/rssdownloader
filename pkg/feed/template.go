package feed

import (
	"runtime"
	"strings"

	"github.com/mmcdole/gofeed"
)

func RenderTemplate(template string, item *gofeed.Item) string {
	if template == "" {
		template = "{title}"
	}

	template = strings.ReplaceAll(template, "{title}", item.Title)
	if runtime.GOOS == "windows" {
		template = sanitizeFileName(template)
	}
	return template
}

func sanitizeFileName(filename string) string {
	invalidChars := []string{":", "\\", "/", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, char, "-")
	}
	return filename
}
