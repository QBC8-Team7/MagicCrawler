package helpers

import (
	"net/url"
	"strconv"
	"strings"
)

func ExtractLastPartInPath(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func GetNextPageLink(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	query := u.Query()
	pageStr := query.Get("page")

	page := 2
	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil || pageInt <= 1 {
			page = 2
		} else {
			page = pageInt + 1
		}
	}

	query.Set("page", strconv.Itoa(page))
	u.RawQuery = query.Encode()
	return u.String(), nil
}
