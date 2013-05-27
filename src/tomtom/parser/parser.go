package parser

import "encoding/xml"
import "tomtom/data"
import "time"
import "log"
//import "fmt"

type item struct {
    XMLName xml.Name `xml:"item"`
    Title string    `xml:"title"`
    Link string     `xml:"link"`
    Description string `xml:"description"`
    EncodedContent string `xml:"encoded"`
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

type link struct {
    Href string `xml:"href,attr"`
}

type entry struct {
    XMLName xml.Name `xml:"entry"`
    Title string    `xml:"title"`
    Link link     `xml:"link"`
    Contents string `xml:"content"`
    Id string `xml:"id"`
    Published string `xml:"published"`
    Updated string `xml:"updated"`
}

type feed struct {
    XMLName xml.Name `xml:"feed"`
    Title string `xml:"title"`
    Link link `xml:"link"`
    Entries []entry `xml:"entry"`
}

func parseTime (contents string) (time.Time, error) {

    if len(contents) == 0 {
        return time.Now(), nil
    }
    
    timeFormats := []string { "Mon, 02 Jan 2006 15:04:05 -0700", "Mon, 02 Jan 2006 15:04:05 GMT", time.RFC3339 }
    var last_error error
    
    for _, format := range timeFormats {
        t, err := time.Parse (format, contents)
        if err == nil {
            return t, nil
        }
        last_error = err
    }
    
    return time.Now(), last_error
}

func parseRSS (contents string) (string, []data.FeedItem, error) {
    r := rss {}
    feedItems := []data.FeedItem {}
    err := xml.Unmarshal([]byte(contents), &r)

    if err != nil {
        log.Printf("Parsing with RSS failed : %s", err)
        return "", []data.FeedItem{}, err
    }
    channel := r.Channel
    currentTime := time.Now()
    var t time.Time
    
    for i, item := range channel.Items {
        if len(item.PubDate) == 0 {
            t = currentTime.Add(-time.Duration(i) * time.Second)
        } else {
            t, err = parseTime (item.PubDate)
        }
       
       if err != nil {
         panic (err)
       }
       
       if len(item.EncodedContent) > len(item.Description) {
           item.Description = item.EncodedContent
       }

       blurb_length := len(item.Description)
       if blurb_length > 240 {
           blurb_length = 240
       }
       
       id := item.Guid
       if len(id) == 0 {
           id = item.Title
       }
       feedItem := data.FeedItem { data.GenerateId(id), item.Title, item.Link, item.Description[:blurb_length] + "...", item.Description, t } 
       feedItems = append (feedItems, feedItem)
    }

    return channel.Title, feedItems, nil
}

func parseFeed (contents string) (string, []data.FeedItem, error) {
    r := feed {}
    feedItems := []data.FeedItem {}
    err := xml.Unmarshal([]byte(contents), &r)

    if err != nil {
        return "", []data.FeedItem{}, err
    }
    
    currentTime := time.Now()
    var t time.Time
    
    for i, entry := range r.Entries {
       if len(entry.Published) == 0 {
           entry.Published = entry.Updated
       }
        if len(entry.Published) == 0 {
            t = currentTime.Add(-time.Duration(i) * time.Second)
        } else {
            t, err = parseTime (entry.Published)
            if err != nil {
                panic (err)
            }
        }
       feedItem := data.FeedItem { data.GenerateId(entry.Id), entry.Title, entry.Link.Href, entry.Contents[:240] + "...", entry.Contents, t } 
       feedItems = append (feedItems, feedItem)
    }

    return r.Title, feedItems, nil
}

func Parse(contents string) (string, []data.FeedItem) {
    title, feedItems, err := parseRSS (contents)

    if (err == nil) {
        return title, feedItems
    }

    title, feedItems, err = parseFeed (contents)
    if (err != nil) {
        panic (err)
    }

    return title, feedItems
}
