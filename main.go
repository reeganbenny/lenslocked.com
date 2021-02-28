package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var washPostXML = []byte(`
	<sitemapindex>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-politics-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-blogs-politics-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-opinions-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-blogs-opinions-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-local-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-blogs-local-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-sports-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-blogs-sports-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-national-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-blogs-national-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-world-sitemap.xml </loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-blogs-world-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-business-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-blogs-business-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-technology-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-blogs-technology-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-lifestyle-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-blogs-lifestyle-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-entertainment-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-blogs-entertainment-sitemap.xml</loc>
		</sitemap>
	</sitemapindex>
`)

// SitemapIndex is for Xml extraction
type SitemapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

//News extracted from xml
type News struct {
	Titles    []string `xml:"url>news>title"`
	Locations []string `xml:"url>loc"`
}

//NewsMap is the value and Title will be the key
type NewsMap struct {
	Topic    string
	Location string
}

func main() {
	var s SitemapIndex
	var n News
	newsmap := make(map[string]NewsMap)

	bytes := washPostXML
	xml.Unmarshal(bytes, &s)
	for _, Location := range s.Locations {
		locationTopic := strings.Split(Location, "/")
		topic := strings.ReplaceAll(locationTopic[3], "-", " ")
		topic = strings.ReplaceAll(topic, "sitemap.xml", "")
		resp, _ := http.Get(Location)
		bytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		xml.Unmarshal(bytes, &n)
		for idx := range n.Titles {
			newsmap[n.Titles[idx]] = NewsMap{topic, n.Locations[idx]}
		}
	}

	for idx, data := range newsmap {
		fmt.Println("\n\n\nTitle := ", idx)
		fmt.Println("\nLocation =", data.Location)
		fmt.Println("\nLocation =", data.Topic)
	}
}
