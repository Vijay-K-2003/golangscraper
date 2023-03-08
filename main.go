package main

import "fmt"

userAgents = []string{

}

func RandomUserAgent() {

}

func discoverLinks() {

}

func getRequest() {

}

func resolveRelativeLinks() {
	
}

func Crawl(targetURL string, baseURL string) []string {
	fmt.Println("Crawling : ", targetURL)
	resp, _ = getRequest(targetURL)
	links := discoverLinks(resp, baseURL)
	foundURLs := []string{}
	for _, link := range links {
		resolveRelativeLinks(link, baseURL)
	}
}

func main() {
	worklist := make(chan []string)
	baseDomain := "https://www.theguardian.com"
	go func(){	worklist <- []string{"https://www.theguardian.com"}} ()
	seen := make(map[string]bool)
	n := 1

	for ; n > 0 ; n-- {
		list := worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func (link string, baseURL string) {
					foundLinks := Crawl(link, baseURL)
	
					if foundLinks != nil {
						worklist <- foundLinks
					}
				}
			}
		}
	}
}
