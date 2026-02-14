package slog

import "log/slog"

func Test() {
	slog.Info("hello world")        // OK
	slog.Info("Hello world")        // want "uppercase"
	slog.Info("привет")             // want "non-Latin"
	slog.Info("hi 😊")               // want "emoji"
	slog.Info("my password is 123") // want "sensitive"
}
