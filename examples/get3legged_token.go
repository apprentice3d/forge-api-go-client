package main

import (
	"context"
	"fmt"
	cfg "github.com/apprentice3d/forge-api-go-client/config"
	"github.com/apprentice3d/forge-api-go-client/oauth"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	clientId := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	//config := &cfg.Configuration{
	//	ClientId:clientId,
	//	ClientSecret:clientSecret,
	//	RedirectUri:"http://localhost:3000/cb",
	//	BasePath:"https://developer.api.autodesk.com",
	//}

	config := cfg.NewConfiguration()
	config.ClientId = clientId
	config.ClientSecret = clientSecret
	config.RedirectUri = "http://localhost:3000/cb"

	client := &oauth.ThreeLeggedApi{
		Configuration: config,
	}

	link, err := client.Authorize("data:read", "")

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Please paste the following link into your browser: " + link)

	getCode := make(chan string)

	server := startBackgroundServer(client.RedirectUri, getCode)

	authCode := <-getCode

	defer server.Shutdown(context.Background())

	fmt.Printf("Received the auth code: %s\n", authCode)

	token, err := client.Gettoken(authCode)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Received following token:\ntoken type = %s,\nexpires in = %d,"+
		"\ntoken itself = %s,\nrefresh token = %s\n",
		token.TokenType, token.ExpiresIn, token.AccessToken, token.RefreshToken)

	info := oauth.NewInformationalApi()

	profile, err := info.AboutMe(*token)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Profile: %+v\n", profile)

	if new_token, err := client.RefreshToken(token.RefreshToken,
		"data:read"); err == nil {
		fmt.Printf("Received following information:\ntoken type = %s,\n expires in = %d,"+
			"\n token itself = %s,\n refresh token = %s\n",
			new_token.TokenType, new_token.ExpiresIn, new_token.AccessToken, new_token.RefreshToken)
	}

}

func startBackgroundServer(serverAddress string, code chan string) *http.Server {
	srv := &http.Server{Addr: ":3000", Handler: http.DefaultServeMux}

	http.HandleFunc("/cb", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Thank you, you may close this page now."))
		r.Body.Close()
		code <- strings.TrimPrefix(r.URL.RawQuery, "code=")
	})

	go srv.ListenAndServe()

	return srv
}
