package main

import (
	"testing"

	"github.com/matryer/is"
)

var testData = []struct {
	state                 GameState
	expectedExtendedBoard GameBoardExtended
}{
	{
		state: GameState{
			Board: Board{
				Height: 2,
				Width:  3,
			},
		},
		expectedExtendedBoard: convertToGBE(VisualGameBoardExtended{
			{Tile{}, Tile{}, Tile{}},
			{Tile{}, Tile{}, Tile{}},
		}),
	},
	{
		state: GameState{
			Board: Board{
				Height: 3,
				Width:  4,
				Snakes: []Battlesnake{
					{
						Head: Coord{X: 2, Y: 0},
						Body: []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
					},
					{
						Head: Coord{X: 2, Y: 1},
						Body: []Coord{{X: 2, Y: 1}},
					},
				},
				Food: []Coord{
					{X: 3, Y: 2},
					{X: 3, Y: 0},
				},
			},
		},
		expectedExtendedBoard: convertToGBE(VisualGameBoardExtended{
			{Tile{}, Tile{}, Tile{}, Tile{Food}},
			{Tile{}, Tile{}, Tile{Head}, Tile{}},
			{Tile{Body}, Tile{Body}, Tile{Head}, Tile{Food}},
		}),
	},
}

func TestNeckAvoidance(t *testing.T) {
	// Arrange
	me := Battlesnake{
		// Length 3, facing right
		Head: Coord{X: 2, Y: 0},
		Body: []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Height: 10,
			Width:  20,
			Snakes: []Battlesnake{me},
		},
		You: me,
	}

	// Act 1,000x (this isn't a great way to test, but it's okay for starting out)
	for i := 0; i < 1000; i++ {
		nextMove := move(state)
		// Assert never move left
		if nextMove.Move == "left" {
			t.Errorf("snake moved onto its own neck, %s", nextMove.Move)
		}
	}
}

func TestCreateGameBoardExtended(t *testing.T) {

	is := is.New(t)

	for _, data := range testData {
		gameBoardExtended := createGameBoardExtended(data.state)
		is.Equal(gameBoardExtended, data.expectedExtendedBoard) // The result gameBoardExtended isn't what we expected it to be
	}
}

func TestVisualConvert(t *testing.T) {
	is := is.New(t)

	for _, data := range testData {
		visual := convertFromGBE(data.expectedExtendedBoard)
		backGBE := convertToGBE(visual)
		is.Equal(backGBE, data.expectedExtendedBoard) // Converting extended game board to visual representation and back should give back the same board.
	}
}

type VisualGameBoardExtended [][]Tile

// convertToGBE takes a visual representation of game board
// (where value at (0, 0) is at the bottom left corner)
// and converts it, so that the field at (0, 0)
// is actually represented on GameBoardExtended at the [0][0] index and so on.
// This function is useful for taking a visual representation of the board
// and converting it to something that logic understands.
func convertToGBE(visual VisualGameBoardExtended) GameBoardExtended {
	gbe := GameBoardExtended{}
	rows := len(visual[0]) // height
	columns := len(visual) // width
	for x := 0; x < rows; x++ {
		line := []Tile{}
		for y := 0; y < columns; y++ {
			line = append(line, visual[columns-y-1][x])
		}
		gbe = append(gbe, line)
	}
	return gbe
}

// convertFromGBE takes an extended game board argument
// (where value at (0, 0) is at index [0][0])
// and converts it, so that the field at (0, 0)
// is at the bottom left corner.
// This function is useful for printing out GameBoardExtended values
// for debugging.
func convertFromGBE(gbe GameBoardExtended) VisualGameBoardExtended {
	visual := VisualGameBoardExtended{}

	rows := len(gbe)       // height
	columns := len(gbe[0]) // width
	for x := 0; x < columns; x++ {
		line := []Tile{}
		for y := 0; y < rows; y++ {
			line = append(line, gbe[y][columns-x-1])
		}
		visual = append(visual, line)
	}
	return visual
}
