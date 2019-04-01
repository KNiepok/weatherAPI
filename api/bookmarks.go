package api

import (
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
)

type requestBody struct {
	Name string
}

func (s *Api) AddBookmark() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body *requestBody

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		if err := s.service.AddBookmark(body.Name); err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		render.NoContent(w, r)
	}
}

func (s *Api) GetBookmarks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locations, err := s.service.GetBookmarks()
		if err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		render.JSON(w, r, locations)
	}
}
