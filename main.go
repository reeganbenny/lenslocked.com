package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
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
		<sitemap>
			<loc>http://www.washingtonpost.com/news-blogs-goingoutguide-sitemap.xml</loc>
		</sitemap>
		<sitemap>
			<loc>http://www.washingtonpost.com/news-goingoutguide-sitemap.xml</loc>
		</sitemap>
	</sitemapindex>
`)

//SitemapIndex is for Xml extraction
type SitemapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

//News extracted from xml
type News struct {
	Titles []string `xml:"url>loc"`
}

func main() {
	var s SitemapIndex
	var n News
	bytes := washPostXML
	xml.Unmarshal(bytes, &s)

	// fmt.Println(s.Locations)

	for _, Location := range s.Locations {
		// fmt.Println(Location)
		resp, _ := http.Get(Location)
		bytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		xml.Unmarshal(bytes, &n)
	}

	// fmt.Println(n.Titles)
}
