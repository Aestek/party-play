package main

import (
	"flag"
	"log"
)

func main() {
	addr := flag.String("addr", "0.0.0.0:8888", "addr")
	flag.Parse()

	log.Fatal(Serve(*addr))
}
