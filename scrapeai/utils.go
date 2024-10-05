package scrapeai

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

// Simple fetch functionality that retrieves data from a given url or returns
// the relevant err
func Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Get body from chromedp headless browswer to collect dynamically rendered
// content
func FetchFromChomedp(url string) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second), // Allow JS content to load
		chromedp.OuterHTML("html", &body),
	)
	if err != nil {
		return "", err
	}
	return body, nil
}

// Takes an html body as a string and returns a goquery document
func GoQueryDocFromBody(body string) (*goquery.Document, error) {
	reader := strings.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
