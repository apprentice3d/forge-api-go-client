package main

import (
	"encoding/json"
	"fmt"
	"github.com/apprentice3d/forge-api-go-client/dm"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	dir, _ := os.Getwd()
	data, err := ioutil.ReadFile(dir + "/dm/hubs_return_sample.json")

	if err != nil {
		log.Fatal("Could not read file: " + err.Error())
	}

	decoder := json.NewDecoder(strings.NewReader(string(data)))

	hub := new(dm.Hubs)

	decoder.Decode(hub)

	fmt.Println(hub)

	//var dat map[string]interface{}

	hubik := dm.Hubs{}
	//json.Unmarshal(data,&dat)
	json.Unmarshal(data, &hubik)

	//fmt.Println(dat)
	fmt.Println(hubik)

	fmt.Println(hubik.Links.Self.Href)
	//fmt.Println(dat["data"])

}
