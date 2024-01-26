package track

import (
	"encoding/json"
	"net/http"

	"github.com/jwtly10/simple-site-tracker/utils/logger"
	"github.com/sirupsen/logrus"
)

func TrackUTMHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.Get()

	if r.Method != http.MethodGet {
		l.Warn().Msgf("Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.Info("Tracking UTMs")
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
	Element   string `json:"element"`
	PageUrl   string `json:"page_url"`
	Timestamp string `json:"timestamp"`
}

func TrackClickHandler(w http.ResponseWriter, r *http.Request) {
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

	l.Info().Msgf("element: %s", clickEvent.Element)
	l.Info().Msgf("page_url: %s", clickEvent.PageUrl)
	l.Info().Msgf("timestamp: %s", clickEvent.Timestamp)
}
