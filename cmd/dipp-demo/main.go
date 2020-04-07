package main

import (
	"github.com/dfuse-io/dipp"
	"net/http"
)

func someHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("some payload"))
}

func main() {
	http.Handle("/", dipp.NewProofMiddleware("super-secret", http.HandlerFunc(someHandler)))
	http.ListenAndServe("127.0.0.1:3000", nil)
}
