package sqlite

import (
	"database/sql"
	"github.com/kniepok/weatherAPI"
	"math"
)

type WeatherStorage struct {
	db *sql.DB
}

func NewWeatherStorage(db *sql.DB) (*WeatherStorage, error) {
	bs := &WeatherStorage{
		db: db,
	}

	err := bs.migrate()

	if err != nil {
		return nil, err
	}

	return bs, nil
}

func (ws *WeatherStorage) StoreWeather(w *weatherapi.Weather) error {
	stmt, err := ws.db.Prepare("INSERT INTO weather_queries(location_name,date,type,temp,min_temp,max_temp) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(w.LocationName, w.Date, w.Type, w.Temp, w.MinTemp, w.MaxTemp); err != nil {
		return err
	}
	return nil
}

// Return statistics based on location
func (ws *WeatherStorage) GetStatistics(l *weatherapi.Location) (*weatherapi.Statistics, error) {
	rows, err := ws.db.Query("SELECT strftime('%m-%Y', date) as month, count(*) as queries_count, AVG(temp) as average_temp, MIN(min_temp) as min_temp, MAX(max_temp) as max_temp FROM weather_queries WHERE location_name=? GROUP BY month ", l.Name)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	queriesCount := 0

	months := make(map[string]*weatherapi.MonthlyStats)

	for rows.Next() {
		var st weatherapi.MonthlyStats
		var t weatherapi.Temperatures
		var count int
		var month string
		if err := rows.Scan(&month, &count, &t.Average, &t.Lowest, &t.Highest); err != nil {
			return nil, err
		}

		t.Average = math.Round(t.Average*100) / 100

		st.Temperatures = &t
		queriesCount += count
		months[month] = &st
	}

	occ, err := ws.getOccurrences(l)

	for m, oc := range occ {
		months[m].TypesOccurrences = oc
	}

	if err != nil {
		return nil, err
	}

	return &weatherapi.Statistics{
		QueriesCount: queriesCount,
		Stats:        months,
	}, nil

}

// Get weather type occurrences in form of map: [month -> [type -> occurencesCount]]
func (ws *WeatherStorage) getOccurrences(l *weatherapi.Location) (map[string]map[weatherapi.WeatherType]int, error) {
	rows, err := ws.db.Query("SELECT strftime('%m-%Y', date) as month, count(*) as queries_count, type FROM weather_queries WHERE location_name=? GROUP BY month,type;", l.Name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mapped := make(map[string]map[weatherapi.WeatherType]int)

	for rows.Next() {
		var month string
		var queriesCount int
		var weatherType weatherapi.WeatherType
		if err := rows.Scan(&month, &queriesCount, &weatherType); err != nil {
			return nil, err
		}
		if mapped[month] == nil {
			mapped[month] = make(map[weatherapi.WeatherType]int)
		}

		mapped[month][weatherType] += queriesCount
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mapped, nil
}

func (ws *WeatherStorage) migrate() error {
	_, err := ws.db.Exec("CREATE TABLE IF NOT EXISTS weather_queries (id integer not null primary key, location_name TEXT, date DATE,type TEXT, temp FLOAT(7,4), min_temp FLOAT(7,4), max_temp FLOAT(7,4));")
	return err
}
