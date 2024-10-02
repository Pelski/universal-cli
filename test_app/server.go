package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL:", r.URL.String())
	fmt.Println("Method:", r.Method)
	fmt.Println("Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}

	if r.Body != nil {
		body, err := io.ReadAll(r.Body)
		if err == nil && len(body) > 0 {
			fmt.Println("Body:")
			fmt.Println(string(body))
		} else {
			fmt.Println("No body received or error reading body.")
		}
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fprintf, err := fmt.Fprintf(w, "Response from server (%s)", r.Method)
	if err != nil {
		log.Println("Error writing response:", err)
	} else {
		log.Println("Response:", fprintf)
	}
	fmt.Println("\n--------------------")
	fmt.Println("")
}

func main() {
	http.HandleFunc("/", requestHandler)

	port := ":8123"
	fmt.Printf("Serving on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
