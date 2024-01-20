package main

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

// Player stores state for a codecrack player and enables two-way communication via websocket
type Player struct {
	Conn      *websocket.Conn
	Username  string
	Number    int
	Combi     [5]Colour
	Board     *Board
	ReadPump  chan Packet[json.RawMessage]
	WritePump chan any
}

// StartReadPump starts the read pump, which allows client websocket packets to be read via the channel
func (p *Player) StartReadPump() {
	for {
		var pack Packet[json.RawMessage]

		// Read next client packet and parse it
		if err := p.Conn.ReadJSON(&pack); err != nil {
			return
		}
		p.ReadPump <- pack
	}
}

// StartWritePump starts the write pump, which enables the server to send packets to the client and avoiding errors if two packets are sent at the same time
func (p *Player) StartWritePump() {
	for {
		pack, ok := <-p.WritePump
		if !ok {
			return
		} else if err := p.Conn.WriteJSON(pack); err != nil {
			return
		}
	}
}

// ReadDeadline waits the provided duration for the client to send a packet, and checks the opcode to see if it is the correct one
func ReadDeadline[T any](p *Player, opcode Opcode, deadline time.Duration) (T, bool) {
	timeout := make(chan struct{})
	var data T

	// Send timeout message if the deadline is reached
	go func() {
		time.Sleep(deadline * time.Second)
		timeout <- struct{}{}
	}()

	select {
	case <-timeout:
		// If no message was received, return false
		return data, false
	case pack := <-p.ReadPump:
		// If message is received in time, parse it and return true
		json.Unmarshal(pack.Data, &data)
		if pack.Opcode != opcode {
			return data, false
		}
		return data, true
	}
}
