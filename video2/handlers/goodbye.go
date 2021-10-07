package handlers

import (
	"log"
	"net/http"
)

type GoodBye struct {
	l *log.Logger
}

func NewBye(x *log.Logger) *GoodBye {
	return &GoodBye{x}
}

func (g *GoodBye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	g.l.Println("goodbye")
	rw.Write([]byte("bye"))
}
