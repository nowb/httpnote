package httpnote

import (
	"net/http"

	"github.com/rs/zerolog"
)

func NewHTTPNoteHandlerFunc(logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.
			Debug().
			Object("req", MapRequest(r)).
			Send()
	}
}
