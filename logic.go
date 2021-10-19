package main

import (
	"log"
)

func createGameBoardExtended(gameState GameState) GameBoardExtended {
	boardHeight := gameState.Board.Height
	boardWidth := gameState.Board.Width
	var gameBoard GameBoardExtended
	for i := 0; i < boardHeight; i++ {
		var line = []Tile{}
		for j := 0; j < boardWidth; j++ {
			line = append(line, Tile{content: Empty})
		}
		gameBoard = append(gameBoard, line)
	}

	// fill content - snakes
	for _, bs := range gameState.Board.Snakes {
		headEx := gameBoard[bs.Head.X][bs.Head.Y]
		headEx.content = Head
		for _, bsCoord := range bs.Body {
			ce := gameBoard[bsCoord.X][bsCoord.Y]
			ce.content = Body
		}
	}
	// fill content - food
	for _, food := range gameState.Board.Food {
		gameBoard[food.X][food.Y].content = Food
	}

	return gameBoard
}

func getNeighbourScore(gb GameBoardExtended, current Coord) int {
	boardHeight := len(gb)
	boardWidth := len(gb[0])

	score := 0
	var x, y int

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			x = current.X + i
			y = current.Y + j
			if x < 0 || y < 0 || x == boardWidth || y == boardHeight {
				score += int(Wall)
			} else {
				score += gb[x][y].value()
			}
		}
	}

	return score
}

func move(state GameState) BattlesnakeMoveResponse {
	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// Step 1 - Don't hit walls.
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height

	leftEnd := 0
	rightEnd := boardWidth - 1
	bottomEnd := 0
	topEnd := boardHeight - 1

  myHead := state.You.Head

	// The head is on left or right ends
	if myHead.X == leftEnd {
		possibleMoves["left"] = false
	}
	if myHead.X == rightEnd {
		possibleMoves["right"] = false
	}
	// The head is on bottom or top ends
	if myHead.Y == bottomEnd {
		possibleMoves["down"] = false
	}
	if myHead.Y == topEnd {
		possibleMoves["up"] = false
	}

	// Step 2 - Don't hit any snakes (including self).

  gameBoardEx := createGameBoardExtended(state)
  
  leftHead := Coord{X: myHead.X - 1, Y: myHead.Y}
  rightHead := Coord{X: myHead.X + 1, Y: myHead.Y}
  upHead := Coord{X: myHead.X, Y: myHead.Y + 1}
  downHead := Coord{X: myHead.X, Y: myHead.Y - 1}

  if possibleMoves["left"] {
    leftHeadField := gameBoardEx[leftHead.X][leftHead.Y]
    possibleMoves["left"] = leftHeadField.isSafeMove()
  }

  if possibleMoves["right"] {
    rightHeadField := gameBoardEx[rightHead.X][rightHead.Y]
    possibleMoves["right"] = rightHeadField.isSafeMove()
  }

  if possibleMoves["up"] {
    upHeadField := gameBoardEx[upHead.X][upHead.Y]
    possibleMoves["up"] = upHeadField.isSafeMove()
  }

  if possibleMoves["down"] {
    downHeadField := gameBoardEx[downHead.X][downHead.Y]
    possibleMoves["down"] = downHeadField.isSafeMove()
  }


	safeMoves := []string{}
	for move, isSafe := range possibleMoves {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

  var maxScore = -9000
  var nextMove string = "down"

	for _, move := range safeMoves {
		var neighbour Coord
		switch move {
		case "up":
			neighbour = upHead
		case "down":
			neighbour = downHead
		case "left":
			neighbour = leftHead
		case "right":
			neighbour = rightHead
		}

    nScore := getNeighbourScore(gameBoardEx, neighbour)
    log.Printf("For possible move '%s' we've calculated score: %d", move, nScore)  
    if nScore > maxScore {
      maxScore = nScore
      nextMove = move
    }
	}

	if len(safeMoves) == 0 {
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	}
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "virtuallight",
		Color:      "#88ff88",
		Head:       "smile",
		Tail:       "bolt",
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}
