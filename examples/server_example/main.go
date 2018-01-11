package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	cfg "github.com/apprentice3d/forge-api-go-client/config"
	"github.com/apprentice3d/forge-api-go-client/oauth"
)

type User struct {
	Profile *oauth.UserProfile
	Token   *oauth.Bearer
}

var (
	users          map[string]User
	clientId       string
	clientSecret   string
	config         cfg.Configuration
	threeLeggedApi oauth.ThreeLeggedApi
	infoApi        *oauth.InformationalApi
)

func init() {
	clientId = os.Getenv("FORGE_CLIENT_ID")
	clientSecret = os.Getenv("FORGE_CLIENT_SECRET")
	config := cfg.NewConfiguration()
	config.ClientId = clientId
	config.ClientSecret = clientSecret
	config.RedirectUri = "http://localhost:3000/cb"
	threeLeggedApi = oauth.ThreeLeggedApi{
		Configuration: config,
	}

	infoApi = oauth.NewInformationalApi()
	users = make(map[string]User)
}

func main() {

	shutdown := make(chan string)

	server := startBackgroundServer(threeLeggedApi.RedirectUri, shutdown)

	log.Println("Starting server: " + server.Addr)
	<-shutdown
	log.Println("Server is shutting down ...")
	defer server.Shutdown(context.Background())
}

func startBackgroundServer(serverAddress string, code chan string) *http.Server {
	srv := &http.Server{Addr: ":3000", Handler: http.DefaultServeMux}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/redirect", redirectHandler)
	http.HandleFunc("/cb", callbackHandler)

	http.HandleFunc("/shutdown", func(writer http.ResponseWriter, request *http.Request) {
		code <- "Shutdown server"
	})

	go srv.ListenAndServe()

	return srv
}

func registerUserSession(authCode string) (string, error) {
	token, err := threeLeggedApi.Gettoken(authCode)

	if err != nil {
		return "", nil
	}

	profile, err := infoApi.AboutMe(*token)
	if err != nil {
		return "", nil
	}

	users[profile.UserId] = User{
		profile,
		token,
	}

	return profile.UserId, nil
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {

	id, err := r.Cookie("session_id")
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", 302)
		return
	}
	profile := users[id.Value].Profile
	if profile == nil {
		log.Println("Could not find a profile with id: " + id.Value)
		http.SetCookie(w, &http.Cookie{
			Name:   "session_id",
			Value:  "expired",
			MaxAge: 0,
		})
		http.Redirect(w, r, "/", 302)
		return
	}
	log.Printf("User with email %s is trying to log in\n", profile.EmailId)

	//just a simple example of filtering access by some criteria - here: having an autodesk.com email
	if !strings.Contains(profile.EmailId, "@autodesk.com") {
		w.Write([]byte("Sorry, it seems that you are not an Autodesk employee"))
		return
	}

	fmt.Fprintf(w, "Hi %s %s [%s]\n", profile.FirstName, profile.LastName, profile.EmailId)
	r.Body.Close()
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	id, err := registerUserSession(strings.TrimPrefix(r.URL.RawQuery, "code="))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}
	cookie := http.Cookie{
		Name:  "session_id",
		Value: id,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/redirect", 302)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	redirectUrl, err := threeLeggedApi.Authorize("data:read", "")
	if err != nil {
		log.Println(err.Error())
	}
	http.Redirect(w, r, redirectUrl, 302)
}
