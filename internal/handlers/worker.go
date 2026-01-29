package handlers

import (
	"log"
	"time"
)

var EmailQueue = make(chan string, 100)

func StartEmailWorker() {
	go func() {
		for email := range EmailQueue {
			time.Sleep(2 * time.Second)
			log.Printf("[Async Worker] Confirmation email sent to %s", email)
		}
	}()
}
