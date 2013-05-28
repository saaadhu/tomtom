package data
import "time"
import "crypto/md5"
import "io"
import "fmt"

type FeedItem struct {
    Id string
    Title string
    Url string
    Blurb string
    Contents string
    Pubdate time.Time
}

type Feed struct {
    Id string
    Title string
    Url  string
    LastFetch time.Time
    LastModified string
}

func GenerateId(str string) string {
    h := md5.New()
    io.WriteString(h, str)
    return fmt.Sprintf("%x", h.Sum([]byte{}));
}

