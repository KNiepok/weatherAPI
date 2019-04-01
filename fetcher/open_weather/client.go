package open_weather

import (
	owm "github.com/briandowns/openweathermap"
	"github.com/kniepok/weatherAPI"
	"time"
)

type Client struct {
	apiKey string
}

func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) FetchWeather(location *weatherapi.Location) (*weatherapi.Weather, error) {
	w, err := owm.NewCurrent("C", "EN", c.apiKey)
	if err != nil {
		return nil, err
	}

	if err = w.CurrentByName(location.Name); err != nil {
		return nil, err
	}

	const dateLayout = "2006-01-02" // YYYY-MM-dd format

	return &weatherapi.Weather{
		LocationName: location.Name,
		Date:         time.Now().Format(dateLayout),
		Type:         weatherapi.WeatherType(w.Weather[0].Description),
		Temp:         w.Main.Temp,
		MinTemp:      w.Main.TempMin,
		MaxTemp:      w.Main.TempMax,
	}, nil
}
