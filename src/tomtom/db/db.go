package db

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "tomtom/data"
    "io/ioutil"
    "os"
    "time"
)

var dataDir string = "/home/saaadhu/code/git/tomtom/data/"


func getConnection() (con *sql.DB) {
    con, err := sql.Open ("mysql", "tomtom_writer:tomtom_writer@/tomtom")
    
    if err != nil {
        panic (err)
    }
    
    return
}

func saveItem (feedId string, feedItemId string, contents string) {
    err := os.MkdirAll(dataDir + feedId, 0777)
    if (err != nil) {
        panic("Could not create dir")
    }
    
    err = ioutil.WriteFile(dataDir + feedId + "/" + feedItemId, []byte(contents), 0777)
    if (err != nil) {
        panic("Could not save file")
    }
}

func readItem (feedId string, feedItemId string) string {
    contents, err := ioutil.ReadFile (dataDir + feedId + "/" + feedItemId)
    
    if (err != nil) {
        panic(err)
    }
    return string (contents)
}

func GetAllFeedsForUser (userid string) []data.Feed {
    con := getConnection()
    defer con.Close()

    rows, err := con.Query ("select feeds.id, feeds.title, feeds.url, server_last_modified, " +
                                "(select count(id) from feeditems " +
                                 "where userfeeds.feed = feeditems.feed and userfeeds.userid=? and " +
                                    "(feeditems.published > (select published from feeditems where id = last_read_item) or userfeeds.last_read_item IS NULL)" +
                                    "group by userfeeds.feed order by feeditems.published desc)," +
                             "last_fetch from feeds, userfeeds where userfeeds.userid=? and feeds.id = userfeeds.feed" , userid, userid)
    
    if err != nil {
        panic (err)
    }
    
    feeds := []data.Feed {}
    for rows.Next() {
        var id, url,title, serverLastModified string
        var lastFetch time.Time
        var unreadItemsCount int
        rows.Scan (&id, &title, &url, &serverLastModified, &unreadItemsCount, &lastFetch)
        feeds = append (feeds, data.Feed { Id : id, Title: title, Url : url, LastModified : serverLastModified, LastFetch : lastFetch, UnreadItemsCount : unreadItemsCount })
    }
    
    return feeds
}


func GetAllFeeds () []data.Feed {
    con := getConnection()
    defer con.Close()

    rows, err := con.Query ("select id, title, url, last_fetch, server_last_modified from feeds")
    
    if err != nil {
        panic (err)
    }
    
    feeds := []data.Feed {}
    for rows.Next() {
        var id, url,title, serverLastModified string
        var lastFetch time.Time
        rows.Scan (&id, &title, &url, &lastFetch, &serverLastModified)
        feeds = append (feeds, data.Feed { Id : id, Title: title, Url : url, LastFetch: lastFetch, LastModified : serverLastModified })
    }
    
    return feeds
}

func GetFeedItems (feedId string) []data.FeedItem {
    con := getConnection()
    defer con.Close()
    
    rows, err := con.Query ("select id, title, url, published from feeditems where feed=? order by published desc", feedId)
    if err != nil {
        panic (err)
    }

    feedItems := []data.FeedItem {}
    for rows.Next() {
        var id, title, url string
        var published time.Time
        rows.Scan (&id, &title, &url, &published)
        
        contents := readItem (feedId, id)
        
        feedItems = append (feedItems, data.FeedItem { Id: id, Title: title, Url: url, Pubdate: published, Contents: contents })
    }
    
    if len(feedItems) > 0 {
        stmt, err := con.Prepare("UPDATE userfeeds SET last_read_item=? WHERE feed=?")

        if err != nil {
            panic (err)
        }

        _, err = stmt.Exec (feedItems[0].Id, feedId)
        if err != nil {
            panic (err)
        }
    }
    
    return feedItems
}

func AddFeed (feed data.Feed, userid string) bool {
    con := getConnection()
    defer con.Close()

    stmt, err := con.Prepare("INSERT IGNORE INTO feeds (id, title, url, last_fetch) VALUES (?,?,?,?)")
    
    if err != nil {
        panic (err)
    }
    
    res, err := stmt.Exec (feed.Id, feed.Title, feed.Url, feed.LastFetch)
    if err != nil {
        panic (err)
    }
    rows_affected, _ := res.RowsAffected()
    row_inserted := rows_affected == 1
    
    stmt, err = con.Prepare ("INSERT IGNORE INTO userfeeds (userid, feed) VALUES (?, ?)")

    if err != nil {
        panic (err)
    }
    
    _, err = stmt.Exec (userid, feed.Id)
    if err != nil {
        panic (err)
    }
    
    return row_inserted
}

func UpdateFeed (feed data.Feed) {
    con := getConnection()
    defer con.Close()

    stmt, err := con.Prepare("UPDATE feeds SET title=?,last_fetch=?, server_last_modified=? WHERE id=?")
    
    if err != nil {
        panic (err)
    }

    _, err = stmt.Exec (feed.Title, feed.LastFetch, feed.LastModified, feed.Id)
    if err != nil {
        panic (err)
    }
}

func SaveFeedItem (feed data.Feed, item data.FeedItem) {
    con := getConnection()
    defer con.Close()

    stmt, err := con.Prepare ("INSERT INTO feeditems (id, feed, title, url, blurb, published) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE title=VALUES(title), url=VALUES(url), blurb=VALUES(blurb), published=VALUES(published)")
    
    if err != nil { 
        panic (err)
    }

    _, err = stmt.Exec (item.Id, feed.Id, item.Title, item.Url, item.Blurb, item.Pubdate)
    if err != nil {
        panic (err)
    }

    saveItem (feed.Id, item.Id, item.Contents)
}
