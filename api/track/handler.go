package track

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"

	"github.com/jwtly10/simple-site-tracker/utils/logger"
)

type Handlers struct {
	repo *Repository
}

func NewHandlers(repo *Repository) *Handlers {
	return &Handlers{repo: repo}
}

type TrackUTMRequest struct {
	UTMSource   string `json:"utm_source"`
	UTMMedium   string `json:"utm_medium"`
	UTMCampaign string `json:"utm_campaign"`
	Track       string `json:"track"`
	PageURL     string `json:"page_url"`
}

// ServeJSHandler serves the JS file for the specific domain
func (h *Handlers) ServeJSHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.Get()
	l.Info().Msg("Serving JS")

	clientKey := r.URL.Path[len("/serve/js/"):]
	if clientKey == "" {
		l.Error().Msg("Missing client key")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check valid clientKey
	domainId, err := h.repo.GetDomainIDFromKey(clientKey)
	if err != nil {
		l.Error().Err(err).Msg("Error getting domain ID from key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if domainId == 0 {
		l.Error().Msg("Invalid client key")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	script := GenerateClientJS(clientKey)
	if script == "" {
		l.Error().Msg("Error generating client JS")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l.Info().Msgf("Serving JS for client key %s", clientKey)
	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte(script))
}

// GenerateClientJS generates the client JS script.
// It returns the client JS script.
func GenerateClientJS(clientKey string) string {
	l := logger.Get()
	templatePath := filepath.Join("templates", "clientScript.js")
	fileContent, err := os.ReadFile(templatePath)
	if err != nil {
		l.Error().Err(err).Msg("Error reading file")
		return ""
	}

	serverURL := os.Getenv("SERVER_URL")
	formattedContent := fmt.Sprintf(string(fileContent), clientKey, serverURL)

	// Minify the JS
	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)

	minified, err := m.String("text/javascript", formattedContent)
	if err != nil {
		return ""
	}

	return minified
}

// TrackUTMHandler handles tracking UTMs.
// It saves the UTMs to the database.
// It returns a 200 status code if successful.
func (h *Handlers) TrackUTMHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.Get()

	if r.Method != http.MethodPost {
		l.Warn().Msgf("Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var utmEvent TrackUTMRequest
	err := json.NewDecoder(r.Body).Decode(&utmEvent)
	if err != nil {
		l.Error().Msgf("Error decoding request: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l.Info().Msgf("Tracking UTMs for request %s", utmEvent)

	// Get the domain of the request
	origin := r.Header.Get("Origin")
	if origin == "" {
		l.Error().Msg("Missing Origin header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	domain := getDomainFromOrigin(origin)

	// We can ignore the error here because we already validated the domain
	// TODO ensure this is true
	domainId, _ := h.repo.GetDomain(domain)
	if domainId == 0 {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// Get the page ID
	page := getPageFromURL(utmEvent.PageURL)
	pageId, err := h.repo.GetPage(domainId, page)
	if err != nil {
		l.Error().Err(err).Msg("Error getting page")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create page if it doesn't exist
	var pageId64 int64
	if pageId == 0 {
		l.Info().Msgf("Page %s does not exist. Creating page", page)
		pageId64, err = h.repo.CreatePage(domainId, page)
		if err != nil {
			l.Error().Err(err).Msg("Error creating page")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		pageId = int(pageId64)
	}

	// Save UTM
	l.Info().Msgf("Saving UTM for page %s", page)
	utmId, err := h.repo.SaveUTM(pageId, utmEvent.UTMSource, utmEvent.UTMMedium, utmEvent.UTMCampaign, utmEvent.Track)
	if err != nil {
		l.Error().Err(err).Msg("Error saving UTM")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l.Info().Msgf("UTM tracked with ID %d", utmId)
	w.WriteHeader(http.StatusOK)
}

type TrackClickRequest struct {
	Element map[string]interface{} `json:"element"`
	URL     string                 `json:"url"`
}

// TrackClickHandler handles tracking clicks.
// It saves the click to the database.
// It returns a 200 status code if successful.
func (h *Handlers) TrackClickHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.Get()
	if r.Method != http.MethodPost {
		l.Warn().Msgf("Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l.Info().Msg("Tracking clicks")

	var clickEvent TrackClickRequest
	err := json.NewDecoder(r.Body).Decode(&clickEvent)
	if err != nil {
		l.Error().Msgf("Error decoding request: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the domain id
	origin := r.Header.Get("Origin")
	if origin == "" {
		l.Error().Msg("Missing Origin header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	domain := getDomainFromOrigin(origin)

	domainId, err := h.repo.GetDomain(domain)
	if domainId == 0 {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// Get the page id
	page := getPageFromURL(clickEvent.URL)
	pageId, err := h.repo.GetPage(domainId, page)
	if err != nil {
		l.Error().Err(err).Msg("Error getting page")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create page if it doesn't exist
	var pageId64 int64
	if pageId == 0 {
		l.Info().Msgf("Page %s does not exist. Creating page", page)
		pageId64, err = h.repo.CreatePage(domainId, page)
		if err != nil {
			l.Error().Err(err).Msg("Error creating page")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		pageId = int(pageId64)
	}
	// Save the click
	l.Info().Msgf("Saving click for page %s", page)
	clickId, err := h.repo.SaveClick(pageId, clickEvent.Element)
	if err != nil {
		l.Error().Err(err).Msg("Error saving click")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l.Info().Msgf("Click tracked with ID %d", clickId)
	w.WriteHeader(http.StatusOK)
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

func getPageFromURL(pageURL string) string {
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		return ""
	}

	path := parsedURL.Path
	if path == "" {
		return "/"
	}

	return path
}
