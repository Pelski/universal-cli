package main

import (
	"fmt"
	"strings"
)

// HandleOperation routes the operation to the appropriate request function
func HandleOperation(operation string, resources []string, flags map[string]interface{}) {
	endpoint := buildEndpoint(resources)

	if debug {
		fmt.Println("> Endpoint:", endpoint)
	}

	switch operation {
	case "get", "list", "show":
		MakeRequest(GET, endpoint, flags)
	case "create", "search", "find":
		MakeRequest(POST, endpoint, flags)
	case "update", "set":
		MakeRequest(PUT, endpoint, flags)
	case "delete", "drop":
		MakeRequest(DELETE, endpoint, nil)
	default:
		fmt.Println("Unknown operation:", operation)
	}
}

// buildEndpoint constructs the API endpoint from resources
func buildEndpoint(resources []string) string {
	return "/" + strings.Join(resources, "/")
}
