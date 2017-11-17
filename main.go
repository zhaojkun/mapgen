package main

import (
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
)

var (
	file = flag.String("file", "example.go", "")
	typ  = flag.String("type", "Backend", "")
)

func main() {
	flag.Parse()
	buf, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}
	st := parse(string(buf), *typ)
	rs := generator(st)
	formatedSource, err := format.Source([]byte(rs))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(formatedSource))
}
