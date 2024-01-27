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
		{Path: "/serve/js/", Handler: trackHandlers.ServeJSHandler},
	}

	allowedOrigins := []string{"http://localhost:5173"}

	for _, route := range routes {
		corsHandler := handleCORS(allowedOrigins, route.Handler)
		router.HandleFunc(route.Path, corsHandler)
	}

	return router
}

func handleCORS(allowedOrigins []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		referrer := r.Header.Get("Referer")
		if referrer == "" {
			http.Error(w, "Missing Referer header", http.StatusUnauthorized)
			return
		}
		if referrer[len(referrer)-1] == '/' {
			referrer = referrer[:len(referrer)-1]
		}

		originAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if referrer == allowedOrigin {
				originAllowed = true
				break
			}
		}

		if r.Method == "OPTIONS" {

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Site-Key, Origin")

			if originAllowed {
				w.Header().Set("Access-Control-Allow-Origin", referrer)
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		if originAllowed {
			w.Header().Set("Access-Control-Allow-Origin", referrer)
		} else {
			http.Error(w, "Origin not allowed", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
