package metrics

import "github.com/labstack/echo/v4"

type MetricsInterface interface {
	EchoMiddleware() echo.MiddlewareFunc
	RecordOperation(operationType string)
	RecordError(errorType, handler string)
	RegisterMetricsHandler(e *echo.Echo, path string)
}
