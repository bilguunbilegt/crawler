package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// PageData holds the data for each page
type PageData struct {
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}

// writeJSONLine writes the PageData to a JSON lines file
func writeJSONLine(outputFile *os.File, data PageData) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Could not marshal data: %v", err)
	}
	_, err = outputFile.WriteString(fmt.Sprintf("%s\n", jsonData))
	if err != nil {
		log.Fatalf("Could not write to file: %v", err)
	}
}

// saveHTMLFile saves the HTML content of the page
func saveHTMLFile(pageDirectory, filename string, content []byte) {
	filePath := filepath.Join(pageDirectory, filename)
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Fatalf("Could not save HTML file: %v", err)
	}
}

// extractTags extracts tags from the URL parts
func extractTags(urlParts []string, regex *regexp.Regexp) []string {
	var tags []string
	if len(urlParts) > 1 {
		tags = append(tags, urlParts[1])
	}
	if len(urlParts) > 2 {
		tags = append(tags, urlParts[2])
	}
	if len(urlParts) > 3 {
		moreTags := strings.Split(urlParts[3], "_")
		for _, tag := range moreTags {
			cleanedTag := regex.ReplaceAllString(strings.ToLower(tag), "")
			tags = append(tags, cleanedTag)
		}
	}
	return tags
}

// cleanText extracts and cleans the visible text from the HTML element
func cleanText(selection *goquery.Selection) string {
	var textBuilder strings.Builder
	selection.Find("p, li").Each(func(i int, s *goquery.Selection) {
		textBuilder.WriteString(s.Text())
		textBuilder.WriteString("\n")
	})
	return strings.TrimSpace(textBuilder.String())
}

func main() {
	startTime := time.Now() // Begin recording the start time

	// Initializing a new collector
	collector := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	// Directory to save HTML pages
	pageDirectory := "wikipages"
	if err := os.MkdirAll(pageDirectory, os.ModePerm); err != nil {
		log.Fatalf("Could not create directory: %v", err)
	}

	// Create JSON lines file for output
	outputFile, err := os.Create("output.jl")
	if err != nil {
		log.Fatalf("Could not create output file: %v", err)
	}
	defer outputFile.Close()

	// Regex for cleaning tags
	regex := regexp.MustCompile("[^a-zA-Z]")

	// Callback for saving HTML content
	collector.OnResponse(func(response *colly.Response) {
		page := strings.Split(response.Request.URL.Path, "/")[2]
		filename := fmt.Sprintf("%s.html", page)
		saveHTMLFile(pageDirectory, filename, response.Body)
	})

	// Callback for extracting page data
	collector.OnHTML("body", func(element *colly.HTMLElement) {
		title := element.DOM.Find("h1").Text()
		text := cleanText(element.DOM.Find("div#mw-content-text"))
		urlParts := strings.Split(element.Request.URL.Path, "/")
		tags := extractTags(urlParts, regex)
		pageData := PageData{
			URL:   element.Request.URL.String(),
			Title: title,
			Text:  text,
			Tags:  tags,
		}
		writeJSONLine(outputFile, pageData)
	})

	// Callback for logging requests
	collector.OnRequest(func(request *colly.Request) {
		log.Println("Visiting", request.URL)
	})

	// List of pages to scrape
	pages := []string{
		"https://en.wikipedia.org/wiki/Robotics",
		"https://en.wikipedia.org/wiki/Robot",
		"https://en.wikipedia.org/wiki/Reinforcement_learning",
		"https://en.wikipedia.org/wiki/Robot_Operating_System",
		"https://en.wikipedia.org/wiki/Intelligent_agent",
		"https://en.wikipedia.org/wiki/Software_agent",
		"https://en.wikipedia.org/wiki/Robotic_process_automation",
		"https://en.wikipedia.org/wiki/Chatbot",
		"https://en.wikipedia.org/wiki/Applications_of_artificial_intelligence",
		"https://en.wikipedia.org/wiki/Android_(robot)",
	}

	// Start scraping each page
	for _, page := range pages {
		collector.Visit(page)
	}

	collector.Wait() // Wait for all requests to finish

	elapsedTime := time.Since(startTime) // Calculate time
	log.Printf("Total time elapsed: %s", elapsedTime)
}
