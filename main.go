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
		fmt.Println("Configuration loaded from file:", viper.ConfigFileUsed())
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
		fmt.Println("Operation:", operation)
		fmt.Println("Resources:", resources)
		fmt.Println("Flags:", flags)
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
		fmt.Println("Endpoint:", endpoint)
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
		fmt.Println("Performing GET on URL:", url)
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
		fmt.Println("Request headers:", req.Header)
		fmt.Println("Query parameters:", req.URL.RawQuery)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer closeBody(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func makePostRequest(endpoint string, bodyData map[string]interface{}) {
	url := viper.GetString("url") + endpoint
	if debug {
		fmt.Println("Performing POST on URL:", url)
	}
	jsonData, _ := json.Marshal(bodyData)

	if debug {
		fmt.Println("Data sent in request:", string(jsonData))
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	setAuthorization(req)
	setCustomHeaders(req)

	if debug {
		fmt.Println("Request headers:", req.Header)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer closeBody(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func makePutRequest(endpoint string, bodyData map[string]interface{}) {
	url := viper.GetString("url") + endpoint
	if debug {
		fmt.Println("Performing PUT on URL:", url)
	}
	jsonData, _ := json.Marshal(bodyData)

	if debug {
		fmt.Println("Data sent in request:", string(jsonData))
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	setAuthorization(req)
	setCustomHeaders(req)

	if debug {
		fmt.Println("Request headers:", req.Header)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer closeBody(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func makeDeleteRequest(endpoint string) {
	url := viper.GetString("url") + endpoint
	if debug {
		fmt.Println("Performing DELETE on URL:", url)
	}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	setAuthorization(req)
	setCustomHeaders(req)

	if debug {
		fmt.Println("Request headers:", req.Header)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer closeBody(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
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
			fmt.Println("Set Bearer Token in Authorization header")
		}
		return
	}

	username := viper.GetString("username")
	password := viper.GetString("password")
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
		if debug {
			fmt.Println("Set Basic Auth with username:", username)
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
			fmt.Println("Set additional headers from configuration:", headers)
		}
	}
}

func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Fatal(err)
	}
}
