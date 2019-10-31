package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/unidoc/unipdf/v3/extractor"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func main() {
	files, _ := filepath.Glob("*.pdf")
	for _, file := range files {
		if len(file) <= 8 { // 1234.pdf
			renameFile(file)
		}
	}

	//fmt.Println("Press the Enter Key to exit")
	//fmt.Scanln() // wait for Enter Key
}

func renameFile(filename string) {
	text, err := extractText(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	trackingNum := getTrackingNum(text)
	invoiceNo := getInvoiceNo(text)

	fmt.Printf("Filename:  %s\n", filename)
	fmt.Printf("Tracking#: %s\n", trackingNum)
	fmt.Printf("Invoice#:  %s\n", invoiceNo)

	if invoiceNo != "" {
		newname := invoiceNo + ".pdf"
		err = os.Rename(filename, newname)
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			return
		}
		fmt.Printf("RenameTo:  %s\n", newname)
	}

	fmt.Println()
}

func extractText(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return "", err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}

	var text string
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return "", err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return "", err
		}

		str, err := ex.ExtractText()
		if err != nil {
			return "", err
		}
		text += str
	}

	return text, nil
}

// TR#:1Z 37Y 059 91 5244 5990
// TRACKING #: 1Z 37Y 059 91 5244 5990

func getTrackingNum(text string) string {
	var patterns = []*regexp.Regexp{
		regexp.MustCompile(`TRACKING #: (.*)`),
		regexp.MustCompile(`TR#:(.*)`),
	}

	for _, re := range patterns {
		matches := re.FindStringSubmatch(text)
		//fmt.Printf("%#v\n", matches)
		if matches != nil {
			return strings.Replace(matches[1], " ", "", -1)
		}
	}

	return ""
}

// Invoice No.: 111-5941991-0261830-R BILLING: F/C Receiver 37Y059
// Purchase No.: 111-5937548-5700226 DESC: pc switch INV-RS

func getInvoiceNo(text string) string {
	var patterns = []*regexp.Regexp{
		regexp.MustCompile(`REF 1:([-\d\w]*)`),
		regexp.MustCompile(`REF 2:([-\d\w]*)`),
		regexp.MustCompile(`Invoice No.: ([-\d\w]*)`),
		regexp.MustCompile(`Purchase No.: ([-\d\w]*)`),
	}

	for _, re := range patterns {
		matches := re.FindStringSubmatch(text)
		//fmt.Printf("%#v\n", matches)
		if matches != nil {
			return matches[1]
		}
	}

	return ""
}
