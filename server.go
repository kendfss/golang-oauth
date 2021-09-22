package main

import (
    "encoding/json"
    "fmt"
    "io/fs"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
    
)

type OAuthAccessResponse struct {
    AccessToken string `json:"access_token"`
}
type APIKey struct {
    ClientID string `json:"clientID"`
    ClientSecret string `json:"clientSecret"`
}
type KeyChain struct {
    Github APIKey `json:"github"`
    Google APIKey `json:"google"`
}

var (
    // credentials map[string]map[string]string
    credentials KeyChain
    // clientID string
    // clientSecret string

    src = http.Dir("src")
    httpClient = http.Client{}
)


func getCredsFrom(path string) error {
    bites, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }
    err = json.Unmarshal(bites, &credentials)
    if e := saveJSONP(credentials, "src/creds.jsonp", "credata"); e != nil {
        log.Printf("Couldn't save creds.jsonp: %s", e)
    }
    return err
}
func pwd() string {
    path, err := os.Getwd()
    if err != nil {
        log.Printf("Couldn't find current working directory: %v\n", err)
    }
    return path
}
func saveJSONP(object interface{}, filename, varname string) error {
    data, err := json.Marshal(object)
    if err != nil {
        return err
    }
    data = append([]byte(fmt.Sprintf("%s = `", varname)), data...)
    data = append(data, []byte("`;")...)
    return ioutil.WriteFile(filename, data, fs.ModePerm)
}

func handleRequests(port int) {
    portNum := strconv.Itoa(port)
    log.Println("On:", "http://localhost:" + portNum)
    server := http.FileServer(src)
    http.Handle("/", server)

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
        reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", credentials.Github.ClientID, credentials.Github.ClientSecret, code)
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
        w.Header().Set("Location", "src/welcome.htm?access_token="+t.AccessToken)
        log.Printf("AccessToken:\n\t%q\n", t.AccessToken)
        w.WriteHeader(http.StatusFound)
    })

    log.Fatal(http.ListenAndServe(":" + portNum, nil))
}

func main() {
    if err:=getCredsFrom("creds.json"); err != nil {
        log.Fatalf("Couldn't parse credentials:\n\t%s\n", err)
    }
    
    // clientID = credentials["github"]["clientID"]
    // clientSecret = credentials["github"]["clientSecret"]

    log.Printf("Serving: %q\n", pwd())
    handleRequests(8000)
}
