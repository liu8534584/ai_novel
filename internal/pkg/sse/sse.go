package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EventType string

const (
	EventMessage     EventType = "message"
	EventStateUpdate EventType = "state_update"
	EventError       EventType = "error"
	EventEnd         EventType = "end"
)

type Message struct {
	Event EventType   `json:"event"`
	Data  interface{} `json:"data"`
}

// SetHeaders sets the necessary headers for SSE.
func SetHeaders(c http.ResponseWriter) {
	c.Header().Set("Content-Type", "text/event-stream")
	c.Header().Set("Cache-Control", "no-cache")
	c.Header().Set("Connection", "keep-alive")
	c.Header().Set("Access-Control-Allow-Origin", "*")
}

// Send sends a message to the client.
func Send(w http.ResponseWriter, msg Message) error {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return err
	}
	
	// Format: event: <event>\ndata: <data>\n\n
	fmt.Fprintf(w, "event: %s\n", msg.Event)
	fmt.Fprintf(w, "data: %s\n\n", string(dataBytes))
	
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// SendText sends a simple text message (chunk of generated text).
func SendText(w http.ResponseWriter, text string) error {
	// For simple text, we can send it as a JSON string to keep format consistent
	// Or just raw text if the frontend expects it. 
	// Based on the guide, "receive text stream", let's assume JSON string for safety.
	return Send(w, Message{
		Event: EventMessage,
		Data:  text,
	})
}
