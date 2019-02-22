package main

import (
	"io/ioutil"
	"jspfmt"
	"log"
)

func main() {
	input, err := ioutil.ReadFile("test.html")
	if err != nil {
		log.Fatalln(err)
	}
	jspfmt.Format("TEST", string(input))
}
