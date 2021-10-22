package main

import (
	"testing"

	"github.com/matryer/is"
)

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

	testData := []struct {
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
			expectedExtendedBoard: GameBoardExtended{
				{Tile{}, Tile{}, Tile{}},
				{Tile{}, Tile{}, Tile{}},
			},
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
			expectedExtendedBoard: GameBoardExtended{
				{Tile{}, Tile{}, Tile{}, Tile{Food}},
				{Tile{}, Tile{}, Tile{Head}, Tile{}},
				{Tile{Body}, Tile{Body}, Tile{Head}, Tile{Food}},
			},
		},
	}

	is := is.New(t)

	for _, data := range testData {
		gameBoardExtended := createGameBoardExtended(data.state)
		is.Equal(gameBoardExtended, data.expectedExtendedBoard) // The result gameBoardExtended isn't what we expected it to be
	}
}
