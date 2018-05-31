package main

import "log"

func main() {
	log.Fatal(Serve("0.0.0.0:8888"))
}
