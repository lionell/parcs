package parcs

import "log"

type Executable interface {
	Init()
	Run()
	Shutdown()
}

func Exec(e Executable) {
	e.Init()
	log.Printf("Initialized successfully")
	e.Run()
	e.Shutdown()
	log.Printf("Shut down successfully")
}
