package task

import (
	"io"
	"log"
	"context"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/swarmpit/agent/swarmpit"
)

func HandleEvents(cli *client.Client) {
	messages, errs := cli.Events(context.Background(), events.ListOptions{})

loop:
	for {
		select {
		case err := <-errs:
			if err != nil && err != io.EOF {
				log.Printf("ERROR: Event channel error: %s", err)
			}
			break loop
		case msg, ok := <-messages:
			if !ok {
				log.Printf("ERROR: Event channel closed.")
				break loop
			}
			swarmpit.SendEvent(swarmpit.EVENT, msg)
		}
	}
	panic("Event collector is broken. Shutdown!!!")
}
