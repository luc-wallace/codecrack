package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MatchPool manages matchmaking and creating new matches
type MatchPool struct {
	active      bool
	playerQueue chan *Player
}

// QueuePlayer initialises a player object from the websocket connection and adds them to the matchmaking queue
func (p *MatchPool) QueuePlayer(conn *websocket.Conn) {
	newPlayer := &Player{
		Conn:      conn,
		ReadPump:  make(chan Packet[json.RawMessage], 1),
		WritePump: make(chan any, 1),
	}

	// Start the websocket pumps
	go newPlayer.StartWritePump()
	go newPlayer.StartReadPump()

	fmt.Printf("queued player: %s\n", conn.RemoteAddr())

	// Close connection if client doesn't respond with details in 10s
	data, ok := ReadDeadline[ParamsFindGame](newPlayer, OpFindGame, 10)
	if !ok {
		conn.Close()
		return
	}

	newPlayer.Combi = data.Combi
	newPlayer.Username = data.Username

	// Send matchmaking packet so client knows it's in queue
	newPlayer.WritePump <- Packet[any]{Opcode: OpMatchmaking}

	// Add player to queue
	p.playerQueue <- newPlayer
}

// Matchmake starts the matchmaking loop that creates a match when at least 2 players are in the matchmaking queue
func (p *MatchPool) Matchmake() {
	for p.active {
		// Receive first two players from the matchmaking queue and initialise match
		p1, p2 := <-p.playerQueue, <-p.playerQueue
		newMatch := &Match{p1, p2, false, sync.Mutex{}}

		// Assign player numbers to distinguish them later
		p1.Number = 1
		p2.Number = 2

		fmt.Printf("new match started: %s and %s\n", p1.Conn.RemoteAddr(), p2.Conn.RemoteAddr())

		// Start new match
		go newMatch.Start()
	}
}

// Match manages the sequencing and networking of a codecrack match
type Match struct {
	p1         *Player
	p2         *Player
	closing    bool
	closeMutex sync.Mutex
}

// Start begins a codecrack match, communicating with both clients from start to end and managing game logic
func (m *Match) Start() {
	// Initialise boards using each player's unique colour combination
	m.ExecuteSync(func(p *Player) {
		var n string
		if p.Number == 1 {
			n = m.p2.Username
			m.p1.Board = NewBoard(m.p2.Combi)
		} else {
			n = m.p1.Username
			m.p2.Board = NewBoard(m.p1.Combi)
		}

		// Send packet so each client knows the gmae has started
		p.WritePump <- Packet[any]{Opcode: OpGameStart, Data: ParamsGameStart{Opponent: n}}
	})

	// If clients don't respond in 10s, cancel the game
	m.ExecuteSync(func(p *Player) {
		_, ok := ReadDeadline[any](p, OpPong, 10)
		if !ok {
			m.End(0)
			return
		}
	})

	if m.closing {
		return
	}

	// Main game loop
	r := 1
	win := [2]bool{}
	for {
		m.ExecuteSync(func(p *Player) {
			// Store this player's opponent
			var opp *Player
			if p.Number == 1 {
				opp = m.p2
			} else {
				opp = m.p1
			}

			// Send round start packet to client
			p.WritePump <- Packet[any]{Opcode: OpRoundStart, Data: ParamsRoundStart{RoundNum: r}}

			// Wait for client to send their row guess within 30s
			data, ok := ReadDeadline[ParamsSubmit](p, OpSubmit, 30)
			if !ok {
				// If client doesn't respond in 30s, send ForceSubmit packet to receive whatever the player entered
				go func() {
					p.WritePump <- Packet[any]{Opcode: OpForceSubmit}
					time.Sleep(100 * time.Millisecond)
				}()
				// If client still doesn't respond, end the game
				data, ok = ReadDeadline[ParamsSubmit](p, OpSubmit, 3)
				if !ok {
					m.End(opp.Number)
					return
				}
			}
			if m.closing {
				return
			}

			// Add the player's row to the board
			correct, check := p.Board.AddRow(data.Combi)

			// If the guess row was the opponent's combination, set win to true
			if correct {
				win[p.Number-1] = true
			}

			// Send board update to player and send their status to opponent
			p.WritePump <- Packet[any]{Opcode: OpBoardUpdate, Data: p.Board.GetRows()}
			opp.WritePump <- Packet[any]{Opcode: OpOpponentStatus, Data: ParamsOpponentStatus{Status: check}}
		})

		if m.closing {
			return
		}

		// Send round end packet to both clients
		m.ExecuteSync(func(p *Player) {
			p.WritePump <- Packet[any]{Opcode: OpRoundEnd}
		})

		// Check if either player won, and end the game
		if win[0] || win[1] {
			winner := 0
			if win[0] && !win[1] {
				winner = 1
			} else if win[1] && !win[0] {
				winner = 2
			}
			m.End(winner)
			return
		}

		r += 1
	}
}

// Pointer is a utility function that creates a pointer from any data type
func Pointer[T any](d T) *T {
	return &d
}

// ExecuteSync runs a function concurrently on both clients in a match, and doesn't return until both functions have returned
func (m *Match) ExecuteSync(cb func(*Player)) {
	doneP1 := make(chan bool)
	doneP2 := make(chan bool)

	go func() {
		cb(m.p1)
		doneP1 <- true
	}()

	go func() {
		cb(m.p2)
		doneP2 <- true
	}()

	<-doneP1
	<-doneP2
}

// End closes a match, sends game end status to both clients and closes both player connections
func (m *Match) End(winner int) {
	// End the closing sequence if the match is already closing
	if m.closing {
		return
	}
	// Ensure that two End functions cannot run at the same time
	m.closing = true
	m.closeMutex.Lock()
	defer m.closeMutex.Unlock()

	m.ExecuteSync(func(p *Player) {
		params := ParamsGameEnd{}
		// Check if player was winner or loser
		if winner == 0 {
			params.Win = nil
		} else if p.Number == winner {
			params.Win = Pointer(true)
		} else {
			params.Win = Pointer(false)
		}

		// Send game end packet
		p.WritePump <- Packet[any]{Opcode: OpGameEnd, Data: params}
		time.Sleep(100 * time.Millisecond)

		// Close websocket pumps as game is over
		close(p.WritePump)
		close(p.ReadPump)
	})

	fmt.Printf("match result %d: %s and %s", winner, m.p1.Conn.RemoteAddr(), m.p2.Conn.RemoteAddr())

	// Close both websocket connections
	m.p1.Conn.Close()
	m.p2.Conn.Close()
}
