package main

import (
	"github.com/lionell/parcs/go/parcs"
	"log"
)

type Echo struct {
	*parcs.Service
}

func (e *Echo) Run() {
	var n int
	log.Printf("BEFORE RECEIVE")
	e.Recv(&n)
	log.Printf("Received %v", n)
	e.Send(n)
}

func main() {
	parcs.Exec(&Echo{parcs.DefaultService()})
}
