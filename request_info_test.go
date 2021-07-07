package httpnote_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/nowb/httpnote"
)

func TestMapRequest(t *testing.T) {
	tests := []struct {
		r              *http.Request
		expectedOutput *httpnote.RequestInfo
	}{
		{
			r: httptest.NewRequest(http.MethodGet, "/test", nil),
			expectedOutput: &httpnote.RequestInfo{
				Method:     http.MethodGet,
				URL:        &httpnote.URLInfo{Path: "/test"},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Host:       "example.com",
				RemoteAddr: "192.0.2.1:1234",
				RequestURI: "/test",
			},
		},
	}

	for _, test := range tests {
		output := httpnote.MapRequest(test.r, true)

		if !reflect.DeepEqual(test.expectedOutput, output) {
			t.Errorf("expected: %v, got: %v", test.expectedOutput, output)
		}
	}
}
