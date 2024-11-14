package main

import (
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler"
)

func main() {
	seeds := []map[string]string{
		{"link": "https://divar.ir/s/tehran-province/buy-apartment", "source": crawler.SOURCE_DIVAR},
		// {"link": "https://divar.ir/s/tehran-province/buy-villa", "source": crawler.SOURCE_DIVAR},
		// {"link": "https://divar.ir/s/tehran-province/rent-apartment", "source": crawler.SOURCE_DIVAR},
		// {"link": "https://divar.ir/s/tehran-province/rent-villa", "source": crawler.SOURCE_DIVAR},
	}

	// Set crawl duration to 10 minutes
	timeout := time.Duration(10) * time.Minute
	crawler.Start(seeds, timeout)
}
