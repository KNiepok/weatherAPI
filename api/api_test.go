package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kniepok/weatherAPI"
	"github.com/kniepok/weatherAPI/mocks"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	exampleLocation = &weatherapi.Location{Name: "Piaseczno"}
	exampleWeather  = &weatherapi.Weather{
		LocationName: "Piaseczno",
		Date:         "2019-04-01",
		Type:         "cloudy",
		Temp:         7.34,
		MinTemp:      1.23,
		MaxTemp:      12.21,
	}

	exampleStats = &weatherapi.Statistics{
		QueriesCount: 11,
		Stats: map[string]*weatherapi.MonthlyStats{"03-2019": {
			Temperatures: &weatherapi.Temperatures{
				Average: 10.0,
				Lowest:  1,
				Highest: 100,
			},
			TypesOccurrences: map[weatherapi.WeatherType]int{"cloudy": 1},
		}},
	}
)

func TestAPI_GetBookmarks(t *testing.T) {
	var locations []*weatherapi.Location
	locations = append(locations, exampleLocation)

	WeatherServiceOK := new(mocks.WeatherService)
	WeatherServiceOK.On("GetBookmarks", mock.Anything, mock.Anything).Return(locations, nil)
	WeatherServiceErr := new(mocks.WeatherService)
	WeatherServiceErr.On("GetBookmarks", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	type fields struct {
		service weatherapi.WeatherService
	}
	type want struct {
		body interface{}
		code int
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "correct query and service result -> 200 and weather in response",
			fields: fields{
				service: WeatherServiceOK,
			},
			want: want{
				code: http.StatusOK,
				body: locations,
			},
		},
		{
			name: "service error -> 500 and appropriate error",
			fields: fields{
				service: WeatherServiceErr,
			},
			want: want{
				code: http.StatusInternalServerError,
				body: ErrServerError(errors.New("error")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantedBody, err := json.Marshal(tt.want.body)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("GET", "/bookmarks", nil)
			if err != nil {
				t.Fatal(err)
			}

			api := New(tt.fields.service)
			rr := httptest.NewRecorder()
			api.router.ServeHTTP(rr, req)

			if rr.Code != tt.want.code {
				t.Fatalf("invalid status code, wanted %v got %v", tt.want.code, rr.Code)
			}

			// strip newline at the end
			r := rr.Body.Bytes()[:len(rr.Body.Bytes())-1]
			if bytes.Compare(r, wantedBody) != 0 {
				t.Fatalf("got bad body, wanted %v got %v", string(wantedBody), rr.Body.String())
			}
		})
	}
}

func TestAPI_AddBookmark(t *testing.T) {
	WeatherServiceOK := new(mocks.WeatherService)
	WeatherServiceOK.On("AddBookmark", mock.Anything, mock.Anything).Return(nil)
	WeatherServiceErr := new(mocks.WeatherService)
	WeatherServiceErr.On("AddBookmark", mock.Anything, mock.Anything).Return(errors.New("error"))

	type fields struct {
		service weatherapi.WeatherService
	}
	type args struct {
		body interface{}
	}
	type want struct {
		code int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "correct query and service result -> 200 and weather in response",
			fields: fields{
				service: WeatherServiceOK,
			},
			args: args{
				body: exampleLocation,
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name: "service error -> 500 and appropriate error",
			fields: fields{
				service: WeatherServiceErr,
			},
			args: args{
				body: exampleLocation,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, err := json.Marshal(tt.args.body)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/bookmarks", bytes.NewReader(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			api := New(tt.fields.service)
			rr := httptest.NewRecorder()
			api.router.ServeHTTP(rr, req)

			if rr.Code != tt.want.code {
				t.Fatalf("invalid status code, wanted %v got %v", tt.want.code, rr.Code)
			}
		})
	}
}

func TestAPI_WeatherQuery(t *testing.T) {
	WeatherServiceOK := new(mocks.WeatherService)
	WeatherServiceOK.On("GetWeather", mock.Anything, mock.Anything).Return(exampleWeather, nil)
	WeatherServiceErr := new(mocks.WeatherService)
	WeatherServiceErr.On("GetWeather", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	type fields struct {
		service weatherapi.WeatherService
	}

	type args struct {
		body interface{}
	}

	type want struct {
		body interface{}
		code int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "correct query and service result -> 200 and weather in response",
			fields: fields{
				service: WeatherServiceOK,
			},
			args: args{
				body: &weatherapi.Location{Name: exampleWeather.LocationName},
			},
			want: want{
				code: http.StatusOK,
				body: exampleWeather,
			},
		},
		{
			name: "service error -> 500 and appropriate error",
			fields: fields{
				service: WeatherServiceErr,
			},
			args: args{
				body: &weatherapi.Location{Name: exampleWeather.LocationName},
			},
			want: want{
				code: http.StatusInternalServerError,
				body: ErrServerError(errors.New("error")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantedBody, err := json.Marshal(tt.want.body)
			if err != nil {
				t.Fatal(err)
			}

			requestBody, err := json.Marshal(tt.args.body)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/query", bytes.NewReader(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			api := New(tt.fields.service)
			rr := httptest.NewRecorder()
			api.router.ServeHTTP(rr, req)

			if rr.Code != tt.want.code {
				t.Fatalf("invalid status code, wanted %v got %v", tt.want.code, rr.Code)
			}

			// strip newline at the end
			r := rr.Body.Bytes()[:len(rr.Body.Bytes())-1]
			if bytes.Compare(r, wantedBody) != 0 {
				t.Fatalf("got bad body, wanted %v got %v", string(wantedBody), rr.Body.String())
			}
		})
	}
}

func TestAPI_GetStatistics(t *testing.T) {
	WeatherServiceOK := new(mocks.WeatherService)
	WeatherServiceOK.On("GetStatistics", mock.Anything, mock.Anything).Return(exampleStats, nil)
	WeatherServiceErr := new(mocks.WeatherService)
	WeatherServiceErr.On("GetStatistics", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	type fields struct {
		service weatherapi.WeatherService
	}

	type args struct {
		query string
	}

	type want struct {
		body interface{}
		code int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "correct query and service result -> 200 and weather in response",
			fields: fields{
				service: WeatherServiceOK,
			},
			args: args{
				query: "Piaseczno, Pl",
			},
			want: want{
				code: http.StatusOK,
				body: exampleStats,
			},
		},
		{
			name: "service error -> 500 and appropriate error",
			fields: fields{
				service: WeatherServiceErr,
			},
			args: args{
				query: "Piaseczno, Pl",
			},
			want: want{
				code: http.StatusInternalServerError,
				body: ErrServerError(errors.New("error")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantedBody, err := json.Marshal(tt.want.body)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("GET", fmt.Sprintf("/statistics?location=%s", tt.args.query), nil)
			if err != nil {
				t.Fatal(err)
			}

			api := New(tt.fields.service)
			rr := httptest.NewRecorder()
			api.router.ServeHTTP(rr, req)

			if rr.Code != tt.want.code {
				t.Fatalf("invalid status code, wanted %v got %v", tt.want.code, rr.Code)
			}

			// strip newline at the end
			r := rr.Body.Bytes()[:len(rr.Body.Bytes())-1]
			if bytes.Compare(r, wantedBody) != 0 {
				t.Fatalf("got bad body, wanted %v got %v", string(wantedBody), rr.Body.String())
			}
		})
	}
}
