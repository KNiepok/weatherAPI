package api

import "github.com/go-chi/chi/middleware"

func (s *Api) SetupRoutes() {

	s.router.Use(middleware.Logger)

	// Bookmarks
	s.router.Get("/bookmarks", s.GetBookmarks())
	s.router.Post("/bookmarks", s.AddBookmark())

	// Weather data
	s.router.Post("/query", s.FetchWeather())
	s.router.Get("/statistics", s.GetStatistics())
}
