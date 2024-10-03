package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// LogResponseDetails logs details about the HTTP response
func LogResponseDetails(resp *http.Response, body []byte) {
	if debug {
		fmt.Println("> Response Status:", resp.Status)
		fmt.Println("> Response Headers:", resp.Header)
	}

	if len(body) == 0 {
		fmt.Println(InterpretStatus(resp.Status))
	} else {
		fmt.Println(string(body))
	}
}

// CloseBody closes the response body
func CloseBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// InterpretStatus provides a user-friendly message for HTTP status codes
func InterpretStatus(status string) string {
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
