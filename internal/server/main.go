package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	log.Printf("port=%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	if err != nil {
		log.Fatal(err)
	}
}
