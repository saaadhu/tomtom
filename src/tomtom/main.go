package main

import "fmt"
import "crypto/md5"
import "io"
import "io/ioutil"
import "net/http"
import "time"
import "os"

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
    url := "http://www.codinghorror.com/blog/index.xml"
    hash := hashUrl(url);
    contents := fetchUrl(url);
    save(hash, contents)
}
