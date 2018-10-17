package main

import (
	"github.com/apprentice3d/forge-api-go-client/md"
	"os"
	"log"
	"net/url"
	"archive/zip"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

var ()

func main() {

	args := os.Args

	if len(args) != 2 {
		log.Fatal("Error: missing argument\n\tUsage: bubble_downloader <urn>")
	}

	urnToGet := args[1]

	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		log.Fatal("Could not get Forge secrets from env")
	}

	mdApi := md.NewAPIWithCredentials(clientID, clientSecret)

	manifest, err := mdApi.GetManifest(urnToGet)
	if err != nil {
		log.Fatal("Could not get the manifest: ", err.Error())
	}

	var svfURN string

	for _, derivative := range manifest.Derivatives {
		if derivative.OutputType == "svf" {
			for _, resource := range derivative.Children {
				if resource.Type == "geometry" {
					for _, child := range resource.Children {
						if child.Role == "graphics" {
							svfURN = child.URN
							break
						}
					}

				}
			}
		}

	}

	log.Println("The urn for the svf file is ", svfURN)

	data, err := mdApi.GetDerivative(urnToGet, url.QueryEscape(svfURN))

	if err != nil {
		log.Fatal("Could not get the derivative, got ", err.Error())
	}

	err = os.MkdirAll("./sample", os.ModePerm)
	file, err := os.Create("./sample/model/Model.svf")
	if err != nil {
		log.Println("Warning: Could not save the SVF file: ", err.Error())
	}
	defer file.Close()
	file.Write(data)




	buffer := bytes.NewReader(data)

	zipper, err := zip.NewReader(buffer, buffer.Size())
	if err != nil {
		log.Fatal("Could not unzip the derivative, got ", err.Error())
	}

	var manifestFile md.LMVManifest
	for _, file := range zipper.File {
		if file.Name == "manifest.json" {
			log.Printf("Reding the manifest file: %s, size = %d", file.Name, file.UncompressedSize64)
			manifest, err := file.Open()
			if err != nil {
				log.Fatal("Could not read the manifest file, got ", err.Error())
			}
			data, err := ioutil.ReadAll(manifest)
			if err != nil {
				log.Fatal("Could not read the manifest content, got ", err.Error())
			}
			err = json.Unmarshal(data, &manifestFile)

			if err != nil {
				log.Fatal("Could not unmarshal data into json, got ", err.Error())
			}

		}

	}


	cacheFile, err := os.Create("./sample/model.manifest")
	if err != nil {
		log.Println("Could not open the manifest file, got ", err.Error())
	}
	defer cacheFile.Close()

	commonPath := svfURN[:getPositionOfLastSlash(svfURN)]

	if err != nil {
		log.Fatal("Could not create the output directory, got ", err.Error())
	}
	for _, asset := range manifestFile.Assets {
		log.Println(commonPath + asset.URI)
		err = os.MkdirAll(filepath.Dir("./sample/model/" + asset.URI), os.ModePerm)
		if err != nil {
			log.Fatalf("Could not create the nested folders, got %s", err.Error())
		}
		file, err := os.Create("./sample/model/" + asset.URI)
		if err != nil {
			log.Fatalf("Could not create file %s, got %s", asset.URI, err.Error())
		}
		data, err := mdApi.GetDerivative(urnToGet, commonPath + asset.URI)
		if err != nil {
			log.Fatalf("Could not get derivative of %s, got %s", asset.URI, err.Error())
		}
		file.Write(data)
		file.Close()
		_, err = cacheFile.WriteString(asset.URI+"\n")
		if err != nil {
			log.Println("Problem writing to Manifest file: ", err.Error())
		}
	}


	log.Println("Looks like everything is ok now.")

}




func getPositionOfLastSlash(fullPath string) int {
	for i:=len(fullPath)-1; i > 0; i-- {
		if fullPath[i] == '/' {
			return i+1
		}
	}
	return 0
}