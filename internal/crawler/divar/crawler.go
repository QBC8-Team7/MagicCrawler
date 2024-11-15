package divar

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/helpers"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
)

const ARCHIVE_PAGE = "archive"
const SINGLE_PAGE = "single"

type Crawler struct {
	Repository repositories.CrawlJobRepository
}

func GetSourceName() string {
	return "divar"
}

func (c Crawler) GetBaseUrl() string {
	return "https://divar.ir"
}

func (c Crawler) CreateCrawlJobArchivePageLink(link string) repositories.RepoResult {
	return c.Repository.CreateCrawlJobArchivePageLink(link, GetSourceName())
}

func (c Crawler) CrawlArchivePage(job sqlc.CrawlJob, wg *sync.WaitGroup, timeoutCh <-chan time.Time) {
	defer wg.Done()

	_, err := c.Repository.UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_PICKED)
	if err != nil {
		fmt.Println(err)
		return
	}

	htmlContent, err := helpers.GetHtml(job.Url)
	if err != nil {
		// TODO - log here
		fmt.Println(err)
		// TODO - maybe we need to put error in db
		// TODO - maybe we need to save resource usage and time
		// TODO - maybe we can add try fields for job
		c.Repository.UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
		return
	}

	links, err := getSinglePageLinksFromArchivePage(htmlContent)
	if err != nil {
		fmt.Println(err)
		// TODO - maybe we need to put error in db
		// TODO - maybe we need to save resource usage and time
		// TODO - maybe we can add try fields for job
		c.Repository.UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
		return
	}

	if len(links) > 0 {
		fmt.Println("links count:", len(links))
		nextLink, err := helpers.GetNextPageLink(job.Url)
		if err != nil {
			fmt.Println(err)
			// TODO - error handling
			// TODO - maybe we need to save resource usage and time
			// TODO - maybe we can add try fields for job
			c.Repository.UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
			return
		}

		// TODO - maybe we need to use transactions to make sure all links with next link inserted successfuly together

		errors := c.Repository.CreateCrawlJobForSinglePageLinks(links, GetSourceName())
		if len(errors) > 0 {
			fmt.Println(errors[0])
			// TODO - error handling
			// TODO - maybe we need to save resource usage and time
			// TODO - maybe we can add try fields for job
			c.Repository.UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
			return
		}

		nextLinkResult := c.Repository.CreateCrawlJobArchivePageLink(nextLink, GetSourceName())
		fmt.Println("next link:", nextLink)
		if nextLinkResult.Err != nil {
			// TODO - log here
			fmt.Println(nextLinkResult.Err)
			// TODO - maybe we need to put error in db
			// TODO - maybe we need to save resource usage and time
			// TODO - maybe we can add try fields for job
			c.Repository.UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
			return
		}

		wg.Add(len(links) + 1)
	}

	c.Repository.UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_DONE)
}

func getSinglePageLinksFromArchivePage(htmlContent string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return []string{}, fmt.Errorf("error parsing html: %s", err)
	}

	var scriptContent string
	doc.Find("[type='application/ld+json']").Each(func(index int, item *goquery.Selection) {
		scriptContent = item.Text()
	})

	if len(scriptContent) == 0 {
		return []string{}, fmt.Errorf("no json-ld script found")
	}

	var items []ArchivePageItem
	err = json.Unmarshal([]byte(scriptContent), &items)
	if err != nil {
		return []string{}, fmt.Errorf("error unmarshalling json: %s", err)
	}

	links := make([]string, len(items))
	for index, item := range items {
		links[index] = item.URL
	}

	return links, nil
}

func (c Crawler) CrawlItemPage(job sqlc.CrawlJob, wg *sync.WaitGroup, timeoutCh <-chan time.Time) (structs.CrawledData, error) {
	defer wg.Done()
	fmt.Println("crawl item page", job.ID)
	c.Repository.UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_DONE)
	return structs.CrawledData{}, nil

	htmlContent, err := helpers.GetHtml(job.Url)
	if err != nil {
		return structs.CrawledData{}, err
	}

	crawledData := structs.CrawledData{}

	// fill general data
	err = c.catchGeneralData(htmlContent, &crawledData)
	if err != nil {
		return structs.CrawledData{}, err
	}

	err = c.catchPublishedAt(htmlContent, &crawledData)
	if err != nil {
		return structs.CrawledData{}, err
	}

	err = c.catchPricesAndSomeOtherData(htmlContent, &crawledData)
	if err != nil {
		return structs.CrawledData{}, err
	}

	// insert data to db
	return crawledData, nil
}

func (c Crawler) catchGeneralData(htmlContent string, crawledData *structs.CrawledData) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return fmt.Errorf("error parsing html: %s", err)
	}

	data := make(map[string]string)
	doc.Find("[type='application/ld+json']").Each(func(index int, item *goquery.Selection) {
		data["script"] = item.Text()
	})

	jsonString, exists := data["script"]
	if !exists {
		return fmt.Errorf("no json-ld script found")
	}

	var items []GeneralFields
	err = json.Unmarshal([]byte(jsonString), &items)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %s", err)
	}

	crawledData.HouseType = getHouseType(items[0].AccommodationCategory)
	crawledData.AdCategory = getAdType(items[0].AccommodationCategory)
	crawledData.Description = items[0].Description
	crawledData.Meterage = helpers.UnsafeAtoi(items[0].FloorSize.Value)
	crawledData.Lat = helpers.ToEnglishDigits(items[0].Geo.Latitude)
	crawledData.Lon = helpers.ToEnglishDigits(items[0].Geo.Longitude)
	crawledData.ImageUrl = items[0].Image
	crawledData.Title = items[0].Name
	crawledData.RoomsCount = helpers.WordNumberToNumber(items[0].NumberOfRooms)
	crawledData.URL = items[0].URL
	crawledData.AdId = helpers.ExtractLastPartInPath(items[0].URL)
	crawledData.City = helpers.ArabicToPersianChars(items[0].WebInfo.CityPersian)
	crawledData.Neighborhood = helpers.ArabicToPersianChars(items[0].WebInfo.DistrictPersian)
	if items[0].WebInfo.DistrictPersian == "" {
		crawledData.Neighborhood = helpers.ArabicToPersianChars(items[0].WebInfo.CityPersian)
	}

	return nil
}

func getAdType(category string) string {
	if strings.Contains(category, "اجاره") {
		return "rent"
	}

	if strings.Contains(category, "فروش") {
		return "buy"
	}

	if strings.Contains(category, "رهن") {
		return "mortgage"
	}

	return ""
}

func getHouseType(category string) string {
	if strings.Contains(category, "آپارتمان") {
		return "apartment"
	}

	if strings.Contains(category, "خانه") || strings.Contains(category, "ویلا") {
		return "villa"
	}

	return ""
}

func (c Crawler) catchPublishedAt(htmlContent string, crawledData *structs.CrawledData) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return fmt.Errorf("error parsing html: %s", err)
	}

	publishedAt := ""
	doc.Find("title").Each(func(index int, item *goquery.Selection) {
		titleValue := item.Text()
		parts := strings.Split(titleValue, "-")
		publishedAt = parts[1]
	})

	if publishedAt == "" {
		return fmt.Errorf("PublishedAt value not found")
	}

	crawledData.PublishedAt = helpers.HumanDateToNormalDate(helpers.ToEnglishDigits(publishedAt))
	return nil
}

func (c Crawler) catchPricesAndSomeOtherData(htmlContent string, crawledData *structs.CrawledData) error {
	startPattern := `"LIST_DATA"\s*:\s*`
	endPattern := `\s*}\s*]\s*}\s*`

	slicedString := helpers.SubStringBetweenTwoRegEx(htmlContent, startPattern, endPattern)
	slicedString = helpers.RemoveLastCurlyBrace(slicedString)

	var widgets []Widget
	err := json.Unmarshal([]byte(slicedString), &widgets)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %v", err)
	}

	var results = make(map[string]interface{})

	for _, widget := range widgets {
		if widget.Items != nil {
			for _, item := range widget.Items {
				if item.IconName == "elevator" {
					if !item.Disabled {
						results["has_elevator"] = true
					} else if item.Disabled {
						results["has_elevator"] = false
					}
				} else if item.IconName == "cabinet" {
					if !item.Disabled {
						results["has_warehouse"] = true
					} else if item.Disabled {
						results["has_warehouse"] = false
					}
				} else if item.Title == "ساخت" {
					results["year"] = item.Value
				}
			}
		} else {
			if widget.Title == "ودیعه" {
				results["prepaid_price"] = widget.Value
			} else if widget.Title == "اجارهٔ ماهانه" {
				results["monthly_rent_price"] = widget.Value
			} else if widget.Title == "طبقه" {
				results["floor_number"] = widget.Value
			} else if widget.Title == "قیمت کل" {
				results["total_price"] = widget.Value
			} else if widget.Title == "قیمت هر متر" {
				results["price_per_meter"] = widget.Value
			}
		}

		// IN DIVAR THESE VALUES WILL BE SHOWN IN TWO STYLES
		// WE HANDLE THEM IN TWO WAYS
		// FIRST UPPER ONE AND SECOND BELOW ONE
		if widget.Credit != nil {
			results["total_price"] = strconv.Itoa(widget.Credit.Value)
		}
		if widget.Rent != nil {
			results["monthly_rent_price"] = strconv.Itoa(widget.Rent.Value)
		}

	}

	if results["total_price"] != nil {
		crawledData.TotalPrice = helpers.CleanPrice(results["total_price"].(string))

	} else {
		crawledData.TotalPrice = ""
	}

	if results["monthly_rent_price"] != nil {
		crawledData.MonthlyRentPrice = helpers.CleanPrice(results["monthly_rent_price"].(string))
	} else {
		crawledData.MonthlyRentPrice = ""
	}

	if results["price_per_meter"] != nil {
		crawledData.PricePerMeter = helpers.CleanPrice(results["price_per_meter"].(string))
	} else {
		crawledData.PricePerMeter = ""
	}

	if results["prepaid_price"] != nil {
		crawledData.PrePaidPrice = helpers.CleanPrice(results["prepaid_price"].(string))
	} else {
		crawledData.PrePaidPrice = ""
	}

	_, exist := results["has_warehouse"]
	if !exist {
		crawledData.HasWarehouse = false
	} else {
		crawledData.HasWarehouse = results["has_warehouse"].(bool)
	}

	_, exist = results["has_elevator"]
	if !exist {
		crawledData.HasElevator = false
	} else {
		crawledData.HasElevator = results["has_elevator"].(bool)
	}

	_, exist = results["year"]
	if !exist {
		crawledData.Year = ""
		crawledData.Age = 0
	} else {
		crawledData.Year = helpers.ToEnglishDigits(results["year"].(string))
		crawledData.Age = helpers.YearToAge(crawledData.Year)
	}

	_, exist = results["floor_number"]
	if !exist {
		crawledData.FloorNumber = 0
	} else {
		crawledData.FloorNumber = helpers.UnsafeAtoi(helpers.ToEnglishDigits(helpers.GetFirstValueOfAPersianRange(results["floor_number"].(string))))
	}

	return nil
}
