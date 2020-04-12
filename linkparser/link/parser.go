package link

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// ParseHTML parses HTML bytes for links
func ParseHTML(htmlBytes []byte) []Link {
	tokenizer := html.NewTokenizer(bytes.NewReader(htmlBytes))
	var links []Link
	var link Link
	var sb strings.Builder
	var inLinkTag bool
	for {
		z := tokenizer.Next()
		switch {
		case z == html.ErrorToken:
			return links
		case z == html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				link.Href = parseLink(token)
				inLinkTag = true
			}
		case z == html.TextToken:
			token := tokenizer.Token()
			if token.Data != "" && inLinkTag {
				sb.WriteString(strings.ReplaceAll(token.Data, "\n", ""))
			}
		case z == html.EndTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				link.Text = strings.TrimSpace(sb.String())
				links = append(links, link)
				link = Link{}
				inLinkTag = false
				sb.Reset()
			}
		}
	}
}

func parseLink(token html.Token) string {
	for _, a := range token.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}

	return ""
}
