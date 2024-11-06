package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// URL of the webpage to crawl
	// url := "https://divar.ir/v/%D8%A2%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86-%DB%B1%DB%B2%DB%B0%D9%85%D8%AA%D8%B1%DB%8C-%D8%B5%D8%A7%D9%84%D8%AD%DB%8C%D9%87/wZJcUdyR"
	// url := "https://divar.ir/v/%D8%A2%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86-%D8%B3%D9%87-%D8%AE%D9%88%D8%A7%D8%A8%D9%87-%D8%B4%D9%85%D8%A7%D9%84-%D8%AC%D8%B1%D8%AF%D9%86/wZpc0SVj"
	// url := "https://divar.ir/v/150%D9%85%D8%AA%D8%B1-%D8%A8%D8%B1%D8%AC-%D8%A8%D8%A7%D8%BA-%D9%86%D9%88%D8%B3%D8%A7%D8%B2-%D9%88%DB%8C%D9%88-%D8%A7%D8%A8%D8%AF%DB%8C-%D8%B3%D8%B9%D8%A7%D8%AF%D8%AA-%D8%A2%D8%A8%D8%A7%D8%AF/wZmM0v34"
	// url := "https://divar.ir/v/%D8%A8%D8%B1%D8%AC-%D9%85%D9%8F%D8%AC%D9%8E%D9%84%D9%84-%D9%88%D9%8B-%D9%87%D9%8F%D8%AA%D9%90%D9%84%DB%8C%D9%86%DA%AF-110%D9%85%D8%AA%D8%B12%D8%AE-2%D8%AA%D8%B1%D8%A7%D8%B3-2%D9%BE%D8%A7%D8%B1%DA%A9%DB%8C%D9%86%DA%AF/wZhYfCej"
	// url := "https://divar.ir/v/%DB%B9%DB%B0-%D9%85%D8%AA%D8%B1-%D8%AF%D9%88-%D8%AE%D9%88%D8%A7%D8%A8%D9%87-%D9%81%D9%88%D9%84-%D8%A7%D9%85%DA%A9%D8%A7%D9%86%D8%A7%D8%AA-%D8%B4%D9%87%D8%B1%DA%A9-%D9%85%D8%A8%D8%B9%D8%AB-%D8%AC%D9%86%D8%AA-%D8%B4%D9%85%D8%A7%D9%84/wZnkiwjN"
	// url := "https://divar.ir/v/%D9%81%D8%B1%D9%88%D8%B4-%D9%88%DB%8C%D9%84%D8%A7-%D8%A8%D8%A7%D8%BA-%D8%AF%D9%85%D8%A7%D9%88%D9%86%D8%AF-%DB%B3%DB%B5%DB%B0-%D9%85%D8%AA%D8%B1/wZpwsRu2"
	// url := "https://divar.ir/v/%D9%81%D8%B1%D9%88%D8%B4%DB%B2%DB%B6%DB%B0%D9%85%D8%AA%D8%B1-%D9%88%DB%8C%D9%84%D8%A7%DB%8C%DB%8C-%D8%AF%D9%88%D8%B7%D8%A8%D9%82%D9%87-%D9%85%D8%AC%D8%B2%D8%A7-%D8%AF%D8%B1-%D9%81%D8%A7%D8%B2%DB%B3-%D8%A7%D9%86%D8%AF%DB%8C%D8%B4%D9%87/wZpo8rRP"
	// url := "https://divar.ir/v/%D9%81%D8%B1%D9%88%D8%B4-%D8%A2%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86/wZlERUNi"
	// url := "https://divar.ir/v/%D8%A2%D9%BE%D8%A7%D8%B1%D8%AA%D9%85%D8%A7%D9%86-80-%D9%85%D8%AA%D8%B1%DB%8C-%DB%8C%D8%A7%D8%AE%DA%86%DB%8C-%D8%A2%D8%A8%D8%A7%D8%AF-%D9%86%D8%A7%D8%B2%DB%8C-%D8%A7%D8%A8%D8%A7%D8%AF/wZqcVWMQ"
	url := "https://divar.ir/v/%D8%AE%D8%A7%D9%86%D9%87-%D9%88%DB%8C%D9%84%D8%A7%DB%8C%DB%8C-%D8%B3%D9%86%D8%AF-%D8%AF%D8%A7%D8%B1/wZmMb4lq"

	ad := crawl(url)

	ad.clean()
	ad.prepareNumberOfRooms()
	ad.computeAgeFromYear(1403)
	ad.prepareTypes()

	fmt.Println("Accommodation Category:", ad.AccommodationCategory)
	fmt.Println("Description:", ad.Description)
	fmt.Println("Floor Size (Value):", ad.FloorSize.Value)
	fmt.Println("Geo (Latitude):", ad.Geo.Latitude)
	fmt.Println("Geo (Longitude):", ad.Geo.Longitude)
	fmt.Println("Image:", ad.Image)
	fmt.Println("Name:", ad.Name)
	fmt.Println("Number of Rooms:", ad.NumberOfRooms)
	fmt.Println("URL:", ad.URL)
	fmt.Println("Web Info (City Persian):", ad.WebInfo.CityPersian)
	fmt.Println("Web Info (District Persian):", ad.WebInfo.DistrictPersian)
	fmt.Println("Year:", ad.Year)
	fmt.Println("AdId:", ad.AdId)
	fmt.Println("Age:", ad.Age)
	fmt.Println("FloorNumber:", ad.FloorNumber)
	fmt.Println("HasWareHouse:", ad.HasWareHouse)
	fmt.Println("HasElevator:", ad.HasElevator)
	fmt.Println("TotalPrice:", ad.TotalPrice)
	fmt.Println("PricePerMeter:", ad.PricePerMeter)
	fmt.Println("PrePaidPrice:", ad.PrePaidPrice)
	fmt.Println("MonthlyRentPrice:", ad.MonthlyRentPrice)
	fmt.Println("AdType:", ad.AdType)
	fmt.Println("HouseType:", ad.HouseType)
	fmt.Println("PublishedAt:", ad.PublishedAt)
	fmt.Println()
}

func crawl(url string) AdItem {
	data := make(map[string]string)

	// Send an HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch page: %d %s", res.StatusCode, res.Status)
	}

	// Parse the HTML using goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("[type='application/ld+json']").Each(func(index int, item *goquery.Selection) {
		data["script"] = item.Text()
	})

	// Check if script content is available
	jsonString, exists := data["script"]
	if !exists {
		log.Fatal("No JSON-LD script found")
	}

	// Unmarshal JSON into struct
	var apartments []AdItem
	err = json.Unmarshal([]byte(jsonString), &apartments)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	// Extract year using takeOutYear function
	apartments[0].Year = takeOutYear(doc)
	apartments[0].AdId = extractAdId(apartments[0].URL)

	doc.Find("title").Each(func(index int, item *goquery.Selection) {
		titleValue := item.Text()
		parts := strings.Split(titleValue, "-")
		apartments[0].PublishedAt = convertHumanDateToNormalDate(toEnglishDigits(parts[1]))
	})

	results := getSomeOtherProperties(doc)
	apartments[0].FloorNumber = toEnglishDigits(results["floor_number"].(string))
	apartments[0].HasWareHouse = results["has_warehouse"].(bool)
	apartments[0].HasElevator = results["has_elevator"].(bool)
	apartments[0].TotalPrice = results["total_price"].(string)
	apartments[0].PricePerMeter = results["price_per_meter"].(string)

	if results["prepaid_price"] != nil {
		apartments[0].PrePaidPrice = results["prepaid_price"].(string)
	} else {
		apartments[0].PrePaidPrice = ""
	}

	if results["monthly_rent_price"] != nil {
		apartments[0].MonthlyRentPrice = results["monthly_rent_price"].(string)
	} else {
		apartments[0].MonthlyRentPrice = ""
	}

	return apartments[0]
}

func getSomeOtherProperties(doc *goquery.Document) map[string]interface{} {
	htmlContent, err := doc.Html()
	if err != nil {
		log.Fatal("Error converting document to HTML:", err)
	}

	startPattern := `"LIST_DATA"\s*:\s*`
	endPattern := `\s*}\s*]\s*}\s*`

	data := cleanLastCurlyBrace(takeOutSomething(htmlContent, startPattern, endPattern))

	var widgets []Widget
	err = json.Unmarshal([]byte(data), &widgets)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Extract the required values
	var results = make(map[string]interface{})

	for _, widget := range widgets {
		// Check for titles and extract associated values
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

		// Extract values from credit and rent objects
		if widget.Credit != nil {
			results["total_price"] = widget.Credit.Value
		}
		if widget.Rent != nil {
			results["monthly_rent_price"] = widget.Rent.Value
		}

	}

	_, exist := results["has_warehouse"]
	if !exist {
		results["has_warehouse"] = false
	}

	_, exist = results["has_elevator"]
	if !exist {
		results["has_elevator"] = false
	}

	_, exist = results["floor_number"]
	if !exist {
		results["floor_number"] = "0"
	}

	return results
}

func cleanLastCurlyBrace(text string) string {
	// Trim spaces from the beginning and end of the string
	trimmedText := strings.TrimSpace(text)

	// Check if the last character is '}' and remove it if it is
	if len(trimmedText) > 0 && trimmedText[len(trimmedText)-1] == '}' {
		// Return the string without the last '}'
		return trimmedText[:len(trimmedText)-1]
	}

	return trimmedText // Return the original string if no '}' is found
}

func takeOutSomething(text string, startPattern string, endPattern string) string {
	regex := regexp.MustCompile(startPattern)
	startMatch := regex.FindStringIndex(text)

	if startMatch == nil {
		return ""
	}

	startIndex := startMatch[1]

	regex = regexp.MustCompile(endPattern)
	endMatch := regex.FindStringIndex(text[startIndex:])

	if endMatch == nil {
		return ""
	}

	endIndex := endMatch[1]
	return text[startIndex : startIndex+endIndex]
}

func takeOutYear(doc *goquery.Document) string {
	// Locate JSON section using regex after "widgetType": "GROUP_INFO_ROW"
	var jsonData string
	doc.Find("script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		scriptContent := s.Text()
		regex := regexp.MustCompile(`"widgetType":\s*"GROUP_INFO_ROW",\s*"items":\s*\[(.*?)\]`)
		match := regex.FindStringSubmatch(scriptContent)
		if len(match) > 1 {
			jsonData = `{"widgetType": "GROUP_INFO_ROW", "items": [` + match[1] + `]}`
			return false // Stop after finding the match
		}
		return true // Continue searching
	})

	if jsonData == "" {
		log.Fatal("Could not find the GROUP_INFO_ROW section in the HTML")
	}

	// Parse the JSON into the GroupInfoRow struct
	var groupInfo GroupInfoRow
	err := json.Unmarshal([]byte(jsonData), &groupInfo)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	// Find the value where title is "ساخت"
	for _, item := range groupInfo.Items {
		if item.Title == "ساخت" {
			year := item.Value
			// Trim spaces and remove "قبل از" from the beginning if present
			year = strings.TrimSpace(year)
			year = strings.Replace(year, "قبل از", "", -1)
			return strings.TrimSpace(year)
		}
	}

	return ""
}

// Define structs to match JSON structure
type AdItem struct {
	AccommodationCategory string `json:"accommodationCategory"`
	AdType                string
	HouseType             string
	Description           string    `json:"description"`
	FloorSize             FloorSize `json:"floorSize"`
	Geo                   Geo       `json:"geo"`
	Image                 string    `json:"image"`
	Name                  string    `json:"name"`
	NumberOfRooms         string    `json:"numberOfRooms"`
	URL                   string    `json:"url"`
	WebInfo               WebInfo   `json:"web_info"`
	Year                  string
	AdId                  string
	Age                   string
	FloorNumber           string
	HasWareHouse          bool
	HasElevator           bool
	TotalPrice            string
	PricePerMeter         string
	PrePaidPrice          string
	MonthlyRentPrice      string
	PublishedAt           string
}

func (ad *AdItem) clean() {
	ad.AccommodationCategory = cleanString(ad.AccommodationCategory)
	ad.Description = cleanString(ad.Description)
	ad.FloorSize.Value = cleanString(ad.FloorSize.Value)
	ad.Geo.Latitude = cleanString(ad.Geo.Latitude)
	ad.Geo.Longitude = cleanString(ad.Geo.Longitude)
	ad.Image = cleanString(ad.Image)
	ad.Name = cleanString(ad.Name)
	ad.NumberOfRooms = cleanString(ad.NumberOfRooms)
	ad.URL = cleanString(ad.URL)
	ad.WebInfo.CityPersian = cleanString(ad.WebInfo.CityPersian)
	ad.WebInfo.DistrictPersian = cleanString(ad.WebInfo.DistrictPersian)
	ad.Year = cleanString(ad.Year)
	ad.AdId = cleanString(ad.AdId)
}

func (ad *AdItem) prepareNumberOfRooms() {
	ad.NumberOfRooms = alphaToNumber(ad.NumberOfRooms)
}

func (ad *AdItem) computeAgeFromYear(currentYear int) {
	year, err := strconv.Atoi(ad.Year)
	if err != nil {
		return
	}

	ad.Age = strconv.Itoa(currentYear - year)
}

func (ad *AdItem) prepareTypes() {
	if strings.Contains(ad.AccommodationCategory, "اجاره") {
		ad.AdType = "rent"
	}

	if strings.Contains(ad.AccommodationCategory, "فروش") {
		ad.AdType = "buy"
	}

	if strings.Contains(ad.AccommodationCategory, "رهن") {
		ad.AdType = "mortgage"
	}

	if strings.Contains(ad.AccommodationCategory, "آپارتمان") {
		ad.HouseType = "apartment"
	}

	if strings.Contains(ad.AccommodationCategory, "خانه") {
		ad.HouseType = "villa"
	}

	if strings.Contains(ad.AccommodationCategory, "ویلا") {
		ad.HouseType = "villa"
	}
}

type FloorSize struct {
	Value string `json:"value"`
}

type Geo struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type WebInfo struct {
	CityPersian     string `json:"city_persian"`
	DistrictPersian string `json:"district_persian"`
}

type Item struct {
	Title    string `json:"title"`
	Value    string `json:"value"`
	ID       int    `json:"id"`
	IconName string `json:"iconName,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

type GroupInfoRow struct {
	WidgetType string `json:"widgetType"`
	Items      []Item `json:"items"`
}

func extractAdId(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

// ---------- utils

func cleanString(str string) string {
	return toEnglishDigits(convertArabicCharsToPersian(cleanArabicErab(str)))
}

func toEnglishDigits(numberStr string) string {
	persianDigits := []string{"۰", "۱", "۲", "۳", "۴", "۵", "۶", "۷", "۸", "۹"}
	arabicDigits := []string{"٠", "١", "٢", "٣", "٤", "٥", "٦", "٧", "٨", "٩"}
	englishDigits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	for i := range persianDigits {
		numberStr = strings.ReplaceAll(numberStr, persianDigits[i], englishDigits[i])
	}
	for i := range arabicDigits {
		numberStr = strings.ReplaceAll(numberStr, arabicDigits[i], englishDigits[i])
	}
	return numberStr
}

func cleanArabicErab(input string) string {
	erabs := []string{"ً", "ٌ", "ٍ", "َ", "ُ", "ِ", "ّ"}
	for _, erab := range erabs {
		input = strings.ReplaceAll(input, erab, "")
	}
	return input
}

func convertArabicCharsToPersian(input string) string {
	arabicChars := []string{"ي", "ك", "ى"}
	persianChars := []string{"ی", "ک", "ی"}

	for i := range arabicChars {
		input = strings.ReplaceAll(input, arabicChars[i], persianChars[i])
	}
	return input
}

func alphaToNumber(str string) string {
	numbersMap := map[string]string{
		"صفر":    "0",
		"یک":     "1",
		"دو":     "2",
		"سه":     "3",
		"چهار":   "4",
		"پنج":    "5",
		"شش":     "6",
		"هفت":    "7",
		"هشت":    "8",
		"نه":     "9",
		"ده":     "10",
		"یازده":  "11",
		"دوازده": "12",
		"سیزده":  "13",
		"چهارده": "14",
		"پانزده": "15",
		"شانزده": "16",
		"هفتده":  "17",
		"هجده":   "18",
		"نوزده":  "19",
		"بیست":   "20",
	}

	number, exist := numbersMap[str]
	if !exist {
		return ""
	}

	return number
}

type ListDataItem struct {
	WidgetType string `json:"widgetType"`
	Credit     struct {
		Value int `json:"value"`
	} `json:"credit"`
	Rent struct {
		Value int `json:"value"`
	} `json:"rent"`
}

type Data struct {
	ListData []ListDataItem `json:"LIST_DATA"`
}

// Define struct types for unmarshalling the JSON
type Widget struct {
	WidgetType string  `json:"widgetType"`
	Items      []Item  `json:"items,omitempty"`
	Credit     *Credit `json:"credit,omitempty"`
	Rent       *Rent   `json:"rent,omitempty"`
	Title      string  `json:"title,omitempty"`
	Value      string  `json:"value,omitempty"`
}

type Credit struct {
	Value            int `json:"value"`
	TransformedValue int `json:"transformedValue"`
}

type Rent struct {
	Value            int `json:"value"`
	TransformedValue int `json:"transformedValue"`
}

func convertHumanDateToNormalDate(persianDate string) string {
	parts := strings.Fields(persianDate)

	return fmt.Sprintf("%s-%s-%s", parts[2], monthNameToMonthNumber(parts[1]), parts[0])
}

func monthNameToMonthNumber(monthName string) string {
	monthName = strings.TrimSpace(monthName)

	months := map[string]int{
		"فروردین":  1,
		"اردیبهشت": 2,
		"خرداد":    3,
		"تیر":      4,
		"مرداد":    5,
		"شهریور":   6,
		"مهر":      7,
		"آبان":     8,
		"آذر":      9,
		"دی":       10,
		"بهمن":     11,
		"اسفند":    12,
	}

	number := months[monthName]
	if number < 10 {
		return fmt.Sprintf("0%d", number)
	} else {
		return fmt.Sprint("%d", number)
	}
}
