package src

import (
    "net/http"
)

// SafPage endpoint to rm request
func SafPage(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Thanks for hitting saf")
}