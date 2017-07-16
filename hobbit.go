// Extremely simple console tool. Takes Long URL as a first argument.
// Returns Bitly short URL to console. That's It.
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
    API_BASEURL     = "https://api-ssl.bitly.com"
    API_AUTH        = "/oauth/access_token"
    API_SHORT       = "/v3/shorten"
    TOKEN_FILE      = ".hobbit_token"

    API_TOKEN_ERROR_INVALID = "INVALID_ARG_ACCESS_TOKEN"
    API_TOKEN_ERROR_MISSING = "MISSING_ARG_ACCESS_TOKEN"
)


// Create config file and write token to It
func SaveConfig(filename string, tk string) error {
    HOME := os.Getenv("HOME")
    var path string = ""

    // If HOME is empty
    if len(HOME) == 0 {
      path = filename
      config_file, err := os.Create(path)
      defer config_file.Close()
      if err != nil {
        log.Println(err)
        return err
      }

      _, err = config_file.WriteString(tk)
      if err != nil {
        log.Println(err)
        return err
      }

      return nil
  }

    // Set Path if HOME is not empty
    path = HOME + "/" + filename

    // Create config file and write token to It
    config_file, err := os.Create(path)
    defer config_file.Close()

    if err != nil {
      log.Println(err)
      log.Fatal("Can't create config file: ", path, " Check free space and permissions")
    }

    _, err = config_file.WriteString(tk)
    if err != nil {
      log.Println(err)
      return err
    }

    return nil
}


// Get Token from file
func ReadToken(filename string) string {
    HOME := os.Getenv("HOME")
    var path string = ""

    if len(HOME) == 0 {
      path = filename
    } else {
      path = HOME + "/" + filename
    }

    f, err := os.Open(path)
    defer f.Close()

    if err != nil {
      log.Println(err)
      return ""
    }

    stat, _ := f.Stat()
    bs := make([]byte, stat.Size())
    f.Read(bs)

    return string(bs)
}


// Authenticate. Returns access_token
func Auth() string {
    username := ""
    password := ""

    // Get username and password from user
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
    var shortUrl, longUrl string
    token := ReadToken(TOKEN_FILE)

    // Check argumets
    if len(os.Args) < 2 {
      fmt.Println("USAGE:", os.Args[0], "http://example.com")
      fmt.Println("NOTE:  'http/https' part is obligatory")
      os.Exit(1)
    }

    longUrl = os.Args[1]

    // Authenticate and save token. 3 attempts
    for i := 0; i < 3; i++ {
      shortUrl = Shorten(token, longUrl)
      if shortUrl == API_TOKEN_ERROR_INVALID || shortUrl == API_TOKEN_ERROR_MISSING {
        log.Println("Bitly access_token is not set/valid")
        token = Auth()
      } else {
        fmt.Printf(shortUrl)
        break
      }
    }

    SaveConfig(TOKEN_FILE, token)
}
