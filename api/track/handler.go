package track

import (
	"encoding/json"
	"net/http"

	"github.com/jwtly10/simple-site-tracker/utils/logger"
)

type Handlers struct {
	repo *Repository
}

func NewHandlers(repo *Repository) *Handlers {
	return &Handlers{repo: repo}
}

// TrackUTMHandler handles tracking UTMs.
// It saves the UTMs to the database.
// It returns a 200 status code if successful.
func TrackUTMHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.Get()

	if r.Method != http.MethodGet {
		l.Warn().Msgf("Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l.Info().Msg("Tracking UTMs")

	utmSource := r.URL.Query().Get("utm_source")
	utmMedium := r.URL.Query().Get("utm_medium")
	utmCampaign := r.URL.Query().Get("utm_campaign")

	if utmSource == "" && utmMedium == "" && utmCampaign == "" {
		l.Warn().Msg("No UTMs found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l.Info().Msgf("utm_source: %s", utmSource)
	l.Info().Msgf("utm_medium: %s", utmMedium)
	l.Info().Msgf("utm_campaign: %s", utmCampaign)

	w.WriteHeader(http.StatusOK)
}

type TrackClickRequest struct {
	Domain    string `json:"domain"`
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

	domainID, err := h.repo.GetDomain(clickEvent.Domain)
	if err != nil {
		l.Error().Msgf("Error getting domain: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	l.Info().Msgf("domain_id: %d", domainID)
	l.Info().Msgf("domain: %s", clickEvent.Domain)
	l.Info().Msgf("element: %s", clickEvent.Element)
	l.Info().Msgf("page_url: %s", clickEvent.PageUrl)
	l.Info().Msgf("timestamp: %s", clickEvent.Timestamp)
}

func validateDomainAgainstKey(r *http.Request) (bool, error) {
	l := logger.Get()

	domain := r.URL.Query().Get("domain")
	key := r.URL.Query().Get("key")

}
