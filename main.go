package main

import (
	"flag"
	"log"
)

func main() {
	token := mustToken()
	_ = token
	//TODO init client

	//TODO init fetcher

	//TODO init processor

	//TODO init consumer and start consumer
}

func mustToken() string {
	var token string
	flag.StringVar(&token, "token", "", "token for telegram client")
	flag.Parse()
	if token == "" {
		log.Fatal("empty token string")
	}
	return token
}
