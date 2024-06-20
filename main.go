package main

import "log"

func main() {
	server := NewApiServer("127.0.0.1", "9000")

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
