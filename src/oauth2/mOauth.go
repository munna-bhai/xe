package main

import (
	//"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubClientId     = "1da3e3e57166dcfd116a"
	githubClientSecret = "806d9f4a49de69272bfaf24e5b4eb9afdebed5d9"
)

var (
	githubOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:9090/Callback",
		ClientID:     githubClientId,
		ClientSecret: githubClientSecret,
		Scopes:       []string{"public_repo", "user:email"},
		Endpoint:     github.Endpoint,
	}

	oauthStateString = "random"
)

func handleMain(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "/home.html")

	//u := "http://localhost:9090/home.html"
	//u := config.AuthCodeURL()
	//http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOAuthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")

	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := githubOAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//fmt.Println(token.AccessToken)
	response, err := http.Get("https://api.github.com/user?access_token=" + url.QueryEscape(token.AccessToken))
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(*token)
	fmt.Fprintf(w, "Content: %s\n", contents)

	// token, err := config.Exchange(context.Background(), code)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// e := json.NewEncoder(w)
	// e.SetIndent("", "  ")
	// e.Encode(*token)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleMain)
	r.HandleFunc("/login", handleLogin)
	r.HandleFunc("/Callback", handleCallback)
	http.Handle("/", r)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}