package main

import (
	"io"
	"net/http"
)

// func return response
func CheckoutDummy(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("id")
	switch key {
	case "42":
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 200, "balance": 100500, "err":""}`)
	case "100500":
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 400, "balance": 0, "err":"bad_balance"}`)
	case "__broken_json":
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `"status": 400`)
	case "__internal_error":
		fallthrough
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
