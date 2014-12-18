package main

import (
	// Stdlib
	"net/http"

	// Vendor
	"github.com/codegangsta/negroni"
)

func main() {
	var (
		token = mustGetenv("WEBHOOK_TOKEN")
	)

	// Prepare the mux.
	mux := http.NewServeMux()
	mux.HandleFunc("/review-request-published", onReviewRequestPublished)
	mux.HandleFunc("/review-published", onReviewPublished)
	mux.HandleFunc("/review-request-closed", onReviewRequestClosed)

	// Setup and start Negroni Classic.
	n := negroni.Classic()
	n.Use(requireAccessToken("webhookToken", token))
	n.UseHandler(mux)
	n.Run(os.Getenv("PORT"))
}

func requireAccessToken(param, value string) negroni.Handler {
	handler := func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.FormValue(param) != value {
			http.Error(rw, "Invalid Access Token", http.StatusForbidden)
			return
		}
		next(rw, r)
	}
	return negroni.HandlerFunc(handler)
}

func mustGetenv(key string) (value string) {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Errorf("environment variable not set: %v", key))
	}
	return v
}
