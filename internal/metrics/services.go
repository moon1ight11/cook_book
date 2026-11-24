package metrics

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (m *Metrics) EchoMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			m.RequestsInFlight.Inc()
			defer m.RequestsInFlight.Dec()

			start := time.Now()

			err := next(c)

			duration := time.Since(start).Seconds()

			endpoint := c.Path()
			method := c.Request().Method
			status := strconv.Itoa(c.Response().Status)

			m.RequestDuration.WithLabelValues(method, endpoint).Observe(duration)
			m.RequestsTotal.WithLabelValues(method, endpoint, status).Inc()

			return err
		}
	}
}

func (m *Metrics) RecordOperation(operationType string) {
	m.OperationsTotal.WithLabelValues(operationType).Inc()
}

func (m *Metrics) RecordError(errorType, handler string) {
	m.ErrorsTotal.WithLabelValues(errorType, handler).Inc()
}

func (m *Metrics) RegisterMetricsHandler(e *echo.Echo, path string) {
	e.GET(path, func(c echo.Context) error {
		promhttp.Handler().ServeHTTP(c.Response(), c.Request())
		return nil
	})
}
