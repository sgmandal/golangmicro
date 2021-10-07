package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	//a handler function, works on which request is received
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("hello world")
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "oops", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(rw, "hello %s", d)
	})

	//example of such handle function
	//localhost:9090/goodbye
	http.HandleFunc("/goodbye", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("goodbye world")
	})

	//listening and serving function
	//9090 obviously specifies the port number
	//nil specifies no servicemux
	/*servicemux is a multiplexer code working on which codeto execute depending
	on what link or url we received, here defaultservicemux is used which is
	inbuilt in golang package
	*/
	http.ListenAndServe(":9090", nil)
}
