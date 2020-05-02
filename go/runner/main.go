package main

import (
	"github.com/lionell/parcs/go/parcs"
	"log"
	"os"
	"strconv"
)

type Program struct {
	*parcs.Runner
}

func (h *Program) Run() {
	n, err := strconv.Atoi(os.Getenv("N"))
	if err != nil {
		log.Fatal(err)
	}
	t, err := h.Start("lionell/factor-py")
	if err != nil {
		log.Fatal(err)
	}
	if err := t.SendAll(n, 1, n+1); err != nil {
		log.Fatal(err)
	}
	var facts []int
	if err := t.Recv(&facts); err != nil {
		log.Fatal(err)
	}
	log.Printf("Factors found %v", facts)
	t.Shutdown()
}

func main() {
	parcs.Exec(&Program{parcs.DefaultRunner()})
}
