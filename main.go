
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// command line flags
	port := flag.Int("port", 80, "port to serve on")
	//parse all flags
	flag.Parse()

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
