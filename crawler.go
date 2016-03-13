package crawler

import (
    "time"
    "math/rand"
    "fmt"
    "golang.org/x/net/html"
    "net/http"
    "strings"
)

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
    // Iterate over all of the Token's attributes until we find an "href"
    for _, a := range t.Attr {
        if a.Key == "id" {
            href = a.Val
            ok = true
        }
    }

    // "bare" return will return the variables (ok, href) as defined in
    // the function definition
    return
}

// Extract all http** links from a given webpage
func crawl(url string, ch chan string, chFinished chan bool) {
    resp, err := http.Get(url)

    defer func() {
        // Notify that we're done after this function
        chFinished <- true
    }()

    if err != nil {
        fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
        return
    }

    b := resp.Body
    defer b.Close() // close Body when the function returns

    z := html.NewTokenizer(b)

    for {
        tt := z.Next()

        switch {
        case tt == html.ErrorToken:
            // End of the document, we're done
            return
        case tt == html.StartTagToken:
            t := z.Token()

            // Check if the token is an <a> tag
            isAnchor := t.Data == "div"
            if !isAnchor {
                continue
            }

            // Extract the href value, if there is one
            ok, url := getHref(t)
            if !ok {
                continue
            }

            // Make sure the url begines in http**
      //      len := utf8.RuneCountInString(url)
            hasProto := strings.Index(url, "") == 0
            if hasProto && (len([]rune(url)) == 7) && strings.Compare("section",url) != 0 {
                ch <- url
            }
        }
    }
}

func initCrawl() {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    foundUrls := make(map[string]bool)
    chUrls := make(chan string)
    chFinished := make(chan bool)

    go crawl("http://imgur.com/r/dankmemes", chUrls, chFinished)

    // Subscribe to both channels
    for allFound := false; !allFound; {
        select {
        case url := <-chUrls:
            foundUrls[url] = true
        case <-chFinished:
            allFound = true;
        }
    }
    var allImages[] string
    for url, _ := range foundUrls {
    //    fmt.Println(" - " + url)
        imagesrc := "http://i.imgur.com/"+url+".jpg"
        allImages = append(allImages, imagesrc)
        fmt.Println(imagesrc)
        //TODO: Store image in an array with title?
    //    addImageStore(imagesrc, title)
    //TODO: Find <P> tag for each div, and pull title

    //TODO: Match Meme with closest keywords to title (ContainsAny?)

    }

    close(chUrls)
    return allImages[r.Int()%60]
}
