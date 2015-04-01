package main

import (
	"github.com/hallas/stacko"
	"github.com/hallas/stacko/test/lib"
	"log"
)

func main() {
	defer func() {
		recover()

		_, err := stacko.NewStacktrace(0)
		if err != nil {
			log.Fatal(err)
		}
	}()

	lib.WillError()
}
