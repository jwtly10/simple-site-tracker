package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	. "github.com/jwtly10/simple-site-tracker/api/track"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlers_ServeTrackJSHandler(t *testing.T) {
	genTmpFile()

	mockRepo := &MockRepository{}
	handlers := NewHandlers(mockRepo)

	mockRepo.On("GetDomainIDFromKey", mock.Anything).Return(1, nil)

	req, err := http.NewRequest("GET", "/serve/js/123", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	handlers.ServeTrackJSHandler(recorder, req)

	assert.Contains(t, recorder.Body.String(), "123")

	assert.Equal(t, http.StatusOK, recorder.Code)

	delTmpFile()
}

func TestHandlers_ServeTrackJSHandler_InvalidClientKey(t *testing.T) {
	mockRepo := &MockRepository{}

	handlers := NewHandlers(mockRepo)

	mockRepo.On("GetDomainIDFromKey", mock.Anything).Return(0, nil)

	req, err := http.NewRequest("GET", "/serve/js/123", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	handlers.ServeTrackJSHandler(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func genTmpFile() {
	cwd, _ := os.Getwd()
	content := []byte("mocked file content %s, %s")
	_ = os.MkdirAll(filepath.Join(cwd, "templates"), 0755)
	filePath := filepath.Join(cwd, "templates", "clientScript.js")
	_ = os.WriteFile(filePath, []byte(content), 0644)
}

func delTmpFile() {
	cwd, _ := os.Getwd()
	filePath := filepath.Join(cwd, "templates", "clientScript.js")
	_ = os.Remove(filePath)
}
