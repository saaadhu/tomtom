package parser

import "encoding/xml"
import "tomtom/data"
import "time"
//import "fmt"

type item struct {
    XMLName xml.Name `xml:"item"`
    Title string    `xml:"title"`
    Link string     `xml:"link"`
    Description string `xml:"description"`
    Guid string `xml:"guid"`
    PubDate string `xml:"pubDate"`
}

type channel struct {
    XMLName xml.Name `xml:"channel"`
    Title string `xml:"title"`
    Link string `xml:"link"`
    Items []item `xml:"item"`
}

type rss struct {
    XMLName xml.Name `xml:"rss"`
    Channel channel `xml:"channel"`
}

func Parse(contents string) []data.FeedItem {
    r := rss {}
    feedItems := []data.FeedItem {}
    err := xml.Unmarshal([]byte(contents), &r)

    if err != nil {
        panic (err)
    }
    channel := r.Channel
    
    for _,item := range channel.Items {
       t, err := time.Parse ("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
       if err != nil {
           panic (err)
       }
       feedItem := data.FeedItem { item.Guid, item.Title, item.Link, "", item.Description, t } 
       feedItems = append (feedItems, feedItem)
    }

    return feedItems
}
