package parcs

import "log"

type Executable interface {
	Init()
	Run()
	Shutdown()
}

func Exec(e Executable) {
	log.Printf("Welcome to PARCS-Go!")
	e.Init()
	log.Printf("Running your program...")
	e.Run()
	e.Shutdown()
	log.Printf("Bye!")
}
