package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nowb/httpnote"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	port := "8080"

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	var encodeBytes bool
	if eb := os.Getenv("ENCODE_BYTES"); eb != "" {
		if ebBool, err := strconv.ParseBool(eb); err == nil {
			encodeBytes = ebBool
		}
	}

	logFilename := "./logs/httpnote-access.log"
	if lf := os.Getenv("LOG_FILENAME"); lf != "" {
		logFilename = lf
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	w := &lumberjack.Logger{
		Filename: logFilename,
		MaxSize:  20,
		Compress: true,
	}
	mw := io.MultiWriter(os.Stderr, w)
	logger := zerolog.New(mw).With().Timestamp().Logger()
	handler := httpnote.NewHTTPNoteHandlerFunc(&logger, encodeBytes)

	logger.Fatal().Err(http.ListenAndServe(":"+port, handler)).Send()
}
