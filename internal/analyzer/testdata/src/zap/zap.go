package zap

import (
	"go.uber.org/zap"
)

func TestZap() {
	logger := zap.NewExample()
	logger.Info("hello world")        // OK
	logger.Info("Hello world")        // want "uppercase"
	logger.Info("привет мир")         // want "non-Latin"
	logger.Info("hi 😊")               // want "emoji"
	logger.Info("my password is 123") // want "sensitive"
}
