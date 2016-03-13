package main

import (
  "net/http"
  "io/ioutil"
  "fmt"
  "strings"
  )

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

//TODO: Rename this later.
func main() {
  resp, _ := http.Get("http://imgur.com/r/dankmemes")
  bytes, _ := ioutil.ReadAll(resp.Body)
  modified_string := string(bytes)

//All titles array
  var allTitles[] string

for i := 0; i < 61; i++ {
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

  fmt.Println(allTitles)
  resp.Body.Close()
}
