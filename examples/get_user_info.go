package main

import (
	"fmt"
	"github.com/apprentice3d/forge-api-go-client/oauth"
	"log"
	"os"
)

func main() {

	log.Fatal("Example not finished yet.")

	clientId := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	client := oauth.NewTwoLeggedClient(clientId, clientSecret)
	token, err := client.Authenticate("data:read")

	if err != nil {
		log.Fatal(err.Error())
	}

	info := oauth.NewInformationalApi()
	profile, err := info.AboutMe(*token)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%v", profile)

}
