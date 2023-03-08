package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.4; Win64; x64) Gecko/20130401 Firefox/60.7",
	"Mozilla/5.0 (Linux; Linux i676 ; en-US) AppleWebKit/602.37 (KHTML, like Gecko) Chrome/53.0.3188.182 Safari/533",
	"Mozilla/5.0 (Android; Android 6.0.1; SM-G928S Build/MDB08I) AppleWebKit/536.43 (KHTML, like Gecko)  Chrome/53.0.1961.231 Mobile Safari/602.4",
	"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_1_7) Gecko/20100101 Firefox/48.4",
	"Mozilla/5.0 (Windows NT 10.1; Win64; x64) AppleWebKit/602.6 (KHTML, like Gecko) Chrome/49.0.2764.174 Safari/603",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 9_8_7) Gecko/20130401 Firefox/59.7",
	"Mozilla/5.0 (Linux; U; Android 5.0; LG-D332 Build/LRX22G) AppleWebKit/602.2 (KHTML, like Gecko)  Chrome/54.0.2014.306 Mobile Safari/600.5",
	"Mozilla/5.0 (Android; Android 5.0.2; HTC 80:number1-2w Build/LRX22G) AppleWebKit/600.31 (KHTML, like Gecko)  Chrome/54.0.3164.157 Mobile Safari/603.5",
	"Mozilla/5.0 (Linux; Android 7.1; Nexus 8P Build/NPD90G) AppleWebKit/601.15 (KHTML, like Gecko)  Chrome/47.0.3854.218 Mobile Safari/600.6",
	"Mozilla/5.0 (Android; Android 5.1; SM-G928L Build/LRX22G) AppleWebKit/602.47 (KHTML, like Gecko)  Chrome/48.0.1094.253 Mobile Safari/537.6",
}

var tokens = make(chan struct{}, 5)

func RandomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func DiscoverLinks(response *http.Response, baseURL string) []string {
	if response != nil {
		doc, _ := goquery.NewDocumentFromResponse(response)
		foundURLs := []string{}

		if doc != nil {

			doc.Find("a").Each(func(i int, s *goquery.Selection) {
				ref, _ := s.Attr("href")
				foundURLs = append(foundURLs, ref)
			})

		} else {
			return []string{}
		}

		return foundURLs
	} else {
		return []string{}
	}
}

func GetRequest(targetURL string) (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", RandomUserAgent())
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func CheckRelative(href string, baseURL string) string {
	if strings.HasPrefix(href, "/") {
		return fmt.Sprintf("%s%s", baseURL, href)
	} else {
		return href
	}
}

func ResolveRelativeLinks(href string, baseURL string) (bool, string) {
	resultHREF := CheckRelative(href, baseURL)

	baseParse, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println(err.Error())
	}

	resultHREFParse, err := url.Parse(resultHREF)
	if err != nil {
		fmt.Println(err.Error())
	}

	if baseParse != nil && resultHREFParse != nil {
		if baseParse.Host == resultHREFParse.Host {
			return true, resultHREF
		} else {
			return false, ""
		}
	}
	return false, ""
}

func Crawl(targetURL string, baseURL string) []string {
	fmt.Println("Crawling : ", targetURL)

	tokens <- struct{}{}

	resp, _ := GetRequest(targetURL)

	<-tokens

	links := DiscoverLinks(resp, baseURL)

	foundURLs := []string{}
	for _, link := range links {
		ok, correctLink := ResolveRelativeLinks(link, baseURL)
		if ok {
			if correctLink != "" {
				foundURLs = append(foundURLs, correctLink)
			}
		}
	}
	// we can also parseHTML to scrape some data
	return foundURLs
}

func main() {
	worklist := make(chan []string)
	baseDomain := "https://www.theguardian.com"
	go func() { worklist <- []string{"https://www.theguardian.com"} }()
	seen := make(map[string]bool)
	n := 1

	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string, baseURL string) {
					foundLinks := Crawl(link, baseURL)

					if foundLinks != nil {
						worklist <- foundLinks
					}
				}(link, baseDomain)
			}
		}
	}
}
