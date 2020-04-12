package sitemap

import "encoding/xml"

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS string `xml:"xmlns,attr"`
	URLs []URL
}

type URL struct {
	XMLName xml.Name `xml:"url"`
	Location string `xml:"loc"`
}

func buildXMLSiteMap(links []string, xmlns string) ([]byte, error) {
	var urls []URL
	for _, l := range links {
		urls = append(urls, URL{
			Location: l,
		})
	}

	urlSet := URLSet{
		XMLNS: xmlns,
		URLs: urls,
	}

	return xml.MarshalIndent(urlSet, "", "\t")
}
