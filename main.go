package main

import (
	"adviserbot/internal/clients/telegram"
	eventconsumer "adviserbot/internal/consumer/event-consumer"
	telegram2 "adviserbot/internal/events/telegram"
	"adviserbot/internal/storage/postgres"
	"flag"
	"log"
	"os"
	"strconv"
)

var (
	storageConn  = os.Getenv("STORAGE_CON")
	batchSize, _ = strconv.Atoi(os.Getenv("BATCH_SIZE"))
	host         = os.Getenv("HOST")
	token        = os.Getenv("TOKEN")
)

func main() {
	telegramClient := telegram.New(host, token)
	storage, err := postgres.New(storageConn)
	if err != nil {
		panic(err)
	}
	eventFethcer := telegram2.New(telegramClient, storage)

	consumer := eventconsumer.New(eventFethcer, eventFethcer, batchSize)
	log.Println("Starting bot...")
	consumer.Start()
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
