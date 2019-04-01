package weatherapi

type WeatherService interface {
	AddBookmark(locationName string) error
	GetBookmarks() ([]*Location, error)
	GetWeather(locationName string) (*Weather, error)
	GetStatistics(locationName string) (*Statistics, error)
}

type WeatherProvider interface {
	FetchWeather(*Location) (*Weather, error)
}

type WeatherStorage interface {
	GetStatistics(*Location) (*Statistics, error)
	StoreWeather(*Weather) error
}

type BookmarksStorage interface {
	AddBookmark(*Location) error
	GetBookmarks() ([]*Location, error)
}

type Location struct {
	Name string `json:"name"`
}

type WeatherType string

type Weather struct {
	LocationName string      `json:"locationName"` // ie Washington, DC
	Date         string      `json:"date"`         // YYYY-MM-dd
	Type         WeatherType `json:"type"`         // cloudy, drizzled, etc
	Temp         float64     `json:"temp"`         // ie 7.34
	MinTemp      float64     `json:"min_temp"`     // ie 1.23
	MaxTemp      float64     `json:"max_temp"`     // ie 12.21
}

type Statistics struct {
	QueriesCount int                      `json:"queriesCount"`
	Stats        map[string]*MonthlyStats `json:"statistics"`
}

type MonthlyStats struct {
	Temperatures     *Temperatures       `json:"temperatures"`
	TypesOccurrences map[WeatherType]int `json:"typesOccurrences"`
}

type Temperatures struct {
	Average float64 `json:"average"`
	Lowest  float64 `json:"lowest"`
	Highest float64 `json:"highest"`
}
