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
import "code.google.com/p/goauth2/oauth"
import "github.com/gorilla/sessions"
import "strconv"
import "net"

var timeout = time.Duration(30 * time.Second)
func dialWithTimeout (network, addr string) (net.Conn, error) {
    return net.DialTimeout (network, addr, timeout)
}

var transport = http.Transport { Dial: dialWithTimeout }
var client = &http.Client { Transport : &transport}
var store = sessions.NewCookieStore([]byte("tomtom-secret-key"))

func fetchUrl(url string, lastModified string) (string, []byte) {
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

func getUserId (w http.ResponseWriter, r *http.Request) string {
    session, err := store.Get(r, "session")
    if err != nil {
        http.Redirect(w, r, "/", http.StatusFound)
    }
    

    id := session.Values["UserId"].(string)
    
    if len(id) == 0 {
        http.Redirect(w, r, "/", http.StatusFound)
    }
    
    return id
}

func removeFeedHandler(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
    type JsonData struct {
        Id string
    }
    var jsonData JsonData
    json.Unmarshal (body, &jsonData)
    
    feedid := jsonData.Id
    userid := getUserId (w, r)

    db.RemoveFeed (feedid, userid)
    listFeedsHandler(w, r)
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

    feed := data.Feed { Id: data.GenerateId (jsonData.Url), Url : jsonData.Url, Title: jsonData.Url }
    userid := getUserId (w, r)
    
    was_inserted := db.AddFeed (feed, userid)
    
    if was_inserted {
        fetchFeed (feed)
    }
    
    listFeedsHandler(w, r)
}

func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header ().Add ("Content-Type", "application/json")
    userid := getUserId (w, r)
    data, err := json.Marshal (db.GetAllFeedsForUser (userid))
    
    if err != nil {
        panic (err)
    }

    fmt.Fprintf (w, "%s", data)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
    feedId := r.URL.Path[6:]
    offset, err := strconv.Atoi (r.FormValue("o"))

    if (err != nil) {
        panic (err)
    }

    w.Header ().Add ("Content-Type", "application/json")
    data, err := json.Marshal (db.GetFeedItems(feedId, offset))

    if err != nil {
        panic (err)
    }

    fmt.Fprintf (w, "%s", data)
}

func importHandler (w http.ResponseWriter, r *http.Request) {
    file, _, err := r.FormFile ("opmlFile")
    
    if err != nil {
        log.Print (err)
        return
    }
    contents, err := ioutil.ReadAll (file)
    
    if err != nil {
        log.Print (err)
        return
    }
    
    userid := getUserId (w, r)
    urls, err := parser.ParseOPML (contents)
    
    if err != nil {
        log.Print (err)
        return
    }
    
    for _, url := range (urls) {

        feed := data.Feed { Id: data.GenerateId (url), Url : url }
        was_inserted := db.AddFeed (feed, userid)

        if was_inserted {
            fetchFeed (feed)
        }
    }
    
    fmt.Fprintf (w, "%s", "Imported");
}

type Pair struct
{
    Feed data.Feed
    FeedItem data.FeedItem
}

func recentFeedItemsHandler (w http.ResponseWriter, r *http.Request) {
    userid := getUserId (w, r)
    feeds, feedItems := db.GetRecentFeedItems (userid)
    
    pairs := []Pair {}
    for i, feed := range feeds {
        pairs = append (pairs, Pair { feed, feedItems[i] })
    }
    
    w.Header ().Add ("Content-Type", "application/json")
    data, err := json.Marshal (pairs)

    if err != nil {
        panic (err)
    }

    fmt.Fprintf (w, "%s", data)
}

var oauthCfg = &oauth.Config {
    ClientId : "",
    ClientSecret : "",
    AuthURL: "https://accounts.google.com/o/oauth2/auth",
    TokenURL: "https://accounts.google.com/o/oauth2/token",
    RedirectURL: "http://localhost:8080/oauth2callback",
    Scope: "https://www.googleapis.com/auth/userinfo.profile",
}

func authenticationHandler (w http.ResponseWriter, r *http.Request) {
    url := oauthCfg.AuthCodeURL("")
    http.Redirect (w, r, url, http.StatusFound)
}

type User struct {
    Id string
    Name string
    Given_Name string
    Family_Name string
    Link string
    Gender string
    Locale string
}

func oauthCallbackHandler (w http.ResponseWriter, r *http.Request) {
    profileInfoURL := "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"
    code := r.FormValue ("code")
    t := oauth.Transport { Config: oauthCfg }
    t.Exchange (code)
    resp, err := t.Client().Get(profileInfoURL)
    if err != nil { 
        panic (err)
    }
    defer resp.Body.Close()

    user := User {}
    contents, err := ioutil.ReadAll(resp.Body)
    json.Unmarshal (contents, &user)
    
    session, _ := store.Get (r, "session")
    session.Values["UserId"] = user.Id
    session.Values["GivenName"] = user.Given_Name
    log.Printf("** %s Logged in **", user.Name)
    session.Save (r, w)
    
    http.Redirect(w, r, "/view/", http.StatusFound)
}

func initWebServer() {
    http.HandleFunc("/feeds/add", addFeedHandler)
    http.HandleFunc("/feeds/remove", removeFeedHandler)
    http.HandleFunc("/feeds", listFeedsHandler)
    http.HandleFunc("/recent", recentFeedItemsHandler)
    http.HandleFunc("/feed/", feedHandler)
    http.HandleFunc("/import", importHandler)
    http.HandleFunc("/", authenticationHandler)
    http.HandleFunc("/oauth2callback", oauthCallbackHandler)
    http.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("/home/saaadhu/code/git/tomtom/src/tomtom/www"))))
    http.ListenAndServe(":8080", nil)
}

func fetchFeed (feed data.Feed) {
    lastModified, contents := fetchUrl (feed.Url, feed.LastModified)
    
    if len(contents) == 0 {
        return
    }

    title, feedItems, err := parser.Parse (string (contents))
    
    if err != nil {
        log.Printf ("%s", err)
        return
    }

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
        log.Printf("Start fetch")
        fetchFeeds ()
    }
    
}
