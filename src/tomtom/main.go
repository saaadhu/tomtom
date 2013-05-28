package main

import "fmt"
import "io/ioutil"
import "tomtom/db"
import "tomtom/parser"
import "tomtom/data"
import "net/http"
import "encoding/json"
import "time"
import "log"

var client = &http.Client {}

func fetchUrl(url string, lastModified string) (string, []byte) {
    log.Printf("Fetching %s", url)
    req, _ := http.NewRequest ("GET", url, nil)

    if len(lastModified) != 0 {
        req.Header.Add ("If-Modified-Since", lastModified)
    }
    res, err := client.Do (req);
    if err != nil {
        log.Print(err)
        return "", []byte{}
    }
    defer res.Body.Close()
    
    if res.StatusCode == 304 {
        return "", []byte{}
    }
    
    contents, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic("Couldn't read contents")
    }
    return res.Header.Get("Last-Modified"), contents
}

func addFeedHandler(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
    type JsonData struct {
        Url string
    }
    var jsonData JsonData
    json.Unmarshal (body, &jsonData)

    feed := data.Feed { Id: data.GenerateId (jsonData.Url), Url : jsonData.Url }
    db.AddFeed (feed)
    fetchFeed (feed)
    
    listFeedsHandler(w, r)
}

func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header ().Add ("Content-Type", "application/json")
    data, err := json.Marshal(db.GetAllFeeds())
    
    if err != nil {
        panic (err)
    }

    fmt.Fprintf (w, "%s", data)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
    feedId := r.URL.Path[6:]
    w.Header ().Add ("Content-Type", "application/json")
    data, err := json.Marshal (db.GetFeedItems(feedId))

    if err != nil {
        panic (err)
    }

    fmt.Fprintf (w, "%s", data)
}

func initWebServer() {
    http.HandleFunc("/feeds/add", addFeedHandler)
    http.HandleFunc("/feeds", listFeedsHandler)
    http.HandleFunc("/feed/", feedHandler)
    http.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("/home/saaadhu/code/git/tomtom/src/tomtom/www"))))
    http.ListenAndServe(":8080", nil)
}

func fetchFeed (feed data.Feed) {
    lastModified, contents := fetchUrl (feed.Url, feed.LastModified)
    
    if len(contents) == 0 {
        return
    }

    title, feedItems := parser.Parse (string (contents))

    feed.Title = title
    feed.LastFetch = time.Now()
    feed.LastModified = lastModified

    db.UpdateFeed (feed)
    for _, feedItem := range feedItems {
        db.SaveFeedItem (feed, feedItem)
    }
}

func fetchFeeds() {
    for _,feed := range db.GetAllFeeds() {
        fetchFeed (feed)
    }
}


func main() {
    go initWebServer()
    
    for ;; {
        time.Sleep (15 * time.Minute)
        fetchFeeds ()
    }
    
}
