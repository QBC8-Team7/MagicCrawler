package divar

type GeneralFields struct {
	AccommodationCategory string    `json:"accommodationCategory"`
	Description           string    `json:"description"`
	FloorSize             FloorSize `json:"floorSize"`
	Geo                   Geo       `json:"geo"`
	Image                 string    `json:"image"`
	Name                  string    `json:"name"`
	NumberOfRooms         string    `json:"numberOfRooms"`
	URL                   string    `json:"url"`
	WebInfo               WebInfo   `json:"web_info"`
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
