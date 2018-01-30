// Extremely simple console tool. Takes Long URL as a first argument.
// Returns Bitly short URL to console. That's It.
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/howeyc/gopass"
)

const (
	apiBaseURL = "https://api-ssl.bitly.com"
	apiAuth    = "/oauth/access_token"
	apiShorten = "/v3/shorten"
)

// Auth performs authentication. Returns access_token
func Auth() string {
	username := ""

	// Get username and password from user
	fmt.Printf("%s: ", "username")
	fmt.Scanln(&username)
	fmt.Printf("%s: ", "password")
	password, err := gopass.GetPasswdMasked()
	if err != nil {
		log.Fatal(err)
	}

	// Create request
	client := &http.Client{}
	urlStr := apiBaseURL + apiAuth
	r, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set auth header
	msg := username + ":" + string(password)
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(msg))
	r.Header.Add("Authorization", authHeader)

	// Get response
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Get token
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Check errors and return access_token
	responseString := string(body)
	_, err = regexp.MatchString("[a-zA-Z0-9]+", responseString)
	if err != nil {
		log.Fatal(responseString)
	}
	return responseString
}

// Shorten long URL
func Shorten(tk string, longurl string) string {
	urlStr := apiBaseURL + apiShorten

	URL, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}

	// Set parameters
	parameters := url.Values{}
	parameters.Add("access_token", tk)
	parameters.Add("longUrl", longurl)
	parameters.Add("format", "txt")
	URL.RawQuery = parameters.Encode()

	// Call API endpoint
	resp, err := http.Get(URL.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func main() {
	var shortURL, longURL string

	token, ok := os.LookupEnv("BITLY_TOKEN")
	if !ok {
		fmt.Println("BITLY_TOKEN environment variable in NOT set. Please provide your Bitly username and pass:")
		token = Auth()
		fmt.Println("Add this to your .bashrc or .zshrc:")
		fmt.Println("export BITLY_TOKEN=" + token)
		os.Exit(1)
	}

	// Check argumets
	if len(os.Args) < 2 {
		fmt.Println("USAGE:", os.Args[0], "http://example.com")
		fmt.Println("NOTE:  'http/https' part is obligatory")
		os.Exit(1)
	}

	longURL = os.Args[1]
	shortURL = Shorten(token, longURL)
	fmt.Printf(shortURL)
}
