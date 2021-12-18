package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte(req.Host))
	})
	log.Fatal(http.ListenAndServe(":3000", nil))
}
