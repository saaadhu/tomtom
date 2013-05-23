package main

import "fmt"
import "crypto/md5"
import "io"
import "io/ioutil"
import "net/http"
import "time"
import "os"
/*
import "tomtom/db"
*/
import "tomtom/data"
import "tomtom/parser"
var dataDir string = "/home/saaadhu/code/git/tomtom/data/"

func hashUrl(url string) string {
    h := md5.New()
    io.WriteString(h, url)
    return fmt.Sprintf("%x", h.Sum([]byte{}));
}

func fetchUrl(url string) []byte {
    res, err := http.Get(url);
    if err != nil {
        panic("Couldn't fetch URL")
    }
    
    defer res.Body.Close()
    
    contents, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic("Couldn't read contents")
    }
    return contents
}

func handle (feed *data.Feed, contents string) {
    
}

func save (urlhash string, contents []byte) {
    err := os.MkdirAll(dataDir + urlhash, 0777)
    if (err != nil) {
        panic("Could not create dir")
    }
    
    filename := time.Now().Format(time.RFC3339)
    err = ioutil.WriteFile(dataDir + urlhash + "/" + filename, contents, 0777)
    if (err != nil) {
        panic("Could not save file")
    }
}

func main() {
    text := `
<?xml version="1.0" encoding="utf-8"?>
<rss xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
    <channel>
        <atom:link href="http://www.codinghorror.com/blog/index.xml" rel="self" type="application/rss+xml" />
        <title>Coding Horror</title>
        <link>http://www.codinghorror.com/blog/</link>
        <description>programming and human factors - Jeff Atwood</description>
        <language>en-us</language>
      
        <lastBuildDate>Mon, 29 Apr 2013 16:45:34 -0700</lastBuildDate>
        <pubDate>Mon, 29 Apr 2013 16:45:34 -0700</pubDate>
        <generator>http://www.typepad.com/</generator>
        <docs>http://blogs.law.harvard.edu/tech/rss</docs>
        
        <image>
            <title>Coding Horror</title>
            <url>http://www.codinghorror.com/blog/images/coding-horror-official-logo-small.png</url>
            <width>100</width>
            <height>91</height>
            <description>Logo image used with permission of the author. (c) 1993 Steven C. McConnell. All Rights Reserved.</description>
            <link>http://www.codinghorror.com/blog/</link>
        </image>
        
        <xhtml:meta xmlns:xhtml="http://www.w3.org/1999/xhtml" name="robots" content="noindex" />

    
        <item>
            <title>So You Don&#39;t Want to be a Programmer After All</title>
            <link>http://www.codinghorror.com/blog/2013/04/so-you-dont-want-to-be-a-programmer-after-all.html</link>
            <description><![CDATA[<p>
I get a surprising number of emails from career programmers who have spent some time in the profession and eventually decided it just isn't for them. Most recently this:
</p>
I've seen less "adept" programmers self-select into related roles at previous jobs and do very well, both financially and professionally. There is a <i>lot</i> of stuff that goes on around programming that is not heads down code writing, where your programming skills are a competitive advantage.
</p>
</table>]]></description>
            <guid>http://www.codinghorror.com/blog/2013/04/so-you-dont-want-to-be-a-programmer-after-all.html</guid>
            <pubDate>Mon, 29 Apr 2013 16:45:34 -0700</pubDate>
        </item>
    </channel>
</rss>`

    for _, feedItem := range parser.Parse(text) {
        fmt.Println (feedItem)
    }
    /*
    for _,feed := range db.GetFeeds() {
        url := feed.Url
        id := feed.Id
        contents := fetchUrl(url);
    }
    */
}
