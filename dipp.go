package dipp

import (
	"encoding/hex"
	"net/http"
	"strings"

	"golang.org/x/crypto/sha3"
)

func NewProofMiddleware(secret string, handler http.Handler) http.Handler {
	return &ProofMiddleware{next: handler, secret: secret}
}

func NewProofMiddlewareFunc(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return NewProofMiddleware(secret, next)
	}
}

// ProofMiddleware will hash the response and add a header. NOTE that
// this middleware will break Trailers added by calling code, after
// their WriteHeader() call.
type ProofMiddleware struct {
	next   http.Handler
	secret string
}

func (m *ProofMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Data-Integrity-Proof") == "true" && strings.ToLower(r.Header.Get("Connection")) != "upgrade" {
		proofWriter := &ProofWriter{
			ResponseWriter: w,
			digestWriter: sha3.NewShake256(),
		}

		_, _ = proofWriter.digestWriter.Write([]byte(m.secret))

		// HTTP Trailing headers MUST be declared BEFORE any writes (see https://golang.org/pkg/net/http/#example_ResponseWriter_trailers)
		w.Header().Set("Trailer", "X-Data-Integrity-Proof")

		m.next.ServeHTTP(proofWriter, r)

		out := make([]byte, 32)
		_, _ = proofWriter.digestWriter.Read(out)
		proof := hex.EncodeToString(out)
		w.Header().Set("X-Data-Integrity-Proof", proof)
	} else {
		m.next.ServeHTTP(w, r)
	}
}

type ProofWriter struct {
	http.ResponseWriter
	code int
	digestWriter sha3.ShakeHash
}

func (w *ProofWriter) Write(data []byte) (int, error) {
	w.digestWriter.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *ProofWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func HashMac(secret string, payload []byte) string {
	digest := sha3.NewShake256()

	_, _ = digest.Write([]byte(secret))
	_, _ = digest.Write(payload)

	out := make([]byte, 32)
	_, _ = digest.Read(out)

	return hex.EncodeToString(out)
}
