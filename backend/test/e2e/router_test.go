//go:build !e2e

package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/router"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthEndpoint(t *testing.T) {
	logger := zerolog.New(nil)
	cfg := common.LoadConfig()

	r := router.NewRouter(cfg, logger, nil, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestRouter_Setup(t *testing.T) {
	logger := zerolog.New(nil)
	cfg := common.LoadConfig()

	r := router.NewRouter(cfg, logger, nil, nil)
	require.NotNil(t, r)

	assert.IsType(t, &chi.Mux{}, r)
}
