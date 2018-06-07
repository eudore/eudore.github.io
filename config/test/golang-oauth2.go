package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

const htmlIndex = `<html><body>
<p><a href="/login">Log in with Google</a></p>
<p><a href="/loginGithub">log in with Github</a></p> 
</body></html>
`

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://47.52.173.119:8085/callback/google",
		ClientID:     os.Getenv("google_client_id"),
		ClientSecret: os.Getenv("google_client_secret"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
)

var (
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://47.52.173.119:8085/callback/github",
		ClientID:     "46eee6bd2e287032ad0c",
		ClientSecret: "b1d5f86930568775d4be31e0329d02fd4c751ec3",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
)

var oauthStateString string

func init() {
	oauthStateString = getRandomString()
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/loginGithub", handleGitlogin)
	http.HandleFunc("/callback/google", handleGoogleCallback)
	http.HandleFunc("/callback/github", handleGithubCallback)
	fmt.Println(http.ListenAndServe(":8085", nil))
}

func getRandomString() string {
	letters := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY")
	result := make([]rune, 16)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGitlogin(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL(oauthStateString)
	fmt.Println("github redirect:",url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "Content: %s\n", contents)
}

func handleGithubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	response, _ := http.PostForm("https://github.com/login/oauth/access_token", url.Values{
		"client_id":     {githubOauthConfig.ClientID},
		"client_secret": {githubOauthConfig.ClientSecret},
		"code":          {code},
	})
	defer response.Body.Close()

	contents, _ := ioutil.ReadAll(response.Body)
	res, _ := http.Get("https://api.github.com/user?" + string(contents))
	defer res.Body.Close()
	con, _ := ioutil.ReadAll(res.Body)
	fmt.Fprintln(w, string(con))
}