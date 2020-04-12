package sitemap

import (
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/ericywl/gophercises/linkparser/link"
)

// BuildSiteMap builds a XML sitemap for the given site
func BuildSiteMap(siteUrl string) (string, error) {
	url := strings.TrimSpace(siteUrl)
	domainName := strings.Trim(url, "/")
	linkSet := map[string]bool{}
	err := recursiveScrapLinks(url, domainName, linkSet)
	if err != nil {
		return "", err
	}

	links := getMapKeys(linkSet)
	sort.Strings(links)
	output, err := buildXMLSiteMap(links, "http://www.sitemaps.org/schemas/sitemap/0.9")
	if err != nil {
		return "", err
	}

	return "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + string(output), nil
}

func getMapKeys(linkSet map[string]bool) []string {
	var keys []string
	for k := range linkSet {
		keys = append(keys, k)
	}

	return keys
}

func recursiveScrapLinks(url string, domainName string, linkSet map[string]bool) error {
	links, err := scrapLinks(url)
	if err != nil {
		return err
	}

	linkSet[url] = true
	links = filterDomain(links, domainName)
	for _, l := range links {
		newUrl := l.Href
		if strings.HasPrefix(newUrl, "/") {
			newUrl = domainName + "/" + newUrl[1:]
		}

		if _, ok := linkSet[newUrl]; ok {
			continue
		}

		err = recursiveScrapLinks(newUrl, domainName, linkSet)
		if err != nil {
			return err
		}
	}

	return nil
}

func scrapLinks(url string) ([]link.Link, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return link.ParseHTML(body), nil
}

func filterDomain(links []link.Link, domainName string) []link.Link {
	var filteredLinks []link.Link
	for _, l := range links {
		if strings.HasPrefix(l.Href, "/") || strings.HasPrefix(l.Href, domainName) {
			filteredLinks = append(filteredLinks, l)
		}
	}

	return filteredLinks
}
