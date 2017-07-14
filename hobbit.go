// Extremely simple console tool.
// Takes Long URL as a first argument.
// Returns Bitly short URL to console.
// That's It.
package main

import (
      "os"
      "log"
      "fmt"
      "regexp"
      "io/ioutil"
      "net/http"
      "net/url"
      "encoding/base64"
)

const(
      API_BASEURL = "https://api-ssl.bitly.com"
      TOKEN_DIR   = ".hobbit"
      TOKEN_FILE  = "token"

      API_AUTH    = "/oauth/access_token"
      API_SHORT   = "/v3/shorten"
)

var(
    HOME = os.Getenv("HOME")
)


// Config file will be created under given path with given filename
// Returns *os.File object and error
func InitConfig(path string, filename string) (*os.File, error) {
  if len(HOME) == 0 {
    log.Println("HOME environment variable is empty. I'll use current folder")
    config_file, err := os.Create(TOKEN_FILE)
    if err != nil {
      return nil, err
    } else {
      return config_file, nil
    }
  }

  // Define absolute path to config
  config_dir       := HOME + "/" + TOKEN_DIR
  config_file_path := config_dir + "/" + TOKEN_FILE

  // Create config directory if not exists
  if _, err := os.Stat(config_dir); os.IsNotExist(err) {
    err := os.Mkdir(config_dir, 0755)
    if err != nil {
      log.Println(err)
      log.Fatal("can't create directory: ", config_dir ," Check free space and permissions")
    }
  }

  // Create config file if not exists
  if _, err := os.Stat(config_file_path); os.IsNotExist(err) {
    config_file, err := os.Create(config_file_path)
    if err != nil {
      log.Println(err)
      log.Fatal("can't create config file: ", config_file_path, " Check free space and permissions")
    }
    return config_file, err
  }

  return os.Open(config_file_path)
}

// Save Token on a filesystem
func SaveToken(tk string, f *os.File) (int, error) {
  return f.WriteString(tk)
}


// Get access_token
func Auth() string {
  username := ""
  password := ""

  // Get username and password from user
  fmt.Println("Looks like access_token is not set. Let's do It!")
  fmt.Printf("%s: ","username")
  fmt.Scanln(&username)
  fmt.Printf("%s: ", "password")
  fmt.Scanln(&password)

  // Create request
  client := &http.Client{}
  urlStr := API_BASEURL + API_AUTH
  r, err := http.NewRequest("POST", urlStr, nil)
  if err != nil {
    log.Fatal(err)
  }
  // Set auth header
  msg := username + ":" + password
  auth_header := "Basic " + base64.StdEncoding.EncodeToString([]byte(msg))
  r.Header.Add("Authorization", auth_header)

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
  response_string := string(body)
  _, err = regexp.MatchString("[a-zA-Z0-9]+", response_string)
  if err != nil {
    log.Fatal(response_string)
  }
  return response_string
}


// Shorten long URL
func Shorten(tk string ,longurl string) string {
  urlStr := API_BASEURL + API_SHORT

  Url, err := url.Parse(urlStr)
  if err != nil {
    log.Fatal(err)
  }

  // Set parameters
  parameters := url.Values{}
  parameters.Add("access_token", tk)
  parameters.Add("longUrl", longurl)
  parameters.Add("format", "txt")
  Url.RawQuery = parameters.Encode()

  // Call API endpoint
  resp, err := http.Get(Url.String())
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
  var token, longUrl string

  f, _ := InitConfig(TOKEN_DIR, TOKEN_FILE)
  defer f.Close()

  // Read existing token from file or get a new one
  stat, _ := f.Stat()
  if stat.Size() > 0 {
    bs := make([]byte, stat.Size())
    f.Read(bs)
    token = string(bs)
  } else {
    token = Auth()
    SaveToken(token, f)
  }

  // Check if argument has been passed to the script
  if len(os.Args) < 2 {
    fmt.Println("USAGE:", os.Args[0], "http://example.com")
    fmt.Println("NOTE:  'http/https' part is obligatory")
    os.Exit(1)
  }

  longUrl = os.Args[1]
  fmt.Printf(Shorten(token, longUrl))
}
