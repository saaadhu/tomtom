package db

/*
import "database/sql"
import "github.com/go-sql-driver/mysql"
*/
import "tomtom/data"


func GetFeeds() []data.Feed {
    feeds := []data.Feed { data.Feed { Id: "ed29f102abda336b49dc9735ebfd919e", Url:"http://www.codinghorror.com/blog/index.xml" } }
    return feeds
}
