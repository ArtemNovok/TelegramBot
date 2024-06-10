package main

import (
	"adviserbot/internal/clients/telegram"
	"flag"
	"log"
)

func main() {
	telegramClient := telegram.New(mustHost(), mustToken())
	_ = telegramClient
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
func mustHost() string {
	var host string
	flag.StringVar(&host, "host", "", "host for telegram client")
	flag.Parse()
	if host == "" {
		log.Fatal("empty host string")
	}
	return host
}
