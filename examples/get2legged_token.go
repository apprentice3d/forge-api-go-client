package main

import (
	"fmt"
	"github.com/apprentice3d/forge-api-go-client/oauth"
	"log"
	"os"
)

func main() {

	clientId := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	client := oauth.NewTwoLeggedClient(clientId, clientSecret)
	token, err := client.Authenticate("data:read")

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Received following information:\ntoken type = %s,\n expires in = %d,"+
		"\n token itself = %s,\n refresh token = %s",
		token.TokenType, token.ExpiresIn, token.AccessToken, token.RefreshToken)

}
