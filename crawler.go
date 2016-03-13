//TODO: webcrawler code
package main
import (
  "net/http"
  "golang.org/x/net/html"
  "io/ioutil"
  "fmt"

  "bufio"
//  "os"

  "strings"
  )

// Get Tag Attributes (Image Title, Description)
func getTitle(t html.Token) (ok bool, alt string) {

  for _, img := range t.Attr {
    if img.Key == "alt" {
      alt = img.Val
      ok = true
    }
  }

  return
}


func main() {
  //
  GetSite()
  // foundUrls := make(map[string]bool)
  // seedUrls := os.Args[1:]
  //
  // // Channels
  // chUrls := make(chan string)
  // chFinished := make(chan bool)
  //
  // // Kick off the crawl process (concurrently)
  // for _, url := range seedUrls {
  //     go crawl(url, chUrls, chFinished)
  // }
  //
  // // Subscribe to both channels
  // for c := 0; c < len(seedUrls); {
  //     select {
  //     case url := <-chUrls:
  //         foundUrls[url] = true
  //     case <-chFinished:
  //         c++
  //     }
  // }

  // We're done! Print the results...

  // fmt.Println("\nFound", len(foundUrls), "unique urls:\n")
  //
  // for url, _ := range foundUrls {
  //     fmt.Println(" - " + url)
  // }
  //
  // close(chUrls)
}

//Crawl Function
func crawl ( url string, ch chan string, chFinished chan bool ) {
  resp, err := http.Get(url)

  defer func() {
    //Notifies that we are finished after function
    chFinished <- true
  }()

  if err != nil {
    fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
    return
  }

  b := resp.Body
  defer b.Close() // close Body when the function returns

  // Parse HTML (Find img anchor tag)
  z := html.NewTokenizer(b)
  for {
    tt := z.Next()
    switch {
    case tt == html.ErrorToken:
      // End of the document, we're done
      return
    case tt == html.StartTagToken:
      t := z.Token()

      isAnchor := t.Data == "img"

      if !isAnchor {
        continue
      }

      //Extract the 'alt' value, if found
      ok, url := getTitle(t)
      if !ok {
        continue
      }

      //Check URL for 'http' at beginning
      hasProto := strings.Index(url, "http") == 0
      if hasProto {
        ch <- url
      }
    }
  }
}

// func main() {
//   //
//   foundUrls := make(map[string]bool)
//   seedUrls := os.Args[1:]
//
//   // Channels
//   chUrls := make(chan string)
//   chFinished := make(chan bool)
//
//   // Kick off the crawl process (concurrently)
//   for _, url := range seedUrls {
//       go crawl(url, chUrls, chFinished)
//   }
//
//   // Subscribe to both channels
//   for c := 0; c < len(seedUrls); {
//       select {
//       case url := <-chUrls:
//           foundUrls[url] = true
//       case <-chFinished:
//           c++
//       }
//   }
//
//   // We're done! Print the results...
//   fmt.Println("\nFound", len(foundUrls), "unique urls:\n")
//
//   for url, _ := range foundUrls {
//
//       fmt.Println(" - " + url)
//   }
//
//   close(chUrls)
// }


func main() {
GetSite()
}

//Sample
func GetSite() {
  resp, _ := http.Get("http://imgur.com/r/dankmemes");
  //bytes, _ := ioutil.ReadAll(resp.Img)

  // for info, _ := range bytes {
  //   if strings.Contains(string(bytes), "<img") {
  //     fmt.Println(" - ", info)
  //   }
  // }
  // fmt.Println("HTML:\n\n", string(bytes))
  //
  // resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      // err
  }

  r, err := bufio.NewReader(bytes.NewReader(body), resp.ContentLength)
  if err != nil {
      // err
  }
}
