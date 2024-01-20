package main

// Board stores game state for a game of Mastermind
type Board struct {
	combi [5]Colour
	rows  []GuessRow
}

// GuessRow represents a row on a Mastermind board
type GuessRow struct {
	Pieces [5]Colour `json:"pieces"`
	Check  [5]Colour `json:"check"`
}

// NewBoard creates a new Mastermind game board
func NewBoard(combi [5]Colour) *Board {
	b := &Board{combi: combi, rows: []GuessRow{}}
	return b
}

// AddRow checks the player's guess against the true combination and adds the row to the board, it returns a boolean
// value of whether the guess was correct or not as well as the check array
func (b *Board) AddRow(guess [5]Colour) (bool, [5]Colour) {
	row := GuessRow{Pieces: guess}

	// Check for black pegs first
	check := [5]Colour{}
	correct := true
	for i := 0; i < 5; i++ {
		if guess[i] == b.combi[i] {
			check[i] = Colour2
		} else {
			correct = false
		}
	}

	// Add white pegs
	remaining := b.combi
	for i := 0; i < 5; i++ {
		if check[i] == Colour2 {
			continue
		}
		for j := 0; j < 5; j++ {
			if remaining[j] == ColourNone {
				continue
			}
			if remaining[j] == guess[i] && remaining[j] != guess[j] {
				remaining[j] = ColourNone
				check[i] = Colour1
				break
			}
		}
	}

	row.Check = check
	b.rows = append(b.rows, row)

	return correct, check
}

// GetRows returns the board rows
func (b *Board) GetRows() []GuessRow {
	return b.rows
}

// GetStatus returns only the 'status' portion of each row (check pegs)
func (b *Board) GetStatus() [][5]Colour {
	status := [][5]Colour{}
	for _, row := range b.rows {
		status = append(status, row.Check)
	}
	return status
}
