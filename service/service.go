package service

import (
	"github.com/kniepok/weatherAPI"
)

type Service struct {
	weatherService   weatherapi.WeatherStorage
	bookmarksService weatherapi.BookmarksStorage
	weatherProvider  weatherapi.WeatherProvider
}

func New(ws weatherapi.WeatherStorage, bs weatherapi.BookmarksStorage, wf weatherapi.WeatherProvider) *Service {
	s := &Service{
		ws, bs, wf,
	}

	return s
}

func (s *Service) AddBookmark(locationName string) error {
	return s.bookmarksService.AddBookmark(&weatherapi.Location{Name: locationName})
}

func (s *Service) GetBookmarks() ([]*weatherapi.Location, error) {
	return s.bookmarksService.GetBookmarks()
}

func (s *Service) GetWeather(locationName string) (*weatherapi.Weather, error) {
	weather, err := s.weatherProvider.FetchWeather(&weatherapi.Location{Name: locationName})

	if err != nil {
		return nil, err
	}

	err = s.weatherService.StoreWeather(weather)

	if err != nil {
		return nil, err
	}

	return weather, nil
}

func (s *Service) GetStatistics(locationName string) (*weatherapi.Statistics, error) {
	return s.weatherService.GetStatistics(&weatherapi.Location{Name: locationName})
}
