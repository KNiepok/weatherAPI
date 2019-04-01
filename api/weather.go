package api

import (
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
)

func (s *Api) FetchWeather() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body *requestBody

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		weatherData, err := s.service.GetWeather(body.Name)

		if err != nil {
			_ = render.Render(w, r, ErrServerError(err))
			return
		}

		render.JSON(w, r, weatherData)
	}
}

func (s *Api) GetStatistics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		locationName := r.URL.Query().Get("location")

		if locationName == "" {
			_ = render.Render(w, r, ErrInvalidRequest(&errorString{"location cannot be empty."}))
			return
		}

		stats, err := s.service.GetStatistics(locationName)

		if err != nil {
			_ = render.Render(w, r, ErrServerError(err))
			return
		}

		render.JSON(w, r, stats)
	}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
