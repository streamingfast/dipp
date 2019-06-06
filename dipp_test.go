package dipp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func RandomHTTPTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "This is a random GET payload")
	}
	return http.HandlerFunc(fn)
}

func TestAuth(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		description  string
		method string
		url          string
		expectedBody string
		expectedCode int
		expectedDIPP bool
		expectedDIPPValue string
	}{
		{
			description:  "DIPP SHALL be found in headers",
			method: 	  "GET",
			url:          "http://example.com/payload",
			expectedBody: "This is a random GET payload",
			expectedCode: 200,
			expectedDIPP: true,
			expectedDIPPValue: "fd212e1422606f966fd1373fc8071f70838b1fb043499b105ee24531d29baedf",
		},
		{
			description:  "DIPP SHALL not be found in headers",
			method: 	  "GET",
			url:          "http://example.com/payload",
			expectedBody: "This is a random GET payload",
			expectedCode: 200,
			expectedDIPP: false,
		},
	}


	dippMiddlewareHandler := NewProofMiddleware("supersecretsstuff", RandomHTTPTestHandler())

	for _, tc := range tests {
		req := httptest.NewRequest(tc.method, tc.url, nil)
		if tc.expectedDIPP {
			req.Header = map[string][]string{
				"X-Data-Integrity-Proof": {"true"},
			}
		}

		recorder := httptest.NewRecorder()
		dippMiddlewareHandler.ServeHTTP(recorder, req)
		res := recorder.Result()
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		assert.Equal(tc.expectedCode, res.StatusCode, tc.description)
		assert.Equal(tc.expectedBody, string(b), tc.description)

		if tc.expectedDIPP {
			assert.Equal(tc.expectedDIPPValue, res.Trailer.Get("X-Data-Integrity-Proof"), tc.description)
		} else {
			assert.Empty(res.Trailer.Get("X-Data-Integrity-Proof"), tc.description)
		}
	}
}
