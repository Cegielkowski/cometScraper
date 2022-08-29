package middleware

import "cometScraper/utils/logger"

// Middleware ...
type Middleware struct {
	logger logger.Logger
}

// NewMiddleware will create new a Middleware object
func NewMiddleware(logger logger.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}
