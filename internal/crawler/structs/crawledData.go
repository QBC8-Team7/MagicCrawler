package structs

type CrawledData struct {
	HouseType        string // (apartment,villa)
	AdCategory       string // (rent, buy, mortgage)
	Description      string
	Meterage         int
	Lat              string
	Lon              string
	ImageUrl         string
	Title            string
	RoomsCount       int
	URL              string
	Year             string
	Age              int
	PublisherAdKey   string
	FloorNumber      int
	HasWarehouse     bool
	HasElevator      bool
	TotalPrice       string
	MonthlyRentPrice string
	PricePerMeter    string
	PrePaidPrice     string
	PublishedAt      string
	City             string
	Neighborhood     string
	SourceName       string
	Author           string
}
