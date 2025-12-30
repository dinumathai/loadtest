package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server at port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	errorCode := getErrorCode(r)
	w.WriteHeader(errorCode)
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
