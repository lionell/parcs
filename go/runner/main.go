package main

import (
	"github.com/lionell/parcs/go/parcs"
	"log"
)

type Bla struct {
	*parcs.Runner
}

func (h *Bla) Run() {
	t, err := h.Start("lionell/factor-py")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sending data over")
	if err := t.SendAll(100, 1, 100); err != nil {
		log.Fatal(err)
	}
	log.Printf("Sent successfully")
	var facts []int
	if err := t.Recv(&facts); err != nil {
		log.Fatal(err)
	}
	log.Printf("Factors found %v", facts)
	t.Shutdown()
}

func main() {
	parcs.Exec(&Bla{parcs.DefaultRunner()})
}
