package middleware

import (
	"net"
	"net/http"
	"net/url"

	"golang.org/x/time/rate"

	"github.com/jwtly10/simple-site-tracker/api/service"
	"github.com/jwtly10/simple-site-tracker/utils/logger"
)

type Middleware struct {
	service *service.Service
}

func NewMiddleware(svc *service.Service) *Middleware {
	return &Middleware{
		service: svc,
	}
}

func (m *Middleware) HandleMiddleware(next http.HandlerFunc, mws ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, mw := range mws {
		next = mw(next)
	}
	return next
}

// LogRequest logs the request.
func (m *Middleware) LogRequest(next http.HandlerFunc) http.HandlerFunc {
	l := logger.Get()
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info().Msgf("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

// RateLimit limits the number of requests per second.
// Limits defined in router config
func (m *Middleware) RateLimit(next http.HandlerFunc, limiter *rate.Limiter) http.HandlerFunc {
	l := logger.Get()
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			l.Error().Msg("Rate limit exceeded")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Checks for the ignore header.
// Ignores the request if the X-Ignore-Tracking header is set to true.
// This is mainly used for testing live sites, to add the header you can install an extension like
// https://modheader.com/ on your browser.
func (m *Middleware) CheckForIgnoreHeader(next http.HandlerFunc) http.HandlerFunc {
	l := logger.Get()
	return func(w http.ResponseWriter, r *http.Request) {
		ignore := r.Header.Get("X-Ignore-Tracking")
		if ignore == "true" {
			l.Info().Msgf("Ignoring request from host: %s", r.Host)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// DomainValidation validates the domain and key pair.
// It returns a 400 status code if the domain and key pair is invalid.
func (m *Middleware) DomainValidation(next http.HandlerFunc) http.HandlerFunc {
	l := logger.Get()
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			http.Error(w, "Missing Origin header", http.StatusUnauthorized)
			return
		}

		l.Info().Msgf("Validating origin: %s", origin)

		siteKey := r.Header.Get("X-Site-Key")
		if siteKey == "" {
			http.Error(w, "Missing X-Site-Key header", http.StatusUnauthorized)
			return
		}

		domain := getDomainFromOrigin(origin)

		if !m.service.ValidateDomainKeyPair(domain, siteKey) {
			http.Error(w, "Invalid domain key pair", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// getDomainFromOrigin returns the domain from the origin.
func getDomainFromOrigin(origin string) string {
	u, err := url.Parse(origin)
	if err != nil {
		return ""
	}

	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return u.Host
	}

	return host
}
