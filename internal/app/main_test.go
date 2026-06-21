package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"argument/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestIndexRoute(t *testing.T) {
	config.Conf, _ = config.LoadConfig("../../config.toml")
	config.Conf.Root = "../../"
	router := SetupHandler()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "New Post")
}
