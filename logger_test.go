package echozap

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapLogger(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/something", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}

	obs, logs := observer.New(zap.DebugLevel)
	config := &Config{
		Log: zap.New(obs),
		Skip2XX: false,
	}

	err := ZapLogger(config)(h)(c)
	assert.Nil(t, err)

	logFields := logs.AllUntimed()[0].ContextMap()

	assert.Equal(t, 1, logs.Len())
	assert.Equal(t, int64(200), logFields["status"])
	assert.NotNil(t, logFields["time"])
	assert.Equal(t, "GET /something", logFields["request"])
	assert.NotNil(t, logFields["host"])
	assert.NotNil(t, logFields["size"])
}

func TestSkip2XX(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/something", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}

	obs, logs := observer.New(zap.DebugLevel)
	config := &Config{
		Log: zap.New(obs),
		Skip2XX: true,
	}

	err := ZapLogger(config)(h)(c)
	assert.Nil(t, err)

	logFields := logs.AllUntimed()
	assert.Equal(t, 0, len(logFields))
}

