
package main

import (
	"time"
	"math/rand"
	"golang.org/x/net/html"
	"strings"
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/websocket"
)

var connections map[*websocket.Conn]bool

func main() {
	port := flag.Int("port", 4000, "")
	dir := flag.String("indexLocation", "", "")

	flag.Parse()
	connections = make(map[*websocket.Conn]bool)

	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)

	http.HandleFunc("/sock", wsHandler)
	log.Printf("PolyHack server started on port %d\n", *port)

	addr := fmt.Sprintf("localhost:%d", *port)
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	_, ok := err.(websocket.HandshakeError)
	if ok {
		http.Error(w, "error", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	connections[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(connections, conn)
			conn.Close()
			return
		}
		log.Println(string(msg))
		message := strings.Split(string(msg), ": ")
		log.Println(message[0])
		log.Println(message[1])
		//DODO -- Go to InitTitle, from there we check the text and load the image.
		imgUrl := initGetTitle(message[1])
		sendAll([]byte(message[0]+": "+imgUrl))
	}
}

func sendAll(message []byte){
	for conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			delete(connections, conn)
			conn.Close()
		}
	}
}

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

func initCrawl(index int) string{
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

    }

    close(chUrls)
    return allImages[index]
}

//Gets title using index (0-60), same position as images array
func initGetTitle(s_input string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

  resp, _ := http.Get("http://imgur.com/r/dankmemes")
  bytes, _ := ioutil.ReadAll(resp.Body)
  modified_string := string(bytes)

//All titles array
  var allTitles[] string

for i := 0; i < 60; i++ {
    if strings.ContainsAny(modified_string,"<p>") == false {
      break
    }

    p_string := GetStringInBetween(modified_string,"<p>","</p>")

    //Store titles into array (slices) -- Removes invalid first <p>, index should match
    if strings.Contains(p_string,"Optimizing your large GIFs...") == false {
      allTitles = append(allTitles, p_string)
    }

    tag := strings.Split(modified_string, "</p>")
    modified_string = strings.Trim(modified_string,tag[0])
  }

  //Check if entered text ContainsAny any element found in array.
  for i := 0; i < len(allTitles); i++ {
    if strings.Contains(allTitles[i], s_input) == true {

			//Crawl images after check
			return initCrawl(i)
      break
    }
  }

  resp.Body.Close()

	//If nothing found, return rand
	return initCrawl(r.Int()%60)
}

//Does what it is titled.
func GetStringInBetween(content, start, end string) (result string) {
	if content != "" && start != "" && end != "" {
		r := strings.Split(content, start)

		if r[1] != "" {
			r = strings.Split(r[1], end)
		}

		result = r[0]
		return
	} else {
		return
	}
}
