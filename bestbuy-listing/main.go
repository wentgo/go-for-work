package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "https://marketplace.bestbuy.ca"
)

var (
	username = os.Getenv("BB_USER")
	password = os.Getenv("BB_PASS")
)

type App struct {
	Client *http.Client
}

func (app *App) getToken() string {
	log.Println("Get Token")

	client := app.Client

	loginURL := baseURL + "/login"

	response, err := client.Get(loginURL)
	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	token, _ := document.Find("input[name='_csrf']").Attr("value")

	return token
}

func (app *App) login() {
	log.Println("Login")

	client := app.Client

	token := app.getToken()

	authURL := baseURL + "/authenticate"

	data := url.Values{
		"_csrf":    {token},
		"username": {username},
		"password": {password},
	}

	response, err := client.PostForm(authURL, data)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// How to detect whether login is successful or not
}

func (app *App) getListings(filename string) {
	log.Println("Downloading ...")

	client := app.Client

	downloadURL := baseURL + "/mmp/shop/offer/export/shop-export?empty=true"

	response, err := client.Get(downloadURL)
	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}
	defer response.Body.Close()

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln("Error creating file. ", err)
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, response.Body)

	log.Println("Completed")
}

func main() {
	jar, _ := cookiejar.New(nil)

	app := App{
		Client: &http.Client{Jar: jar},
	}

	app.login()
	app.getListings("bestbuy-listing.csv")
}
