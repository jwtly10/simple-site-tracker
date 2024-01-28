package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/jwtly10/simple-site-tracker/api/track"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlers_TrackUTMHandler(t *testing.T) {

	mockRepo := &MockRepository{}
	handlers := NewHandlers(mockRepo)

	mockRepo.On("GetDomain", mock.Anything).Return(1, nil)
	mockRepo.On("GetPage", mock.Anything, "/about").Return(1, nil)
	mockRepo.On("SaveUTM", 1, "test_source", "test_medium", "test_campaign", "test_track").Return(42, nil)

	data := `{"utm_source":"test_source","utm_medium":"test_medium","utm_campaign":"test_campaign","track":"test_track","page_url":"http://localhost:3000/about"}`

	req, err := http.NewRequest("POST", "/api/v1/track/utm", strings.NewReader(data))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost:3000")

	recorder := httptest.NewRecorder()

	handlers.TrackUTMHandler(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestHandlers_TrackPageViewHandler(t *testing.T) {
	mockRepo := &MockRepository{}
	handlers := NewHandlers(mockRepo)

	mockRepo.On("GetDomain", mock.Anything).Return(2, nil)
	mockRepo.On("GetPage", mock.Anything, mock.Anything).Return(3, nil)
	mockRepo.On("SavePageView", 2, 3).Return(42, nil)

	data := `{"page_url":"http://localhost:3000/about"}`

	req, err := http.NewRequest("POST", "/api/v1/track/pageview", strings.NewReader(data))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost:3000")

	recorder := httptest.NewRecorder()

	handlers.TrackPageViewHandler(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestHandlers_TrackClickHandler(t *testing.T) {
	mockRepo := &MockRepository{}
	handlers := NewHandlers(mockRepo)

	mockRepo.On("GetDomain", mock.Anything).Return(2, nil)
	mockRepo.On("GetPage", mock.Anything, mock.Anything).Return(3, nil)
	mockRepo.On("SaveClick", 3, mock.Anything).Return(42, nil)

	data := `{"element":{"tag":"span","id":"","classList":[],"textContent":"Generate video","parentElement":{"tag":"button","id":"","classList":["ant-btn","css-dev-only-do-not-override-6ynzfo","ant-btn-primary","generate"],"textContent":"Generate video"}},"url":"http://localhost:5173/generate"}`

	req, err := http.NewRequest("POST", "/api/v1/track/click", strings.NewReader(data))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost:5173")

	recorder := httptest.NewRecorder()

	handlers.TrackClickHandler(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestHandlers_MissingOriginHeader(t *testing.T) {
	mockRepo := &MockRepository{}
	handlers := NewHandlers(mockRepo)

	data := `{"utm_source":"test_source","utm_medium":"test_medium","utm_campaign":"test_campaign","track":"test_track","page_url":""}`

	req, err := http.NewRequest("POST", "/api/v1/track/utm", strings.NewReader(data))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	// Missing origin header

	recorder := httptest.NewRecorder()

	handlers.TrackUTMHandler(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestHandlers_InvalidMethod(t *testing.T) {
	mockRepo := &MockRepository{}
	handlers := NewHandlers(mockRepo)

	req, err := http.NewRequest("GET", "/api/v1/track/utm", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	handlers.TrackUTMHandler(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandlers_InvalidBody(t *testing.T) {
	mockRepo := &MockRepository{}
	handlers := NewHandlers(mockRepo)

	data := `{"utm_source":"test_source","utm_medium":"test_medium","utm_campaign":"test_campaign","track":"test_track","page_url":"http://"`

	req, err := http.NewRequest("POST", "/api/v1/track/utm", strings.NewReader(data))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost:3000")

	recorder := httptest.NewRecorder()

	handlers.TrackUTMHandler(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
