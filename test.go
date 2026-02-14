package main

import (
	"log/slog"
)

func main() {
	slog.Info("Hello world")           // должно ругаться на uppercase
	slog.Info("Привет")                // должно ругаться на non-Latin и uppercase
	slog.Info("hi 😊")                  // должно ругаться на emoji
	slog.Info("my password is secret") // sensitive
}
