package search

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type CrawlData struct {
	Url          string
	Success      bool
	ResponseCode int
	CrawlData    ParsedBody
}

type ParsedBody struct {
	CrawlTime       time.Duration
	PageTitle       string
	PageDescription string
	Headings        string
	Links           Links
}

type Links struct {
	Internal []string
	External []string
}

// runCrawl performs a crawl operation on the given URL.
// It sends a GET request to the URL, checks the response, and parses the body if it's HTML.
// If the response is not 200 or the content type is not HTML, it returns a CrawlData struct with Success set to false.
// If the body is successfully parsed, it returns a CrawlData struct with the parsed data and Success set to true.
//
// Parameters:
//   inputUrl: The URL to crawl.
//
// Returns:
//   CrawlData: A struct containing the URL, whether the crawl was successful, the response code, and the parsed data.
func runCrawl(inputUrl string) CrawlData {
	resp, err := http.Get(inputUrl)
	baseUrl, err := url.Parse(inputUrl)
	// Check if error or if response code is empty
	if err != nil || resp == nil {
		fmt.Println("something went wrong fetching the body")
		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: 0,
			CrawlData:    ParsedBody{},
		}
	}
	defer resp.Body.Close() // defer: execute after the function returns
	// Check for 200
	if resp.StatusCode != 200 {
		fmt.Println("Non 200 code found")
		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: resp.StatusCode,
			CrawlData:    ParsedBody{},
		}
	}
	// Check for html
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		// response in html
		data, err := parseBody(resp.Body, baseUrl)
		if err != nil {
			return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParsedBody{}}

		}
		return CrawlData{
			Url:          inputUrl,
			Success:      true,
			ResponseCode: resp.StatusCode,
			CrawlData:    data,
		}
	} else {
		// Response is not html
		return CrawlData{
			Url:          inputUrl,
			Success:      true,
			ResponseCode: resp.StatusCode,
			CrawlData:    ParsedBody{},
		}
	}
}

// parseBody parses the body of a web page and extracts various information.
// It parses the HTML, gets all internal and external links, the page title and description, and all h1 headings.
// It also measures the time it takes to perform the parsing and extraction.
// If the HTML cannot be parsed, it returns an empty ParsedBody struct and an error.
//
// Parameters:
//   body: The body of the web page as an io.Reader.
//   baseUrl: The base URL of the page to compare with absolute URLs.
//
// Returns:
//   ParsedBody: A struct containing the crawl time, page title, page description, h1 headings, and links.
//   error: An error that will be non-nil if there was an issue parsing the HTML.
func parseBody(body io.Reader, baseUrl *url.URL) (ParsedBody, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return ParsedBody{}, err
	}
	start := time.Now()
	// Get Links
	links := getLinks(doc, baseUrl)
	// Get Page Title & Description
	title, description := getPageData(doc)
	// Get H1 Tags
	headings := getPageHeadings(doc)
	// Return the time & data

	end := time.Now()

	return ParsedBody{
		CrawlTime:       end.Sub(start),
		PageTitle:       title,
		PageDescription: description,
		Headings:        headings,
		Links:           links,
	}, nil
}

// Depth first search (DFS) - recursive function to scanning the html tree
// getLinks traverses the HTML node tree to find and return all internal and external links in the page.
// It uses a recursive helper function, findLinks, to traverse the tree.
// If the node is nil, it returns an empty Links struct.
// The function distinguishes between absolute and relative URLs, and checks if absolute URLs are from the same host.
// It ignores links that start with "#", "mail", "tel", ".pdf", and ".md".
//
// Parameters:
//   node: The root of the HTML node tree.
//   baseUrl: The base URL of the page to compare with absolute URLs.
//
// Returns:
//   Links: A struct containing two slices of strings, one for internal links and one for external links.
func getLinks(node *html.Node, baseUrl *url.URL) Links {
	links := Links{}
	if node == nil {
		return links
	}

	var findLinks func(*html.Node)

	findLinks = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					url, err := url.Parse(attr.Val)
					if err != nil || strings.HasPrefix(url.String(), "#") || strings.HasPrefix(url.String(), "mail") || strings.HasPrefix(url.String(), "tel") || strings.HasPrefix(url.String(), ".pdf") || strings.HasPrefix(url.String(), ".md") {
						continue
					}

					// q: what is a absolute url?
					// a: An absolute URL is a full URL that includes everything you need to find a specific page, including the protocol (HTTP or HTTPS), the domain name, and the path. (e.g. https://www.example.com/page.html)

					// q: What is a relative URL?
					// a: A relative URL is a URL that only includes the path to a specific page, without the protocol or domain name. (e.g. /page.html)
					if url.IsAbs() {
						if isSameHost(url.String(), baseUrl.String()) {
							links.Internal = append(links.Internal, url.String())
						} else {
							links.External = append(links.External, url.String())
						}
					} else {
						rel := baseUrl.ResolveReference(url)
						links.Internal = append(links.Internal, rel.String())
					}

				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			findLinks(child)
		}
	}

	findLinks(node)
	return links
}

// isSameHost checks if the host of the absoluteUrl is the same as the host of the baseUrl.
// It parses both URLs and compares their hosts.
// If either URL cannot be parsed, it returns false.
//
// Parameters:
//
//	absoluteUrl: The absolute URL to check.
//	baseUrl: The base URL to compare with.
//
// Returns:
//
//	bool: True if the hosts of both URLs are the same, false otherwise.
func isSameHost(absoluteUrl string, baseUrl string) bool {
	absURL, err := url.Parse(absoluteUrl)
	if err != nil {
		return false
	}

	baseURLParsed, err := url.Parse(baseUrl)

	if err != nil {
		return false
	}

	return absURL.Host == baseURLParsed.Host
}

// getPageData traverses the HTML node tree to find and return the title and description of the page.
// It uses a recursive helper function, findMetaAndTitle, to traverse the tree.
// If the node is nil, it returns empty strings.
//
// Parameters:
//   node: The root of the HTML node tree.
//
// Returns:
//   string: The title of the page. If no title is found, returns an empty string.
//   string: The description of the page. If no description is found, returns an empty string.
func getPageData(node *html.Node) (string, string) {
	if node == nil {
		return "", ""
	}
	// Find title and description
	title, description := "", ""

	var findMetaAndTitle func(*html.Node)
	findMetaAndTitle = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "title" {
			// check if empty
			if node.FirstChild == nil {
				title = " "
			} else {
				title = node.FirstChild.Data
			}
		} else if node.Type == html.ElementNode && node.Data == "meta" {
			var name, content string
			for _, attr := range node.Attr {
				if attr.Key == "name" {
					name = attr.Val
				} else if attr.Key == "content" {
					content = attr.Val
				}
			}
			if name == "description" {
				description = content
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findMetaAndTitle(child)
	}

	findMetaAndTitle(node)

	return title, description
}

// getPageHeadings traverses the HTML node tree to find and return all h1 headings of the page.
// It uses a recursive helper function, findH1, to traverse the tree.
// If the node is nil, it returns an empty string.
// The headings are returned as a comma-separated string.
//
// Parameters:
//   n: The root of the HTML node tree.
//
// Returns:
//   string: The h1 headings of the page. If no h1 headings are found, returns an empty string.
func getPageHeadings(n *html.Node) string {
	if n == nil {
		return ""
	}
	var headings strings.Builder
	var findH1 func(*html.Node)
	findH1 = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			// Check if node is empty
			if n.FirstChild != nil {
				headings.WriteString(n.FirstChild.Data)
				headings.WriteString(", ")

			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findH1(c)
		}
	}

	// Remove the last comma
	return strings.TrimSuffix(headings.String(), ",")
}
