# Go Web Crawler

## Overview

This project implements a web crawler using the Go programming language and the Colly framework. The crawler extracts text information from specific Wikipedia pages related to intelligent systems and robotics, saves the HTML content of each page, and writes the extracted data to a JSON lines file. It also measures and logs the time taken to complete the scraping process.

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/bilguunbilegt/crawler.git
   cd crawler
   ```

2. Initialize the Go module:
   ```bash
   go mod init crawler
   ```

3. Install dependencies:
   ```bash
   go get -u github.com/gocolly/colly/v2
   go get -u github.com/PuerkitoBio/goquery
   ```

## Running the Crawler

To run the crawler and extract data, use the following command:
```bash
go run crawler.go
```

The extracted data will be saved in `output.jsonl` and the HTML content will be saved in the `wikipages` directory. The total time taken for the scraping process will be logged in the terminal output.

## Testing

To run unit tests for the critical components, use the following command:
```bash
go test
```

## Directory Structure

```
go-web-crawler/
├── wikipages/
├── crawler.go
├── crawler_test.go
├── go.mod
├── go.sum
└── README.md
```

## Explanation of Code

- **crawler.go**: Contains the main logic for the web crawler.
  - Initializes the Colly collector.
  - Defines callbacks for saving HTML content and extracting page data.
  - Measures and logs the time taken for the scraping process.

- **crawler_test.go**: Contains unit tests for critical components of the crawler.
  - Tests the `writeJSONLine` function.
  - Tests the `saveHTMLFile` function.

