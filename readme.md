golang-oauth
---
  
  
This repo contains an Oath2 authenticator in pure go, html and js.  
It's an extension of a tutorial by Soham Kamani<sup>[1](#references)</sup> which also implements support for authentication via google, and use of json/p to manage credentials. 
Credentials aren't secured, yet. Do not use this in production.  
  
  
### Usage
You must have [go](https://golang.org/) to use this. It was developed in Go 1.17.
Github
1. [Register your new application on Github](https://github.com/settings/applications/new). 
    - In the "callback URL" field, enter `http://localhost:8080/oauth/redirect`
    - Request a client secret
1. Save your credentials in a file called `golang-oauth/golang-creds.json` using the following syntax:
    ```json
    {
        "github": {
            "clientID": "",
            "clientSecret": "",
        },
        "google": {
            "clientID": "",
            "clientSecret": "",
        }
    }
    ```
1. Start the server by executing `go run server.go`
1. Navigate to `http://localhost:8000` on your browser.
  
  

### References
1. [Implementing OAuth 2.0 with Go(Golang)](https://www.sohamkamani.com/golang/oauth/)
