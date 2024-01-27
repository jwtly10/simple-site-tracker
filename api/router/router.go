package router

import (
	"net/http"

	"github.com/jwtly10/simple-site-tracker/api/track"
	"github.com/jwtly10/simple-site-tracker/utils/logger"
)

type Route struct {
	Path    string
	Handler http.HandlerFunc
}

type Routes []Route

func NewRouter(trackHandlers *track.Handlers) *http.ServeMux {
	router := http.NewServeMux()

	routes := Routes{
		{Path: "/api/v1/track/utm", Handler: logRequest(track.TrackUTMHandler)},
		{Path: "/api/v1/track/click", Handler: logRequest(trackHandlers.TrackClickHandler)},
	}

	for _, route := range routes {
		router.HandleFunc(route.Path, route.Handler)
	}

	return router
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	l := logger.Get()
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info().Msgf("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
