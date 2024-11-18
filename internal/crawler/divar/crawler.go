package divar

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/helpers"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
)

type DivarCrawler struct {
	Repository repositories.CrawlerRepository
	Logger     *logger.AppLogger
}

func (c DivarCrawler) GetLogger() *logger.AppLogger {
	return c.Logger
}

func GetSourceName() string {
	return "divar"
}

func (c DivarCrawler) GetSourceName() string {
	return "divar"
}

func (c DivarCrawler) GetBaseUrl() string {
	return "https://divar.ir"
}

func (c DivarCrawler) GetRepository() repositories.CrawlerRepository {
	return c.Repository
}

func (c DivarCrawler) CreateCrawlJobArchivePageLink(link string) repositories.RepoResult {
	return c.Repository.CreateCrawlJobArchivePageLink(link, GetSourceName())
}

func (c DivarCrawler) GetSinglePageLinksFromArchivePage(htmlContent string) ([]string, error) {
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

func (c DivarCrawler) CrawlItemPage(job sqlc.CrawlJob) (structs.CrawledData, error) {
	c.GetLogger().Infof(" | crawl single page with %s crawler | JobID: %d", c.GetSourceName(), job.ID)

	htmlContent, err := helpers.GetHtml(job.Url)
	if err != nil {
		return structs.CrawledData{}, err
	}

	c.GetLogger().Infof(" | extracted html | JobID: %d", job.ID)

	crawledData := structs.CrawledData{
		SourceName: GetSourceName(),
	}

	var errors []error

	// fill general data
	err = c.catchGeneralData(htmlContent, &crawledData)
	if err != nil {
		errors = append(errors, err)
	}

	c.GetLogger().Infof(" | extracted general properties | JobID: %d", job.ID)

	err = c.catchPublishedAt(htmlContent, &crawledData)
	if err != nil {
		errors = append(errors, err)
	}

	c.GetLogger().Infof(" | extracted published at property | JobID: %d", job.ID)

	err = c.catchPricesAndSomeOtherData(htmlContent, &crawledData)
	if err != nil {
		errors = append(errors, err)
	}

	c.GetLogger().Infof(" | extracted price values and some other properties | JobID: %d", job.ID)

	if len(errors) > 0 {
		_, err = c.Repository.UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
		if err != nil {
			fmt.Println(err)
		}
		return structs.CrawledData{}, errors[0]
	}

	c.GetLogger().Infof(" | crawling single page finished | JobID: %d", job.ID)

	return crawledData, nil
}

func (c DivarCrawler) catchGeneralData(htmlContent string, crawledData *structs.CrawledData) error {
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
	crawledData.PublisherAdKey = helpers.ExtractLastPartInPath(items[0].URL)
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

func (c DivarCrawler) catchPublishedAt(htmlContent string, crawledData *structs.CrawledData) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return fmt.Errorf("error parsing html: %s", err)
	}

	publishedAt := ""
	doc.Find("title").Each(func(index int, item *goquery.Selection) {
		titleValue := item.Text()
		parts := strings.Split(titleValue, "-")
		publishedAt = parts[len(parts)-1]
	})

	if publishedAt == "" {
		return fmt.Errorf("PublishedAt value not found")
	}

	crawledData.PublishedAt = helpers.HumanDateToNormalDate(helpers.ToEnglishDigits(publishedAt))
	return nil
}

func (c DivarCrawler) catchPricesAndSomeOtherData(htmlContent string, crawledData *structs.CrawledData) error {
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
		crawledData.Year = strings.TrimSpace(strings.ReplaceAll(crawledData.Year, "قبل از", ""))
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
