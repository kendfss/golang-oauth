package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    
)

var (
    credentials map[string]string
    clientID string
    clientSecret string

    src = http.Dir("src")
    counter int
    mutex = &sync.Mutex{}
    httpClient = http.Client{}
)

type OAuthAccessResponse struct {
    AccessToken string `json:"access_token"`
}
func getCredsFrom(path string) error {
    bites, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }
    err = json.Unmarshal(bites, &credentials)
    return err
}
func pwd() string {
    path, err := os.Getwd()
    if err != nil {
        log.Printf("Couldn't find current working directory: %v\n", err)
    }
    return path
}
func read(path string) (string, error) {
    lines, err := readLines(path)
    if err != nil {
        log.Printf("Couldn't read %q\n\t%s\n", path, err)
        return "", err
    }
    str := ""
    for _, line := range lines {
        str += line
    }
    return str, nil
}
func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range lines {
        fmt.Fprintln(w, line)
    }
    return w.Flush()
}

func homePage(w http.ResponseWriter, r *http.Request) {
    path := filepath.Join("src", "index.htm")
    // path := filepath.Join("home.htm")
    text, err := read(path)
    if err != nil {
        log.Printf("Couldn't read %q\n", "home.htm")
    }
    fmt.Fprintf(w, text)
    fmt.Println("Endpoint Hit: homePage")
}


func incrementCounter(w http.ResponseWriter, r *http.Request) {
    mutex.Lock()
    counter++
    fmt.Fprintf(w, strconv.Itoa(counter))
    mutex.Unlock()
}

func handleRequests(port int) {
    portNum := strconv.Itoa(port)
    log.Println("On:", "http://localhost:" + portNum)
    fs := http.FileServer(src)
    http.Handle("/", fs)

        

    // Create a new redirect route route
    http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
        // First, we need to get the value of the `code` query param
        err := r.ParseForm()
        if err != nil {
            fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
            w.WriteHeader(http.StatusBadRequest)
        }
        code := r.FormValue("code")
        log.Printf("code:\n\t%q\n", code)


        // Next, lets for the HTTP request to call the github oauth enpoint
        // to get our access token
        reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
        log.Printf("Request url:\n\t%q\n", reqURL)
        req, err := http.NewRequest(http.MethodPost, reqURL, nil)
        if err != nil {
            fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
            w.WriteHeader(http.StatusBadRequest)
        }
        // We set this header since we want the response
        // as JSON
        req.Header.Set("accept", "application/json")

        // Send out the HTTP request
        res, err := httpClient.Do(req)
        if err != nil {
            fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
            w.WriteHeader(http.StatusInternalServerError)
        }
        defer res.Body.Close()

        // Parse the request body into the `OAuthAccessResponse` struct
        var t OAuthAccessResponse
        if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
            fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
            w.WriteHeader(http.StatusBadRequest)
        }

        // Finally, send a response to redirect the user to the "welcome" page
        // with the access token
        w.Header().Set("Location", "/welcome.html?access_token="+t.AccessToken)
        log.Printf("AccessToken:\n\t%q\n", AccessToken)
        w.WriteHeader(http.StatusFound)
    })

    log.Fatal(http.ListenAndServe(":" + portNum, nil))
}

func main() {
    if err:=getCredsFrom("creds.json"); err != nil {
        log.Fatalf("Couldn't parse credentials:\n\t%s\n", err)
    }
    
    clientID = credentials["clientID"]
    clientSecret = credentials["clientSecret"]

    log.Printf("Serving: %q\n", pwd())
    handleRequests(8000)
}
