package main

import (
	"fmt"

	"github.com/samredway/scrapeai/scrapeai"
)

func main() {
	result, err := scrapeai.Scrape(
		scrapeai.ScrapeAiRequest{
			Url:    "https://www.google.com/about/careers/applications/jobs/results#!t=jo&jid=127025001&",
			Prompt: "Extract the job titles from the page.",
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
