package main

import (
	"fmt"
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
						Body: []Coord{{X: 2, Y: 1}, {X: 2, Y: 2}, {X: 1, Y: 2}},
					},
				},
				Food: []Coord{
					{X: 3, Y: 2},
					{X: 3, Y: 0},
				},
			},
		},
		expectedExtendedBoard: convertToGBE(VisualGameBoardExtended{
			{Tile{}, Tile{Body}, Tile{Body}, Tile{Food}},
			{Tile{}, Tile{}, Tile{Head}, Tile{}},
			{Tile{Body}, Tile{Body}, Tile{Head}, Tile{Food}},
		}),
	},
}

var neckTestData = []struct {
	state GameState
	neck  string
}{
	// Neck on the left
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
				},
			},
			You: Battlesnake{
				Head: Coord{X: 2, Y: 0},
				Body: []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
			},
		},
		neck: "left",
	},
	// Neck on the right
	{
		state: GameState{
			Board: Board{
				Height: 3,
				Width:  4,
				Snakes: []Battlesnake{
					{
						Head: Coord{X: 0, Y: 0},
						Body: []Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
					},
				},
			},
			You: Battlesnake{
				Head: Coord{X: 0, Y: 0},
				Body: []Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			},
		},
		neck: "right",
	},
	// Neck up
	{
		state: GameState{
			Board: Board{
				Height: 3,
				Width:  4,
				Snakes: []Battlesnake{
					{
						Head: Coord{X: 0, Y: 0},
						Body: []Coord{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}},
					},
				},
			},
			You: Battlesnake{
				Head: Coord{X: 0, Y: 0},
				Body: []Coord{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}},
			},
		},
		neck: "up",
	},
	// Neck down
	{
		state: GameState{
			Board: Board{
				Height: 3,
				Width:  4,
				Snakes: []Battlesnake{
					{
						Head: Coord{X: 1, Y: 2},
						Body: []Coord{{X: 1, Y: 2}, {X: 1, Y: 1}, {X: 2, Y: 1}},
					},
				},
			},
			You: Battlesnake{
				Head: Coord{X: 1, Y: 2},
				Body: []Coord{{X: 1, Y: 2}, {X: 1, Y: 1}, {X: 2, Y: 1}},
			},
		},
		neck: "down",
	},
}

func TestNeckAvoidance(t *testing.T) {
	for _, data := range neckTestData {
		fmt.Println("testing", data.neck)
		nextMove := move(data.state)
		// Assert never move to the nexk
		if nextMove.Move == data.neck {
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

var testCheckFoodData = []struct {
	id                string
	extendedBoard     GameBoardExtended
	head              Coord
	expectedLength    int
	expectedDirection string
}{
	{
		id: "Check for reachable food",
		extendedBoard: convertToGBE(VisualGameBoardExtended{
			{Tile{Body}, Tile{Body}, Tile{}, Tile{}},
			{Tile{Body}, Tile{Head}, Tile{Body}, Tile{Food}},
			{Tile{}, Tile{}, Tile{Body}, Tile{}},
			{Tile{Food}, Tile{}, Tile{Head}, Tile{Food}},
		}),
		head:              Coord{X: 1, Y: 2},
		expectedLength:    3,
		expectedDirection: "down",
	},
	{
		id: "Check for distant reachable food",
		extendedBoard: convertToGBE(VisualGameBoardExtended{
			{Tile{Body}, Tile{Body}, Tile{}, Tile{}},
			{Tile{}, Tile{Head}, Tile{Body}, Tile{Food}},
			{Tile{}, Tile{Body}, Tile{Body}, Tile{}},
			{Tile{}, Tile{Body}, Tile{}, Tile{}},
			{Tile{}, Tile{Body}, Tile{Body}, Tile{}},
			{Tile{}, Tile{}, Tile{Body}, Tile{}},
			{Tile{Food}, Tile{}, Tile{Head}, Tile{Food}},
		}),
		head:              Coord{X: 1, Y: 5},
		expectedLength:    6,
		expectedDirection: "left",
	},
	{
		id: "Head is food",
		extendedBoard: convertToGBE(VisualGameBoardExtended{
			{Tile{Body}, Tile{Body}, Tile{}, Tile{}},
			{Tile{Body}, Tile{Food}, Tile{Body}, Tile{Food}},
			{Tile{}, Tile{}, Tile{Body}, Tile{}},
			{Tile{Food}, Tile{}, Tile{Head}, Tile{Food}},
		}),
		head:              Coord{X: 1, Y: 2},
		expectedLength:    0,
		expectedDirection: "",
	},
}

func TestCheckFood(t *testing.T) {
	for _, data := range testCheckFoodData {
		t.Run(data.id, func(t *testing.T) {
			is := is.NewRelaxed(t)
			shortestLength, direction := checkFood(data.extendedBoard, data.head)
			is.Equal(shortestLength, data.expectedLength) // The shortest path has a different length
			is.Equal(direction, data.expectedDirection)   // The direction for the shortest path is different
		})
	}
}
