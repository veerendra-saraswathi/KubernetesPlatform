package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	port := ":8084"
	fmt.Println("Starting health server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

