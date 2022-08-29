package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	appMiddleware "cometScraper/delivery/middleware"
	"cometScraper/entity"
	"cometScraper/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	mockLogger := new(mocks.Logger)
	cid := appMiddleware.NewMiddleware(mockLogger).RequestID()
	h := cid(handler)
	err := h(c)

	require.NoError(t, err)
	assert.NotNil(t, rec.Header().Get(entity.RequestIDHeader))
}
