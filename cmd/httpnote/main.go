package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nowb/httpnote"
	"github.com/rs/zerolog"
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

	zerolog.TimeFieldFormat = time.RFC3339Nano

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	handler := httpnote.NewHTTPNoteHandlerFunc(&logger, encodeBytes)

	logger.Fatal().Err(http.ListenAndServe(":"+port, handler)).Send()
}
