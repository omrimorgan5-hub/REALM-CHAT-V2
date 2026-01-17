package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 1. Define a route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Backend is running!")
	})

	// 2. Start the server
	fmt.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}