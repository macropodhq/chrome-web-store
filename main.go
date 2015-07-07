package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	appID := flag.String("a", "", "appId")
	pkgFileName := flag.String("p", "", "path to zipped extension")
	credentails := flag.String("c", "", "path to google api credentials json")
	email := flag.String("e", "", "developer account email")
	flag.Parse()

	if *appID == "" || *pkgFileName == "" || *credentails == "" || *email == "" {
		fmt.Println("Not all required paramters have been specified")
		flag.PrintDefaults()
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(*credentails)
	if err != nil {
		log.Fatal(err)
	}
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/chromewebstore")
	if err != nil {
		log.Fatal(err)
	}
	conf.Subject = *email

	client := conf.Client(oauth2.NoContext)

	pkg, err := os.Open(*pkgFileName)
	if err != nil {
		log.Panicln(err)
	}

	req, err := http.NewRequest("PUT", "https://www.googleapis.com/upload/chromewebstore/v1.1/items/"+*appID, pkg)
	if err != nil {
		log.Panicln(err)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Panicln(err)
	}

	defer res.Body.Close()

	fmt.Println("Status", res.Status)

	dec := json.NewDecoder(res.Body)

	response := struct {
		Kind        string              `json:"kind"`
		UploadState string              `json:"uploadState"`
		ItemError   []map[string]string `json:"itemError"`
	}{}
	err = dec.Decode(&response)
	if err != nil {
		log.Panicln(err)
	}

	if response.UploadState == "SUCCESS" {
		fmt.Println("Upload Complete")
	} else {
		fmt.Printf("%+v", response)
		os.Exit(1)
	}

	res, err = client.PostForm("https://www.googleapis.com/chromewebstore/v1.1/items/"+*appID+"/publish", url.Values{})

	if err != nil {
		log.Panicln(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Print(string(body))
}
