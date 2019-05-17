package dipp

import (
	"encoding/hex"
	"net/http"
	"strings"

	"golang.org/x/crypto/sha3"
)

func IntegrityProofMiddleware(secret string, handler http.Handler) http.Handler {
	return &ProofMiddleware{next: handler, secret: secret}
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
		writer := &ProofWriter{
			ResponseWriter: w,
		}

		m.next.ServeHTTP(writer, r)

		proof := HashMac(m.secret, writer.buffer)
		w.Header().Set("X-Data-Integrity-Proof", proof)

		if writer.code != 0 {
			w.WriteHeader(writer.code)
		}

		_, _ = writer.Write(writer.buffer)
	} else {
		m.next.ServeHTTP(w, r)
	}
}

type ProofWriter struct {
	http.ResponseWriter

	code int

	buffer []byte
}

func (w *ProofWriter) Write(data []byte) (int, error) {
	w.buffer = append(w.buffer, data...)
	return len(data), nil
}

func (w *ProofWriter) WriteHeader(code int) {
	w.code = code
}

func HashMac(secret string, payload []byte) string {
	digest := sha3.NewShake256()

	_, _ = digest.Write([]byte(secret))
	_, _ = digest.Write(payload)

	out := make([]byte, 32)
	_, _ = digest.Read(out)

	return hex.EncodeToString(out)
}
