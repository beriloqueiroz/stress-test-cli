package main

import (
	"math/rand/v2"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		test := rand.IntN(100)
		// time.Sleep(time.Millisecond * 10)
		if test <= 20 {
			w.WriteHeader(200)
		} else if test <= 40 {
			w.WriteHeader(400)
		} else if test <= 60 {
			w.WriteHeader(404)
		} else if test <= 80 {
			w.WriteHeader(429)
		} else {
			w.WriteHeader(300)
		}
	})
	http.ListenAndServe(":8080", mux)
}
