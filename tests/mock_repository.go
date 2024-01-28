package tests

import (
	. "github.com/jwtly10/simple-site-tracker/api/track"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SavePageView(domainId, pageId int) (int64, error) {
	args := m.Called(domainId, pageId)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockRepository) SaveDomain(domain, key string) (int64, error) {
	args := m.Called(domain, key)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockRepository) GetDomain(domain string) (int, error) {
	args := m.Called(domain)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetDomainIDFromKey(key string) (int, error) {
	args := m.Called(key)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetDomainKeyPair(domain string) (DomainKeyPair, error) {
	args := m.Called(domain)
	return args.Get(0).(DomainKeyPair), args.Error(1)
}

func (m *MockRepository) GetPage(domainID int, pageURL string) (int, error) {
	args := m.Called(domainID, pageURL)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) CreatePage(domainID int, pageURL string) (int64, error) {
	args := m.Called(domainID, pageURL)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockRepository) SaveIPAddress(ipAddress string) (int64, error) {
	args := m.Called(ipAddress)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockRepository) SaveUTM(pageID int, utmSource, utmMedium, utmCampaign, track string) (int64, error) {
	args := m.Called(pageID, utmSource, utmMedium, utmCampaign, track)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockRepository) SaveClick(pageID int, element map[string]interface{}) (int64, error) {
	args := m.Called(pageID, element)
	return int64(args.Int(0)), args.Error(1)
}
