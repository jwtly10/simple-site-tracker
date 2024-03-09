package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/jwtly10/simple-site-tracker/api/middleware"
	"github.com/jwtly10/simple-site-tracker/api/track"
	"golang.org/x/time/rate"
)

type Route struct {
	Path    string
	Handler http.HandlerFunc
}

type Routes []Route

func NewRouter(trackHandlers *track.Handlers, middleware *middleware.Middleware) *http.ServeMux {
	router := http.NewServeMux()

	//  Max 50 requests per hour
	allowedReqPerHour:= 50
	secondsPerHour := 3600
	ratePerSecond := allowedReqPerHour / secondsPerHour
	burst := 50

	limiter := rate.NewLimiter(rate.Limit(ratePerSecond), burst)

	routes := Routes{
		{Path: "/api/v1/track/utm", Handler: middleware.HandleMiddleware(
			middleware.RateLimit(trackHandlers.TrackUTMHandler, limiter),
			middleware.DomainValidation,
			middleware.LogRequest)},
		{Path: "/api/v1/track/click", Handler: middleware.HandleMiddleware(
			middleware.RateLimit(trackHandlers.TrackClickHandler, limiter),
			middleware.DomainValidation,
			middleware.LogRequest)},
		{Path: "/api/v1/track/pageview", Handler: middleware.HandleMiddleware(
			middleware.RateLimit(trackHandlers.TrackPageViewHandler, limiter),
			middleware.DomainValidation,
			middleware.LogRequest)},
		{Path: "/serve/js/", Handler: middleware.HandleMiddleware(
			middleware.RateLimit(trackHandlers.ServeTrackJSHandler, limiter),
			middleware.CheckForIgnoreHeader)},
	}

	origins := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins := strings.Split(origins, ",")

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
