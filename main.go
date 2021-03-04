package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
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

var wg sync.WaitGroup

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

func newsRoutine(c chan News, Location string) {
	defer wg.Done()
	var n News
	resp, _ := http.Get(Location)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &n)
	resp.Body.Close()
	c <- n
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	var s SitemapIndex
	// var n News
	bytes := washPostXML
	xml.Unmarshal(bytes, &s)

	newsmap := make(map[string]NewsMap)
	queue := make(chan News, 50)

	for _, Location := range s.Locations {
		wg.Add(1)
		go newsRoutine(queue, Location)
	}
	wg.Wait()
	close(queue)

	for elem := range queue {
		for idx := range elem.Titles {
			locationTopic := strings.Split(elem.Locations[idx], "/")
			topic := strings.ReplaceAll(locationTopic[3], "-", " ")
			topic = strings.ReplaceAll(topic, "sitemap.xml", "")
			newsmap[elem.Titles[idx]] = NewsMap{topic, elem.Locations[idx]}
		}
	}
	p := NewsMapPage{News: newsmap}
	t, _ := template.ParseFiles("newstracker.html")
	t.Execute(w, p)

}

func main() {

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)

}
