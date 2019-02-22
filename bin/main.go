package main

import (
	"io/ioutil"
	"log"

	"github.com/awmottaz/jspfmt/jspfmt"
)

func main() {
	input, err := ioutil.ReadFile("test.html")
	if err != nil {
		log.Fatalln(err)
	}
	jspfmt.Format("TEST", string(input))
}
