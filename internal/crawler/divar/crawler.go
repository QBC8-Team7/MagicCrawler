package divar

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/helpers"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
)

type Crawler struct{}

func (c Crawler) CrawlArchivePage(link string) []string {
	// TODO - Implement crawl archive page

	fmt.Println("ArchivePage")
	return []string{
		"https://divar.ir/v/%D8%A8%D8%B1%D8%AC-%D9%85%D9%8F%D8%AC%D9%8E%D9%84%D9%84-%D9%88%D9%8B-%D9%87%D9%8F%D8%AA%D9%90%D9%84%DB%8C%D9%86%DA%AF-110%D9%85%D8%AA%D8%B12%D8%AE-2%D8%AA%D8%B1%D8%A7%D8%B3-2%D9%BE%D8%A7%D8%B1%DA%A9%DB%8C%D9%86%DA%AF/wZhYfCej",
		// "https://divar.ir/v/%DB%B9%DB%B0-%D9%85%D8%AA%D8%B1-%D8%AF%D9%88-%D8%AE%D9%88%D8%A7%D8%A8%D9%87-%D9%81%D9%88%D9%84-%D8%A7%D9%85%DA%A9%D8%A7%D9%86%D8%A7%D8%AA-%D8%B4%D9%87%D8%B1%DA%A9-%D9%85%D8%A8%D8%B9%D8%AB-%D8%AC%D9%86%D8%AA-%D8%B4%D9%85%D8%A7%D9%84/wZnkiwjN",
		// "https://divar.ir/v/%D9%81%D8%B1%D9%88%D8%B4-%D9%88%DB%8C%D9%84%D8%A7-%D8%A8%D8%A7%D8%BA-%D8%AF%D9%85%D8%A7%D9%88%D9%86%D8%AF-%DB%B3%DB%B5%DB%B0-%D9%85%D8%AA%D8%B1/wZpwsRu2",
		// "https://divar.ir/v/%D9%81%D8%B1%D9%88%D8%B4%DB%B2%DB%B6%DB%B0%D9%85%D8%AA%D8%B1-%D9%88%DB%8C%D9%84%D8%A7%DB%8C%DB%8C-%D8%AF%D9%88%D8%B7%D8%A8%D9%82%D9%87-%D9%85%D8%AC%D8%B2%D8%A7-%D8%AF%D8%B1-%D9%81%D8%A7%D8%B2%DB%B3-%D8%A7%D9%86%D8%AF%DB%8C%D8%B4%D9%87/wZpo8rRP",
		// "https://divar.ir/v/%D9%81%D8%B1%D9%88%D8%B4-%D8%A2%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86/wZlERUNi",
		// "https://divar.ir/v/%D8%A2%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86-80-%D9%85%D8%AA%D8%B1%DB%8C-%DB%8C%D8%A7%D8%AE%DA%86%DB%8C-%D8%A2%D8%A8%D8%A7%D8%AF-%D9%86%D8%A7%D8%B2%DB%8C-%D8%A7%D8%A8%D8%A7%D8%AF/wZqcVWMQ",
		// "https://divar.ir/v/%D8%AE%D8%A7%D9%86%D9%87-%D9%88%DB%8C%D9%84%D8%A7%DB%8C%DB%8C-%D8%B3%D9%86%D8%AF-%D8%AF%D8%A7%D8%B1/wZmMb4lq",
		// "https://divar.ir/v/%D8%A2%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86-%DB%B1%DB%B2%DB%B0%D9%85%D8%AA%D8%B1%DB%8C-%D8%B5%D8%A7%D9%84%D8%AD%DB%8C%D9%87/wZJcUdyR",
		// "https://divar.ir/v/%D8%A2%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86-%D8%B3%D9%87-%D8%AE%D9%88%D8%A7%D8%A8%D9%87-%D8%B4%D9%85%D8%A7%D9%84-%D8%AC%D8%B1%D8%AF%D9%86/wZpc0SVj",
		// "https://divar.ir/v/150%D9%85%D8%AA%D8%B1-%D8%A8%D8%B1%D8%AC-%D8%A8%D8%A7%D8%BA-%D9%86%D9%88%D8%B3%D8%A7%D8%B2-%D9%88%DB%8C%D9%88-%D8%A7%D8%A8%D8%AF%DB%8C-%D8%B3%D8%B9%D8%A7%D8%AF%D8%AA-%D8%A2%D8%A8%D8%A7%D8%AF/wZmM0v34",
	}

}

func (c Crawler) CrawlItemPage(link string) (structs.CrawledData, error) {
	htmlContent, err := helpers.GetHtml(link)
	if err != nil {
		return structs.CrawledData{}, err
	}

	crawledData := structs.CrawledData{}

	// fill general data
	err = c.CatchGeneralData(htmlContent, &crawledData)
	if err != nil {
		return structs.CrawledData{}, err
	}

	err = c.CatchPublishedAt(htmlContent, &crawledData)
	if err != nil {
		return structs.CrawledData{}, err
	}

	err = c.CatchPricesAndSomeOtherData(htmlContent, &crawledData)
	if err != nil {
		return structs.CrawledData{}, err
	}

	return crawledData, nil
}

func (c Crawler) CatchGeneralData(htmlContent string, crawledData *structs.CrawledData) error {
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

func (c Crawler) CatchPublishedAt(htmlContent string, crawledData *structs.CrawledData) error {
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

func (c Crawler) CatchPricesAndSomeOtherData(htmlContent string, crawledData *structs.CrawledData) error {
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
