package main

import (
	"fmt"
)

// Cell state of the board
type Cell int

// NONE cell is a cell that no player put mark on yet
// OCCUPIED_X and OCCUPIED_Y values are for cells with X and O respectively
const (
	NONE Cell = iota
	OCCUPIED_X
	OCCUPIED_O
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
	PLAYER_X Player = iota
	PLAYER_O
)

// Board represents tic tac toe game state
// board is a square grid of cells
type Board struct {
	size     int
	grid     [][]Cell
	nextTurn Player
}

// MakeBoard creates a board of size rows and size columns
func MakeBoard(size int) *Board {
	grid := make([][]Cell, size)
	for i := 0; i < size; i++ {
		row := make([]Cell, size)
		for j := 0; j < size; j++ {
			row[j] = NONE
		}
		grid[i] = row
	}
	return &Board{size, grid, PLAYER_X}
}

// SetValue for the board at x, y to the given value
// return error if the cell already has a value distinct from NONE
// or the coordinates are outside of the grid
func (b *Board) SetValue(value Cell, x, y int) error {
	err := b.validateCoordinates(x, y)
	if err != nil {
		return err
	}
	if b.grid[y][x] != NONE {
		return fmt.Errorf("%d, %d is already occupied", x, y)
	}
	b.grid[y][x] = value
	return nil
}

// GetValue returns value of the board at x, y
func (b *Board) GetValue(x, y int) (Cell, error) {
	err := b.validateCoordinates(x, y)
	if err != nil {
		return 0, err
	}
	return b.grid[y][x], nil
}

// return nil if x, y coordinates are valid for this board
// otherwise return error
func (b *Board) validateCoordinates(x, y int) error {
	var err error
	if x < 0 || x >= b.size {
		err = fmt.Errorf("X coordinate must be within [0, %d]", b.size-1)
	}
	if y < 0 || y >= b.size {
		err = fmt.Errorf("Y coordinate must be within [0, %d], %s", b.size-1, err)
	}
	return err
}

// GetBoardState returns state in which current board is, together with the
// player that should make next turn
func (b *Board) GetBoardState() (BoardState, Player) {
	state := TIE
	for row := 0; row < b.size; row++ {
		state = b.getLineState(0, row, 1, 0, state)
	}
	for col := 0; col < b.size; col++ {
		state = b.getLineState(col, 0, 0, 1, state)
	}
	// check main diagonal
	state = b.getLineState(0, 0, 1, 1, state)
	// check secondary diagonal
	state = b.getLineState(b.size-1, 0, -1, 1, state)
	return state, b.nextTurn
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
	spotted := NONE
	for i := y; i >= 0 && i < b.size; i += dy {
		for j := x; j >= 0 && j < b.size; j += dx {
			cell := b.grid[i][j]
			if cell == NONE {
				// if there is an empty cell left, we continue playing
				return PLAYING
			}
			if spotted == NONE {
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
	case OCCUPIED_O:
		return O_WON
	case OCCUPIED_X:
		return X_WON
	default:
		// just to make compiler happy, this should never be reached
		return TIE
	}
}

func (b *Board) String() string {
	res := ""
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			val, _ := b.GetValue(i, j)
			switch val {
			case NONE:
				res += "-"
			case OCCUPIED_X:
				res += "X"
			case OCCUPIED_O:
				res += "O"
			}
			res += " "
		}
		res += "\n"
	}
	// todo: move to corresponding String methods on the types
	boardStateStr := ""
	state, nextTurn := b.GetBoardState()
	switch state {
	case TIE:
		boardStateStr = "TIE"
	case PLAYING:
		boardStateStr = "PLAYING"
	case O_WON:
		boardStateStr = "Player O wins!"
	case X_WON:
		boardStateStr = "Player X wins!"
	}
	nextTurnStr := ""
	res += "Board state: " + boardStateStr + "\n"
	switch nextTurn {
	case PLAYER_O:
		nextTurnStr = "Player O"
	case PLAYER_X:
		nextTurnStr = "Player X"
	}
	res += "Next turn: " + nextTurnStr
	return res
}
