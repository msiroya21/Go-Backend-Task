package middleware

import (
	"time"

	"go-backend-task/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("X-Request-ID", requestID)

		c.Locals("requestid", requestID)

		return c.Next()
	}
}

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		requestID := c.Locals("requestid")
		if requestID == nil {
			requestID = ""
		}

		logger.Logger.Info("Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("request_id", requestID.(string)),
			zap.Duration("duration", duration),
		)

		return err
	}
}
