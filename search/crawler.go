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

// runCrawl is a function that performs a web crawl on a given URL.
// It sends a GET request to the URL, checks the response for errors, and parses the body if the response is HTML.
// If there is an error sending the request, the response is nil, the status code is not 200, or the content type is not text/html, it returns a CrawlData struct with Success set to false.
// If the body is successfully parsed, it returns a CrawlData struct with Success set to true and the parsed data.
// The function prints an error message and returns a CrawlData struct with Success set to false if there is an error parsing the body.
//
// Parameters:
// inputUrl string: The URL to perform the web crawl on.
//
// Returns:
// CrawlData: A struct containing the URL, whether the crawl was successful, the response code, and the parsed data from the body.
func runCrawl(inputUrl string) CrawlData {
	resp, err := http.Get(inputUrl)
	baseUrl, _ := url.Parse(inputUrl)
	// Check for error or if response is empty
	if err != nil || resp == nil {
		fmt.Println(err)
		fmt.Println("something went wrong fetch the body")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: 0, CrawlData: ParsedBody{}}
	}
	defer resp.Body.Close()
	// Check if response code is not 200
	if resp.StatusCode != 200 {
		fmt.Println(err)
		fmt.Println("status code is not 200")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParsedBody{}}
	}
	// Check the content type is text/html
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		// response is HTML
		data, err := parseBody(resp.Body, baseUrl)
		if err != nil {
			fmt.Println("something went wrong getting data from html body")
			return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParsedBody{}}
		}
		return CrawlData{Url: inputUrl, Success: true, ResponseCode: resp.StatusCode, CrawlData: data}
	} else {
		// response is not HTML
		fmt.Println("non html response detected")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParsedBody{}}
	}
}

// parseBody is a function that parses the body of a web page and extracts various information from it.
// It parses the body into an HTML node tree, extracts all the links, the title and description, and the h1 headings from the tree,
// and records the time it took to perform these operations.
// The function returns a ParsedBody struct containing the extracted information and the time it took to extract it.
// If there is an error parsing the body, the function prints an error message and returns an empty ParsedBody struct and the error.
//
// Parameters:
// body io.Reader: The body of the web page to parse.
// baseUrl *url.URL: The base URL of the web page to resolve relative URLs against and to compare with for determining if a link is internal or external.
//
// Returns:
// ParsedBody, error: A struct containing the extracted information and the time it took to extract it, and an error object that describes an error that occurred during the function's execution.
func parseBody(body io.Reader, baseUrl *url.URL) (ParsedBody, error) {
	doc, err := html.Parse(body)
	if err != nil {
		fmt.Println(err)
		fmt.Println("something went wrong parsing body")
		return ParsedBody{}, err
	}
	start := time.Now()
	// Get the links from the doc
	links := getLinks(doc, baseUrl)
	// Get the page title description
	title, desc := getPageData(doc)
	// Get the H1 tags for the page
	headings := getPageHeadings(doc)

	// Record timings
	end := time.Now()
	// Return the data
	return ParsedBody{
		CrawlTime:       end.Sub(start),
		PageTitle:       title,
		PageDescription: desc,
		Headings:        headings,
		Links:           links,
	}, nil
}

// getLinks is a function that extracts all the internal and external links from a given HTML node and its children.
// Depth First Search (DFS) of the html tree structure. This is a recursive function to scan the full tree.
// It uses a recursive function to traverse the HTML node tree and find all anchor tags.
// The href attribute of each anchor tag is parsed into a URL.
// If the URL is absolute and has the same host as the base URL, it is considered an internal link.
// If the URL is absolute and has a different host, it is considered an external link.
// If the URL is relative, it is resolved against the base URL and considered an internal link.
// The function ignores URLs that are a hashtag/anchor, mail link, telephone link, javascript link, or a PDF or MD file.
// The function returns a Links struct containing slices of internal and external links.
//
// Parameters:
// node *html.Node: The root HTML node to start the search from.
// baseUrl *url.URL: The base URL to resolve relative URLs against and to compare with for determining if a link is internal or external.
//
// Returns:
// Links: A struct containing slices of internal and external links.
func getLinks(node *html.Node, baseUrl *url.URL) Links {
	links := Links{}
	if node == nil {
		return links
	}
	var findLinks func(*html.Node)
	findLinks = func(node *html.Node) {
		// Check if the current node is an `html.ElementNode` and if it has a tag name of "a" (i.e., an anchor tag).
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					url, err := url.Parse(attr.Val)
					// Check for errors or if url is 1)a hashtag/anchor 2) is mail link, 3) is a telephone link, 4)is a javascript link 5) is a PDF or MD file
					if err != nil || strings.HasPrefix(url.String(), "#") || strings.HasPrefix(url.String(), "mail") || strings.HasPrefix(url.String(), "tel") || strings.HasPrefix(url.String(), "javascript") || strings.HasSuffix(url.String(), ".pdf") || strings.HasSuffix(url.String(), ".md") {
						continue
					}
					// If url is absolute then test if internal or extend before append. Else add the baseUrl append as internal
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
		// Recursively call function to do Depth First Search of entire tree
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			findLinks(child)
		}
	}
	findLinks(node)

	return links
}

// isSameHost is a function that checks if two URLs have the same host.
// It parses both URLs and compares their hosts.
// If there is an error parsing either URL, it returns false.
// If the hosts are the same, it returns true. Otherwise, it returns false.
//
// Parameters:
// absoluteURL string: The first URL to compare.
// baseURL string: The second URL to compare.
//
// Returns:
// bool: True if the hosts of both URLs are the same, false otherwise.
func isSameHost(absoluteURL string, baseURL string) bool {
	absURL, err := url.Parse(absoluteURL) // Parse the absolute URL. Example: https://example.com/path => Host: example.com
	if err != nil {
		return false
	}

	baseURLParsed, err := url.Parse(baseURL) // Parse the base URL. Example: https://example.com/path => Host: example.com
	if err != nil {
		return false
	}

	return absURL.Host == baseURLParsed.Host
}

// getPageData is a function that extracts the title and description from a given HTML node and its children.
// It uses a recursive function to traverse the HTML node tree and find the title and meta elements.
// The content of the title element and the content attribute of the meta element with name "description" are returned.
// If the HTML node is nil, it returns two empty strings.
// If the title element or the description meta element are not found, it returns an empty string for the respective value.
//
// Parameters:
// node *html.Node: The root HTML node to start the search from.
//
// Returns:
// string, string: The title and description of the page, respectively.
func getPageData(node *html.Node) (string, string) {
	if node == nil {
		return "", ""
	}
	// Find the page title and description
	title, desc := "", ""
	var findMetaAndTitle func(*html.Node)
	findMetaAndTitle = func(node *html.Node) {
		// Recursive function to search for `meta` elements in the HTML tree and extracts their `name` and `content` attributes.
		if node.Type == html.ElementNode && node.Data == "title" {
			// Check if first child is empty
			if node.FirstChild == nil {
				title = ""
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
				desc = content
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			findMetaAndTitle(child)
		}
	}
	findMetaAndTitle(node)
	return title, desc
}

// getPageHeadings is a function that extracts all the h1 headings from a given HTML node and its children.
// It uses a recursive function to traverse the HTML node tree and find all h1 elements.
// The content of each h1 element is appended to a string, separated by commas.
// If the HTML node is nil, it returns an empty string.
// The function removes the last comma and space from the concatenated string before returning it.
//
// Parameters:
// n *html.Node: The root HTML node to start the search from.
//
// Returns:
// string: A string containing all the h1 headings, separated by commas.
func getPageHeadings(n *html.Node) string {
	if n == nil {
		return ""
	}
	// Find all h1 elements and concatenate their content
	var headings strings.Builder
	var findH1 func(*html.Node)
	findH1 = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			// Check if first child is empty
			if n.FirstChild != nil {
				headings.WriteString(n.FirstChild.Data)
				headings.WriteString(", ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findH1(c)
		}
	}
	findH1(n)
	// Remove the last comma and space from the concatenated string & return
	return strings.TrimSuffix(headings.String(), ", ")
}
