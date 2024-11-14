package helpers

import "strings"

func ExtractLastPartInPath(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}
