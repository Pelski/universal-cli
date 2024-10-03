package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// HTTPMethod defines a type for HTTP methods
type HTTPMethod string

// Constants representing HTTP methods
const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
)

// MakeRequest sends an HTTP request based on the method and parameters
func MakeRequest(method HTTPMethod, endpoint string, data map[string]interface{}) {
	url := viper.GetString("url") + endpoint
	if debug {
		fmt.Printf("> Performing %s on URL: %s\n", method, url)
	}

	var req *http.Request
	var err error

	if method == GET || method == DELETE {
		req, err = http.NewRequest(string(method), url, nil)
		if err != nil {
			log.Fatal(err)
		}
		if method == GET && data != nil {
			q := req.URL.Query()
			for key, value := range data {
				q.Add(key, fmt.Sprintf("%v", value))
			}
			req.URL.RawQuery = q.Encode()
		}
	} else {
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatal("Error marshaling JSON:", err)
		}
		if debug {
			fmt.Println("> Data sent in request:", string(jsonData))
		}
		req, err = http.NewRequest(string(method), url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
	}

	setAuthorization(req)
	setCustomHeaders(req)

	if debug {
		fmt.Println("> Request headers:", req.Header)
		if method == GET {
			fmt.Println("> Query parameters:", req.URL.RawQuery)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer CloseBody(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	LogResponseDetails(resp, body)
}

// setAuthorization sets the Authorization header based on config
func setAuthorization(req *http.Request) {
	tokenPath := viper.GetString("token")
	if tokenPath != "" {
		tokenBytes, err := os.ReadFile(tokenPath)
		if err != nil {
			log.Fatal("Error reading token file:", err)
		}
		token := strings.TrimSpace(string(tokenBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		if debug {
			fmt.Println("> Set Bearer Token in Authorization header")
		}
		return
	}

	username := viper.GetString("username")
	password := viper.GetString("password")
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
		if debug {
			fmt.Println("> Set Basic Auth with username:", username)
		}
	}
}

// setCustomHeaders adds additional headers from configuration
func setCustomHeaders(req *http.Request) {
	if viper.IsSet("headers") {
		headers := viper.GetStringMapString("headers")
		for key, value := range headers {
			req.Header.Set(key, value)
		}
		if debug {
			fmt.Println("> Set additional headers from configuration:", headers)
		}
	}
}
