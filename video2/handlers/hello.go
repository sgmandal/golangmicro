package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//creating a struct
type Hello struct {
	l *log.Logger
}

//receiving the logger
//initializing the struct to a given variable in main program
func NewHello(x *log.Logger) *Hello {
	return &Hello{x}
}

//this func implements an interface inside http package
//rw is an interface hence such declaration
//r's declaration is such because its a struct instance
//and in go it is good practice to create an instance of struct as a pointer
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	//same as log.logger
	//but here we've created our own logger linked to that interface
	//hence h.l...
	h.l.Println("hello world")

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "oops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "hello %s", d)
}
