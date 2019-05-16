package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// IsArn checks if a string starts with `arn:aws:`.
func IsArn(possibleArn string) bool {
	arnPrefix := "arn:aws:"
	return strings.HasPrefix(possibleArn, arnPrefix)
}

// GetPublicIP sends a request to https://www.ipify.org/ for the end user's local ip address in text format.
func GetPublicIP() string {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		log.Fatalf("unable to get IP address: %v\n", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to read response: %v\n", err)
	}

	return string(body)
}
