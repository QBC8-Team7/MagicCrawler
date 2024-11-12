package helpers

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetHtml(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout after 10 seconds
	}

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Googlebot/2.1 (+http://www.google.com/bot.html)")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching the URL: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK HTTP status: %s", res.Status)
	}

	htmlContent, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(htmlContent), nil
}
