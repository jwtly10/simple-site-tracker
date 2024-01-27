package router

import (
	"net/http"

	"github.com/jwtly10/simple-site-tracker/api/router/middleware"
	"github.com/jwtly10/simple-site-tracker/api/track"
)

type Route struct {
	Path    string
	Handler http.HandlerFunc
}

type Routes []Route

func NewRouter(trackHandlers *track.Handlers, middleware *middleware.Middleware) *http.ServeMux {
	router := http.NewServeMux()

	routes := Routes{
		{Path: "/api/v1/track/utm", Handler: middleware.HandleMiddleware(trackHandlers.TrackUTMHandler, middleware.DomainValidation, middleware.LogRequest)},
		{Path: "/api/v1/track/click", Handler: middleware.HandleMiddleware(trackHandlers.TrackClickHandler, middleware.DomainValidation, middleware.LogRequest)},
		// Serve the JS with a variable int he URL
		{Path: "/serve/js/", Handler: trackHandlers.ServeJSHandler},
	}

	for _, route := range routes {
		corsHandler := handleCORS(route.Handler)
		router.HandleFunc(route.Path, corsHandler)
	}

	return router
}

func handleCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Site-Key, Origin")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
