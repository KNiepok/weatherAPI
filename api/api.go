package api

import (
	"github.com/go-chi/chi"
	"github.com/kniepok/weatherAPI"
	"log"
	"net"
	"net/http"
	"os"
)

type Api struct {
	router  chi.Router
	server  *http.Server
	service weatherapi.WeatherService
}

var (
	port = "8080"
	host = "0.0.0.0"
)

func init() {
	if val := os.Getenv("PORT"); len(val) > 0 {
		port = val
	}

	if val := os.Getenv("HOST"); len(val) > 0 {
		host = val
	}
}

func New(weatherService weatherapi.WeatherService) *Api {
	s := &Api{
		router:  chi.NewRouter(),
		service: weatherService,
	}

	s.SetupRoutes()
	return s
}

// ListenAndServe will listen for requests
func (s *Api) ListenAndServe() error {
	s.server = &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: s.router,
	}

	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	go func() {
		if err = s.server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
