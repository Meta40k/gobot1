package main

import (
	"fmt"
	"os"
)

func main() {

	p := &program{}
	err := p.Initialize(os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERR: %s\n", err)
		os.Exit(1)
	}
	err = p.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERR: %s\n", err)
		os.Exit(1)
	}
	p.wait()

}
