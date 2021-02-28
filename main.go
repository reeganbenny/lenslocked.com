package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
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

//NewsMapPage is the struct used to pass to html Page
type NewsMapPage struct {
	News map[string]NewsMap
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

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
		log.Println("Extracting news for ", topic)
		xml.Unmarshal(bytes, &n)
		log.Println("Length of Titles :-", len(n.Titles))
		for idx := range n.Titles {
			if val, ok := newsmap[n.Titles[idx]]; ok {
				fmt.Println("Repetition value in newsmap := ", val)
			}
			newsmap[n.Titles[idx]] = NewsMap{topic, n.Locations[idx]}
		}
		log.Println("Current newsmap size :=", len(newsmap))
	}
	p := NewsMapPage{News: newsmap}
	t, err := template.ParseFiles("newstracker.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(t.Execute(w, p))
}

func main() {

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)

}
