package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

//SitemapIndex is for Xml extraction
type SitemapIndex struct {
	Locataions []Location `xml:"sitemap"`
}

//Location is for xml extration
type Location struct {
	Loc string `xml:"loc"`
}

func (L Location) String() string {
	return fmt.Sprintf(L.Loc)
}

func main() {
	resp, _ := http.Get("https://www.washingtonpost.com/sitemaps/index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var s SitemapIndex
	xml.Unmarshal(bytes, &s)

	fmt.Println(s.Locataions)
}
