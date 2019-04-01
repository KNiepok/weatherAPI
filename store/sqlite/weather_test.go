package sqlite

import (
	"database/sql"
	"github.com/kniepok/weatherAPI"
	"io/ioutil"
	"os"
	"testing"
)

var (
	testWeather = &weatherapi.Weather{
		LocationName: "Washington, DC",
		Date:         "2019-04-01",
		Type:         "cloudy",
		Temp:         7.34,
		MinTemp:      1.23,
		MaxTemp:      12.21,
	}

	location = &weatherapi.Location{Name: "Washington, DC"}
)

func TestWeatherStorage(t *testing.T) {
	f, err := ioutil.TempFile(".", "")
	if err != nil {
		t.Fatalf("could not create file for db: %v", err)
	}
	f.Close()
	defer os.Remove(f.Name())

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		t.Fatalf("could not open db: %v", err)
	}

	storage, err := NewWeatherStorage(db)
	if err != nil {
		t.Fatalf("could not create storage: %v", err)
	}

	defer storage.Close()

	err = storage.StoreWeather(testWeather)
	if err != nil {
		t.Fatalf("could not add weather result to storage: %v", err)
	}

	stats, err := storage.GetStatistics(location)
	if err != nil {
		t.Fatalf("could not get stats: %v", err)
	}

	monthMap := make(map[string]*weatherapi.MonthlyStats)
	occurrences := make(map[weatherapi.WeatherType]int)
	occurrences[testWeather.Type] = 1

	monthMap["04-2019"] = &weatherapi.MonthlyStats{
		Temperatures: &weatherapi.Temperatures{
			Average: testWeather.Temp,
			Lowest:  testWeather.MinTemp,
			Highest: testWeather.MaxTemp,
		},
		TypesOccurrences: occurrences,
	}
	assertedStats := &weatherapi.Statistics{
		QueriesCount: 1,
		Stats:        monthMap,
	}

	if assertedStats.QueriesCount != stats.QueriesCount {
		t.Fatalf("invalid queries count: expected %v , got %v", assertedStats.QueriesCount, stats.QueriesCount)
	}

	if assertedStats.Stats["04-2019"].Temperatures.Average != stats.Stats["04-2019"].Temperatures.Average {
		t.Fatalf("invalid queries count: expected %v , got %v", assertedStats.QueriesCount, stats.QueriesCount)
	}
}
