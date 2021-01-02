package utils

import (
	"go.uber.org/zap"
)

// Lightweight logger.
var logger, _ = zap.NewDevelopment()

// Sugar is an API for logging from Logger.
var Sugar = logger.Sugar()
