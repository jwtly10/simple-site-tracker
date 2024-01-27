package service

import (
	"github.com/jwtly10/simple-site-tracker/api/track"
	"github.com/jwtly10/simple-site-tracker/utils/logger"
)

type Service struct {
	repo track.Repository
}

func NewService(repo *track.Repository) *Service {
	return &Service{
		repo: *repo,
	}
}

// ValidateDomainKeyPair validates the domain and key pair.
// It returns true if the domain and key pair is valid.
func (s *Service) ValidateDomainKeyPair(domain string, siteKey string) bool {
	l := logger.Get()

	keyPair, err := s.repo.GetDomainKeyPair(domain)
	if err != nil {
		l.Error().Err(err).Msg("Error getting domain key pair")
		return false
	}

	l.Info().Msgf("Validating domain key pair: %s %s", keyPair.Domain, keyPair.SiteKey)

	if keyPair.Domain == "" {
		l.Error().Msg("Domain not found")
		return false
	}

	if keyPair.SiteKey != siteKey {
		l.Error().Msg("Invalid site key")
		return false
	}

	return true
}
