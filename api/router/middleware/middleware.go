package middleware

import (
	"net"
	"net/http"
	"net/url"

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

// DomainValidation validates the domain and key pair.
// It returns a 400 status code if the domain and key pair is invalid.
func (m *Middleware) DomainValidation(next http.HandlerFunc) http.HandlerFunc {
	l := logger.Get()
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			http.Error(w, "Missing Origin header", http.StatusBadRequest)
			return
		}
		l.Info().Msgf("Validating origin: %s", origin)

		siteKey := r.Header.Get("X-Site-Key")
		if siteKey == "" {
			http.Error(w, "Missing X-Site-Key header", http.StatusBadRequest)
			return
		}

		domain := getDomainFromOrigin(origin)

		if !m.service.ValidateDomainKeyPair(domain, siteKey) {
			http.Error(w, "Invalid domain key pair", http.StatusBadRequest)
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
