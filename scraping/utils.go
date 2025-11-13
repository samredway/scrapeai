// Package scraping contains utility functions for web scraping
package scraping

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

// Simple fetch functionality that retrieves data from a given url or returns
// the relevant err. Timeouts should be set on the context
func Fetch(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
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

// Get body from chromedp headless browswer to collect dynamically rendered content
func FetchFromChromedp(ctx context.Context, url string) (string, error) {
	// Suppress chromedp's internal logging by using a no-op logger
	chromedpCtx, cancel := chromedp.NewContext(
		ctx,
		chromedp.WithLogf(func(string, ...any) {}),
	)
	defer cancel()

	var body string
	err := chromedp.Run(chromedpCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second), // Allow JS content to load
		chromedp.OuterHTML("html", &body),
	)
	if err != nil {
		return "", err
	}
	return body, nil
}

// FetchWithZyteStaticProxy fetches a URL using Zyte's static proxy.
// The ZYTE_API_KEY environment variable must be set.
// The proxy endpoint defaults to api.zyte.com:8011 but can be overridden
// via the ZYTE_PROXY_ENDPOINT environment variable.
func FetchWithZyteStaticProxy(ctx context.Context, targetURL string) (string, error) {
	apiKey := os.Getenv("ZYTE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("ZYTE_API_KEY is not set in the environment")
	}

	proxyEndpoint := os.Getenv("ZYTE_PROXY_ENDPOINT")
	if proxyEndpoint == "" {
		proxyEndpoint = "http://api.zyte.com:8011"
	}

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	// Add Zyte-Browser-Html header for rendered HTML (optional but recommended)
	req.Header.Set("Zyte-Browser-Html", "true")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	// Parse proxy endpoint URL and embed API key in the URL for authentication
	// Format: http://<api_key>:@api.zyte.com:8011
	// The empty password after the colon matches curl's --proxy-user format
	parsedProxy, err := url.Parse(proxyEndpoint)
	if err != nil {
		return "", fmt.Errorf("invalid proxy endpoint: %w", err)
	}

	// Set user info: API key as username, empty password
	parsedProxy.User = url.UserPassword(apiKey, "")
	proxyURL := parsedProxy

	// Configure HTTP client to use the proxy
	// Note: Zyte's proxy performs SSL interception/TLS termination, presenting its own
	// certificate instead of the target server's certificate. We skip TLS verification
	// for the target connection since the proxy handles the actual TLS to the target.
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("requesting through proxy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf(
			"proxy request failed with status %d: %s",
			resp.StatusCode,
			body,
		)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	return string(body), nil
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

// StripNonTextTags removes elements that don't contain text from a copy of the
// given goquery document
// Returns a new document with non-text elements removed
func StripNonTextTags(doc *goquery.Document) *goquery.Document {
	docCopy := goquery.CloneDocument(doc)
	docCopy.Find("script, style, link, meta").Remove()
	docCopy.Find("*").Each(func(i int, s *goquery.Selection) {
		if strings.TrimSpace(s.Text()) == "" {
			s.Remove()
		}
	})
	return docCopy
}

// Helper function to get the HTML string from a GoQuery document
func GetDocumentHTML(doc *goquery.Document) (string, error) {
	html, err := doc.Html()
	if err != nil {
		return "", err
	}
	return html, nil
}
