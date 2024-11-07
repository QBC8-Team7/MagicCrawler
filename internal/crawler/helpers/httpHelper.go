package helpers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func GetHtml(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout after 10 seconds
	}

	res, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching the URL: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK HTTP status: %s", res.Status)
	}

	htmlContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(htmlContent), nil
}
