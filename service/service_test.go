package service

import (
	"errors"
	"github.com/kniepok/weatherAPI"
	"github.com/kniepok/weatherAPI/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_GetBookmarks(t *testing.T) {
	bookmarksStorageOK := new(mocks.BookmarksStorage)
	var locations []*weatherapi.Location
	bookmarksStorageOK.On("GetBookmarks", mock.Anything, mock.Anything).Return(locations, nil)
	bookmarksStorageErr := new(mocks.BookmarksStorage)
	bookmarksStorageErr.On("GetBookmarks", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	type fields struct {
		bookmarkStorage weatherapi.BookmarksStorage
		weatherStorage  weatherapi.WeatherStorage
		weatherProvider weatherapi.WeatherProvider
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "correctly working bookmark storage",
			fields: fields{
				bookmarkStorage: bookmarksStorageOK,
			},
			wantErr: false,
		},
		{
			name: "error storage",
			fields: fields{
				bookmarkStorage: bookmarksStorageErr,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(
				tt.fields.weatherStorage,
				tt.fields.bookmarkStorage,
				tt.fields.weatherProvider,
			)
			if _, err := s.GetBookmarks(); (err != nil) != tt.wantErr {
				t.Errorf("Service.GetBookmarks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_AddBookmark(t *testing.T) {
	bookmarksStorageOK := new(mocks.BookmarksStorage)
	bookmarksStorageOK.On("AddBookmark", mock.Anything, mock.Anything).Return(nil)
	bookmarksStorageErr := new(mocks.BookmarksStorage)
	bookmarksStorageErr.On("AddBookmark", mock.Anything, mock.Anything).Return(errors.New("error"))

	type fields struct {
		bookmarkStorage weatherapi.BookmarksStorage
		weatherStorage  weatherapi.WeatherStorage
		weatherProvider weatherapi.WeatherProvider
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "correctly working bookmark storage",
			fields: fields{
				bookmarkStorage: bookmarksStorageOK,
			},
			wantErr: false,
		},
		{
			name: "error storage",
			fields: fields{
				bookmarkStorage: bookmarksStorageErr,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(
				tt.fields.weatherStorage,
				tt.fields.bookmarkStorage,
				tt.fields.weatherProvider,
			)
			if err := s.AddBookmark("locationName"); (err != nil) != tt.wantErr {
				t.Errorf("Service.GetBookmarks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetStatistics(t *testing.T) {
	var stats *weatherapi.Statistics
	weatherStorageOK := new(mocks.WeatherStorage)
	weatherStorageOK.On("GetStatistics", mock.Anything, mock.Anything).Return(stats, nil)
	weatherStorageErr := new(mocks.WeatherStorage)
	weatherStorageErr.On("GetStatistics", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	type fields struct {
		bookmarkStorage weatherapi.BookmarksStorage
		weatherStorage  weatherapi.WeatherStorage
		weatherProvider weatherapi.WeatherProvider
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "correctly working weather storage",
			fields: fields{
				weatherStorage: weatherStorageOK,
			},
			wantErr: false,
		},
		{
			name: "error storage",
			fields: fields{
				weatherStorage: weatherStorageErr,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(
				tt.fields.weatherStorage,
				tt.fields.bookmarkStorage,
				tt.fields.weatherProvider,
			)
			if _, err := s.GetStatistics("locationName"); (err != nil) != tt.wantErr {
				t.Errorf("Service.GetStatistics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetWeather(t *testing.T) {
	// weather storage mocks
	weatherStorageOK := new(mocks.WeatherStorage)
	weatherStorageOK.On("StoreWeather", mock.Anything, mock.Anything).Return(nil)
	weatherStorageErr := new(mocks.WeatherStorage)
	weatherStorageErr.On("StoreWeather", mock.Anything, mock.Anything).Return(errors.New("error"))

	var w *weatherapi.Weather
	// weather provider mocks
	weatherProviderOK := new(mocks.WeatherProvider)
	weatherProviderOK.On("FetchWeather", mock.Anything, mock.Anything).Return(w, nil)
	weatherProviderErr := new(mocks.WeatherProvider)
	weatherProviderErr.On("FetchWeather", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	type fields struct {
		bookmarkStorage weatherapi.BookmarksStorage
		weatherStorage  weatherapi.WeatherStorage
		weatherProvider weatherapi.WeatherProvider
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "both services are ok",
			fields: fields{
				weatherStorage:  weatherStorageOK,
				weatherProvider: weatherProviderOK,
			},
			wantErr: false,
		},
		{
			name: "error storage",
			fields: fields{
				weatherStorage:  weatherStorageErr,
				weatherProvider: weatherProviderOK,
			},
			wantErr: true,
		},
		{
			name: "error provider",
			fields: fields{
				weatherStorage:  weatherStorageOK,
				weatherProvider: weatherProviderErr,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(
				tt.fields.weatherStorage,
				tt.fields.bookmarkStorage,
				tt.fields.weatherProvider,
			)
			if _, err := s.GetWeather("locationName"); (err != nil) != tt.wantErr {
				t.Errorf("Service.GetWeather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
