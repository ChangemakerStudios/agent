package swarmpit

import (
	"log"
	"time"
	"bytes"
	"net/http"
	"encoding/json"
	"github.com/swarmpit/agent/setup"
)

var arg = setup.GetArgs()

type EventType string

const (
	EVENT EventType = "event"
	STATS EventType = "stats"
	EMPTY           = ""
	TAB             = "\t"

	// Swarmpit reads the shared secret from this header. It deliberately is not
	// Authorization: swarmpit parses that as a JWT and rejects anything it
	// cannot verify before the event endpoint's own access rule is consulted.
	TOKEN_HEADER = "X-Swarmpit-Event-Token"
	CONTENT_TYPE = "application/json; charset=utf-8"
)

type Event struct {
	EventType EventType   `json:"type"`
	Message   interface{} `json:"message"`
}

func SendEvent(eventType EventType, message interface{}) {
	event := Event{EventType: eventType, Message: message}
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(EMPTY, TAB)
	encoder.Encode(event)

	if eventType == STATS && arg.Debug.Stats == true {
		log.Printf("DEBUG: Host stats: %s", buffer)
	}

	if eventType == EVENT && arg.Debug.Event == true {
		log.Printf("DEBUG: Docker event: %s", buffer)
	}

	request, err := http.NewRequest(http.MethodPost, arg.EventEndpoint, buffer)
	if err != nil {
		log.Printf("ERROR: Event request creation failed: %s", err)
		return
	}
	request.Header.Set("Content-Type", CONTENT_TYPE)
	if arg.EventToken != EMPTY {
		request.Header.Set(TOKEN_HEADER, arg.EventToken)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("ERROR: Event sending failed: %s", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		log.Printf("ERROR: Event rejected: swarmpit returned 401. Check that SWARMPIT_EVENT_TOKEN matches the app's.")
	}
}

func HealthCheck() {
	for {
		<-time.After(5 * time.Second)
		_, err := http.Get(arg.HealthCheckEndpoint)

		if err == nil {
			log.Printf("INFO: Swarmpit OK")
			break;
		}
	}
}
