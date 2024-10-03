package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var debug bool

func initConfig(configPath string) {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("configuration")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading configuration: %s \n", err))
	}
	if debug {
		fmt.Println("> Configuration loaded from file:", viper.ConfigFileUsed())
	}
}

func main() {
	args := os.Args[1:]
	configPath := ""
	var newArgs []string

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "--config" {
			if i+1 < len(args) {
				configPath = args[i+1]
				i += 2
			} else {
				fmt.Println("Missing value for --config")
				os.Exit(1)
			}
		} else if strings.HasPrefix(arg, "--config=") {
			configPath = strings.TrimPrefix(arg, "--config=")
			i++
		} else if arg == "--ucli-debug" {
			debug = true
			i++
		} else if strings.HasPrefix(arg, "--ucli-debug=") {
			value := strings.TrimPrefix(arg, "--ucli-debug=")
			debug = value == "true" || value == "1"
			i++
		} else {
			newArgs = append(newArgs, arg)
			i++
		}
	}
	args = newArgs

	initConfig(configPath)

	if len(args) < 1 {
		fmt.Println("You need to provide an operation, e.g. get, create, update, delete")
		os.Exit(1)
	}

	operation := args[0]
	var resources []string
	var flagsArgs []string

	for i := 1; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flagsArgs = args[i:]
			break
		} else {
			resources = append(resources, arg)
		}
	}

	flags := parseDynamicFlags(flagsArgs)

	if debug {
		fmt.Println("> Operation:", operation)
		fmt.Println("> Resources:", resources)
		fmt.Println("> Flags:", flags)
	}

	handleOperation(operation, resources, flags)
}

func parseDynamicFlags(args []string) map[string]interface{} {
	flags := make(map[string]interface{})
	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			key := strings.TrimPrefix(arg, "--")
			var value string
			if strings.Contains(key, "=") {
				parts := strings.SplitN(key, "=", 2)
				key = parts[0]
				value = parts[1]
			} else {
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
					i++
					value = args[i]
				} else {
					value = ""
				}
			}
			flags[key] = parseValue(value)
		}
		i++
	}
	return flags
}

func parseValue(value string) interface{} {
	// Try to parse as int
	if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
		return intVal
	}
	// Try to parse as bool
	if boolVal, err := strconv.ParseBool(value); err == nil {
		return boolVal
	}
	// Try to parse as float
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}
	// Try to parse as JSON (array or object)
	trimmedValue := strings.TrimSpace(value)
	if strings.HasPrefix(trimmedValue, "[") || strings.HasPrefix(trimmedValue, "{") {
		var jsonData interface{}
		err := json.Unmarshal([]byte(trimmedValue), &jsonData)
		if err == nil {
			return jsonData
		}
	}
	// Return as string
	return value
}

func handleOperation(operation string, resources []string, flags map[string]interface{}) {
	endpoint := buildEndpoint(resources)

	if debug {
		fmt.Println("> Endpoint:", endpoint)
	}

	switch operation {
	case "get", "list", "show":
		makeGetRequest(endpoint, flags)
	case "create", "search", "find":
		makePostRequest(endpoint, flags)
	case "update", "set":
		makePutRequest(endpoint, flags)
	case "delete", "drop":
		makeDeleteRequest(endpoint)
	default:
		fmt.Println("Unknown operation:", operation)
	}
}

func buildEndpoint(resources []string) string {
	endpoint := "/" + strings.Join(resources, "/")
	return endpoint
}

func makeGetRequest(endpoint string, params map[string]interface{}) {
	url := viper.GetString("url") + endpoint
	if debug {
		fmt.Println("> Performing GET on URL:", url)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = q.Encode()

	setAuthorization(req)
	setCustomHeaders(req)

	if debug {
		fmt.Println("> Request headers:", req.Header)
		fmt.Println("> Query parameters:", req.URL.RawQuery)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer closeBody(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	logResponseDetails(resp, body)
}

func makePostRequest(endpoint string, bodyData map[string]interface{}) {
	url := viper.GetString("url") + endpoint
	if debug {
		fmt.Println("> Performing POST on URL:", url)
	}
	jsonData, _ := json.Marshal(bodyData)

	if debug {
		fmt.Println("> Data sent in request:", string(jsonData))
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	setAuthorization(req)
	setCustomHeaders(req)

	if debug {
		fmt.Println("> Request headers:", req.Header)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer closeBody(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	logResponseDetails(resp, body)
}

func makePutRequest(endpoint string, bodyData map[string]interface{}) {
	url := viper.GetString("url") + endpoint
	if debug {
		fmt.Println("> Performing PUT on URL:", url)
	}
	jsonData, _ := json.Marshal(bodyData)

	if debug {
		fmt.Println("> Data sent in request:", string(jsonData))
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	setAuthorization(req)
	setCustomHeaders(req)

	if debug {
		fmt.Println("> Request headers:", req.Header)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer closeBody(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	logResponseDetails(resp, body)
}

func makeDeleteRequest(endpoint string) {
	url := viper.GetString("url") + endpoint
	if debug {
		fmt.Println("> Performing DELETE on URL:", url)
	}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	setAuthorization(req)
	setCustomHeaders(req)

	if debug {
		fmt.Println("> Request headers:", req.Header)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer closeBody(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	logResponseDetails(resp, body)
}

func logResponseDetails(resp *http.Response, body []byte) {
	if debug {
		fmt.Println("> Response Status:", resp.Status)
		fmt.Println("> Response Headers:", resp.Header)
	}

	if body == nil || len(body) == 0 {
		fmt.Println(interpretStatus(resp.Status))
	} else {
		fmt.Println(string(body))
	}
}

func setAuthorization(req *http.Request) {
	tokenPath := viper.GetString("token")
	if tokenPath != "" {
		tokenBytes, err := ioutil.ReadFile(tokenPath)
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

func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func interpretStatus(status string) string {
	switch status {
	case "200 OK":
		return "[200] Operation completed successfully."
	case "201 Created":
		return "[201] Resource created successfully."
	case "202 Accepted":
		return "[202] Request accepted, processing in progress."
	case "204 No Content":
		return "[204] Operation completed successfully, no content to display."
	case "301 Moved Permanently":
		return "[301] The resource has been moved permanently to a new URL."
	case "302 Found":
		return "[302] The resource is temporarily located at a different URL."
	case "304 Not Modified":
		return "[304] The resource has not been modified since the last request."
	case "400 Bad Request":
		return "[400] The request was invalid. Please check the input data."
	case "401 Unauthorized":
		return "[401] Authorization required or token is invalid."
	case "403 Forbidden":
		return "[403] You do not have permission to access this resource."
	case "404 Not Found":
		return "[404] The resource could not be found."
	case "405 Method Not Allowed":
		return "[405] The HTTP method used is not allowed for this resource."
	case "408 Request Timeout":
		return "[408] The server timed out waiting for the request."
	case "409 Conflict":
		return "[409] There was a conflict with the request, such as duplicate data."
	case "410 Gone":
		return "[410] The resource is no longer available."
	case "413 Payload Too Large":
		return "[413] The request payload is too large to be processed."
	case "415 Unsupported Media Type":
		return "[415] The media type of the request is not supported."
	case "429 Too Many Requests":
		return "[429] Too many requests have been made in a short period. Please try again later."
	case "500 Internal Server Error":
		return "[500] The server encountered an error."
	case "502 Bad Gateway":
		return "[502] The server received an invalid response from an upstream server."
	case "503 Service Unavailable":
		return "[503] The server is temporarily unavailable. Please try again later."
	case "504 Gateway Timeout":
		return "[504] The server did not receive a timely response from an upstream server."
	default:
		return fmt.Sprintf("[%s] Unknown status.", status)
	}
}
