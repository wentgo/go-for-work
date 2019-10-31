package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetTrackingNumber(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return ""
	}

	var trackingNum string

	// Find the review items
	doc.Find(".carrierRelatedInfo-trackingId-text").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		trackingNum = s.Text()
	})

	return strings.Replace(trackingNum, "Tracking ID: ", "", -1)
}

func main() {
	filename := "114-4939539-3641016.html"
	trackingNum := GetTrackingNumber(filename)
	fmt.Printf("%s\n", trackingNum)
}
