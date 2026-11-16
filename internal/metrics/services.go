package metrics

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// мидлвар для сбора метрик
func (m *Metrics) EchoMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// увеличиваем счетчик текущих запросов и гарантированно его уменьшаем при выходе
			m.RequestsInFlight.Inc()
			defer m.RequestsInFlight.Dec()

			// засекаем время начала обработки
			start := time.Now()

			// идем дальше в логику приложения
			err := next(c)

			// высчитываем длительность запроса
			duration := time.Since(start).Seconds()

			// замечаем путь запроса и его метаданные
			endpoint := c.Path()
			method := c.Request().Method
			status := strconv.Itoa(c.Response().Status)

			// записываем длительность запроса и увеличиваем счетчик выполненных запросов
			m.RequestDuration.WithLabelValues(method, endpoint).Observe(duration)
			m.RequestsTotal.WithLabelValues(method, endpoint, status).Inc()

			return err
		}
	}
}

// запись операции
func (m *Metrics) RecordOperation(operationType string) {
	m.OperationsTotal.WithLabelValues(operationType).Inc()
}

// запись ошибки
func (m *Metrics) RecordError(errorType, handler string) {
	m.ErrorsTotal.WithLabelValues(errorType, handler).Inc()
}

// регистрация эндпоинта прометея
func (m *Metrics) RegisterMetricsHandler(e *echo.Echo, path string) {
	e.GET(path, func(c echo.Context) error {
		promhttp.Handler().ServeHTTP(c.Response(), c.Request())
		return nil
	})
}
