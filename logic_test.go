package main

import (
	"testing"
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
  // Creating the GameState
  state := GameState{
    Board: Board{
      Height: 10,
      Width: 20,
    },
	}
  gameBoardExtended := createGameBoardExtended(state)
  if len(gameBoardExtended) != 10 {
    t.Errorf("game board of incorrect height: %s", len(gameBoardExtended))
  }
  for i := 0; i < 10; i++ {
    if len(gameBoardExtended[i]) != 20 {
      t.Errorf("game board of incorrect widgth: %s", len(gameBoardExtended[i]))
    }
  }
}
