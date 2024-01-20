package main

// Opcode distinguishes network actions from either the server or a client
type Opcode int

// Server-to-client packets
const (
	OpPing Opcode = iota
	OpMatchmaking
	OpGameStart
	OpGameEnd
	OpRoundStart
	OpRoundEnd
	OpBoardUpdate
	OpOpponentStatus
	OpForceSubmit
)

// Client-to-server packets
const (
	OpPong Opcode = iota
	OpFindGame
	OpSubmit
)

// Packet represents a websocket packet
type Packet[T any] struct {
	Opcode Opcode `json:"op"`
	Data   T      `json:"d,omitempty"`
}

type ParamsFindGame struct {
	Username string    `json:"username"`
	Combi    [5]Colour `json:"combi"`
}

type ParamsGameStart struct {
	Opponent string `json:"opponent"`
}

type ParamsRoundStart struct {
	RoundNum int `json:"roundNum"`
}

type BoardUpdate struct {
	Board []GuessRow `json:"board"`
}

type ParamsSubmit struct {
	Combi [5]Colour `json:"combi"`
}

type ParamsOpponentStatus struct {
	Status [5]Colour `json:"status"`
}

type ParamsGameEnd struct {
	Win   *bool `json:"win"`
	Error bool  `json:"error"`
}

type ParamsForceSubmit struct {
	Combi [5]Colour `json:"combi"`
}
