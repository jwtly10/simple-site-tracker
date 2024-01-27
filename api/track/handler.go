package track

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"

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
	PageURL     string `json:"page_url"`
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
	utmId, err := h.repo.SaveUTM(pageId, utmEvent.UTMSource, utmEvent.UTMMedium, utmEvent.UTMCampaign)
	if err != nil {
		l.Error().Err(err).Msg("Error saving UTM")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l.Info().Msgf("UTM saved with ID %d", utmId)

	w.WriteHeader(http.StatusOK)
}

type TrackClickRequest struct {
	Element   string `json:"element"`
	PageUrl   string `json:"page_url"`
	Timestamp string `json:"timestamp"`
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

	if clickEvent.Element == "" || clickEvent.Timestamp == "" {
		l.Error().Msg("Invalid request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l.Info().Msgf("Click tracked")

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
