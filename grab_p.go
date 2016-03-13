package main

import (
  "net/http"
//  "golang.org/x/net/html"
  "io/ioutil"
  "fmt"
  //"bufio"
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

func main() {
  resp, _ := http.Get("http://imgur.com/r/dankmemes")
  bytes, _ := ioutil.ReadAll(resp.Body)
  allHtml := strings.Split(string(bytes), "<p>")
  for i:=0; i<60; i++{
    for j:=0; j<10; j++{
      fmt.Print(allHtml[i][j])
    }
    fmt.Println("")
  }
  //p_string := GetStringInBetween(string(bytes),"<p>","</p>")

  //fmt.Println(p_string)
//  fmt.Println("HTML:\n\n", string(bytes))

  resp.Body.Close()
}
