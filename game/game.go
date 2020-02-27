package game

import (
	"fmt"
)

type BoardState int

const (
	PLAYING BoardState = iota
	X_WON
	O_WON
	TIE
)

type Player int

const (
	NO_PLAYER Player = iota
	PLAYER_X
	PLAYER_O
)

func (p Player) String() string {
	switch p {
	case PLAYER_O:
		return "O"
	case PLAYER_X:
		return "X"
	default:
	case NO_PLAYER:
		return "-"
	}
	return "" // shut the fuck up
}

// Board represents tic tac toe game state
// board is a square grid of cells
type Board struct {
	Size     int
	Grid     [][]Player
	NextTurn Player
	State    BoardState
}

// MakeBoard creates a board of size rows and size columns
func MakeBoard(size int) *Board {
	grid := make([][]Player, size)
	for i := 0; i < size; i++ {
		row := make([]Player, size)
		for j := 0; j < size; j++ {
			row[j] = NO_PLAYER
		}
		grid[i] = row
	}
	return &Board{size, grid, PLAYER_X, PLAYING}
}

// MakeMove for the board at x, y to the given owner
// return error if the cell already has a owner distinct from NONE
// or the coordinates are outside of the grid
func (b *Board) MakeMove(player Player, x, y int) error {
	if player == NO_PLAYER {
		return fmt.Errorf("Only X and O players can make moves")
	}
	if player != b.NextTurn {
		return fmt.Errorf("Cannot make a move: it's not your turn, it's player %s", b.NextTurn)
	}
	err := b.ValidateCoordinates(x, y)
	if err != nil {
		return err
	}
	if b.Grid[y][x] != NO_PLAYER {
		return fmt.Errorf("%d, %d is already occupied", x, y)
	}
	if b.State != PLAYING {
		return fmt.Errorf("Cannot make a move: the game is over")
	}
	b.Grid[y][x] = player
	b.State = b.calcBoardState()
	if b.State != PLAYING {
		b.NextTurn = NO_PLAYER
	} else if player == PLAYER_X {
		b.NextTurn = PLAYER_O
	} else {
		b.NextTurn = PLAYER_X
	}
	return nil
}

// GetValue returns value of the board at x, y
func (b *Board) GetValue(x, y int) (Player, error) {
	err := b.ValidateCoordinates(x, y)
	if err != nil {
		return 0, err
	}
	return b.Grid[y][x], nil
}

func (b *Board) GetState() BoardState {
	return b.State
}

func (b *Board) GetNextTurn() Player {
	return b.NextTurn
}

// return nil if x, y coordinates are valid for this board
// otherwise return error
func (b *Board) ValidateCoordinates(x, y int) error {
	var err error
	if x < 0 || x >= b.Size {
		err = fmt.Errorf("X coordinate must be within [0, %d]", b.Size-1)
	}
	if y < 0 || y >= b.Size {
		err = fmt.Errorf("Y coordinate must be within [0, %d], %s", b.Size-1, err)
	}
	return err
}

func (b *Board) calcBoardState() BoardState {
	state := TIE
	for row := 0; row < b.Size; row++ {
		state = b.getLineState(0, row, 1, 0, state)
	}
	for col := 0; col < b.Size; col++ {
		state = b.getLineState(col, 0, 0, 1, state)
	}
	// check main diagonal
	state = b.getLineState(0, 0, 1, 1, state)
	// check secondary diagonal
	state = b.getLineState(b.Size-1, 0, -1, 1, state)
	return state
}

// wrapper around calcLineState that allows passing previously calculated state (in some other line)
// In case there is already a line with a WON state, there is no point in calculating state of
// any other line, we can just return that state
// If there is a line with PLAYING state, there could still potentially be a WON line, so we need to
// check for that. If we can't find a WON, we stay with PLAYING
func (b *Board) getLineState(x, y, dx, dy int, knownState BoardState) BoardState {
	// if we already know that someone won, there is no point in searching
	if knownState == O_WON || knownState == X_WON {
		return knownState
	}
	nextState := b.calcLineState(x, y, dx, dy)
	if nextState == TIE {
		// if the line to check is a tie, return previously known state which could be either also a tie,
		// or PLAYING
		return knownState
	}
	return nextState
}

// calculate state of the given line. The line is specified as start point at x and y,
// and is followed by adding dx, and dy to x and y respectively, until the end of the board
// is reached
// If there is a single NONE spot, then the board is in PLAYING state
// if all the cells are occupied by a single player, then board is in WON state, by the
// said player
// Otherwise (the line is a mix of different player cells) it's a tie
func (b *Board) calcLineState(x, y, dx, dy int) BoardState {
	spotted := NO_PLAYER
	for i := y; i >= 0 && i < b.Size; i += dy {
		for j := x; j >= 0 && j < b.Size; j += dx {
			cell := b.Grid[i][j]
			if cell == NO_PLAYER {
				// if there is an empty cell left, we continue playing
				return PLAYING
			}
			if spotted == NO_PLAYER {
				// this is the first cell we encounter, remember its value
				spotted = cell
			} else if cell != spotted {
				// we found a cell that has a different value from that we've seen
				// in this line, so it's a tie
				return TIE
			}
			if dx == 0 {
				break
			}
		}
		if dy == 0 {
			break
		}
	}
	switch spotted {
	case PLAYER_O:
		return O_WON
	case PLAYER_X:
		return X_WON
	default:
		// just to make compiler happy, this should never be reached
		return TIE
	}
}

func (b *Board) String() string {
	res := ""
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			player, _ := b.GetValue(i, j)
			res += player.String() + " "
		}
		res += "\n"
	}
	// todo: move to corresponding String methods on the types
	boardStateStr := ""
	switch b.State {
	case TIE:
		boardStateStr = "TIE"
	case PLAYING:
		boardStateStr = "PLAYING"
	case O_WON:
		boardStateStr = "Player O wins!"
	case X_WON:
		boardStateStr = "Player X wins!"
	}
	res += boardStateStr + "\n"
	res += "Next turn: " + b.NextTurn.String()
	return res
}
