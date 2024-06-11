package eventconsumer

import (
	"adviserbot/internal/events"
	"log"
	"sync"
	"time"
)

type Consumer struct {
	fetcher events.Fethcer

	processor events.Processor

	bathSize int
}

func New(fetcher events.Fethcer, processor events.Processor, bathSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		bathSize:  bathSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.bathSize)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		log.Println("ty ty ty ")
		c.HandleEvents(gotEvents)
	}
}

func (c *Consumer) HandleEvents(events []events.Event) {
	var wg sync.WaitGroup
	for _, event := range events {
		log.Printf("got new event %s\n", event.Text)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := c.processor.Process(event); err != nil {
				log.Printf("can't handle event:%s ", err.Error())
			}
		}()
	}
	wg.Wait()
}
