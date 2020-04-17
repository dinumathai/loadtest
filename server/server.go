package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	errorCode := getErrorCode(r)
	w.WriteHeader(errorCode)
	fmt.Fprintf(w, "Returning code %d", errorCode)
}

func getErrorCode(r *http.Request) int {
	errorCodeInputs, ok := r.URL.Query()["error-code"]
	if !ok || len(errorCodeInputs[0]) < 1 {
		return 200
	}

	errorCodeStr := errorCodeInputs[0]
	errorCode, err := strconv.ParseUint(errorCodeStr, 0, 64)
	if errorCode < 200 || err != nil {
		return 200
	}
	return int(errorCode)
}
