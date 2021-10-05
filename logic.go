package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

// Points scale = -8 <--> +8
// Food +2
// Wall -1
// Another|Our Snake Head -1
// Another|Our Snake Body -2
// Another|Our Snake Tail -1
// Empty 0

import (
	"log"
	"math/rand"
)



func createGameBoardExtended(gameState GameState) GameBoardExtended {
	boardHeight := gameState.Board.Height
	boardWidth := gameState.Board.Width
	var gameBoard GameBoardExtended
	for i := 0; i < boardHeight; i++ {
		var line = []CoordExtended{}
		for j := 0; j < boardWidth; j++ {
			line = append(line, CoordExtended{content: Empty})
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

  // calculate score
	for _, bs := range gameState.Board.Snakes {
    for _, bsCoord := range bs.Body {
      ce := gameBoard[bsCoord.X][bsCoord.Y]
      ce.totalScore = getNeighbourScore(gameBoard, bsCoord)
    }
  }

	return gameBoard
}

func getNeighbourScore(gb GameBoardExtended, current Coord) int {
  score := 0
  var x, y int
  
  me := gb[current.X][current.Y]
  score += me.score()

  // TODO calculate score based on neighbours (mind the wall!)
  // left neighbour
  x = current.X - 1
  y = current.Y

  // right neighbour
  x = current.X + 1
  y = current.Y

  // upper neighbour
  x = current.X
  y = current.Y + 1

  // lower neighbour
  x = current.X
  y = current.Y - 1

  // upper left neighbour
  x = current.X - 1
  y = current.Y + 1 

  // upper right neighbour
  x = current.X + 1
  y = current.Y + 1 

  // lower left neighbour
  x = current.X - 1
  y = current.Y - 1 

  // lower right neighbour
  x = current.X + 1
  y = current.Y - 1 

  return score
}

func move(state GameState) BattlesnakeMoveResponse {
	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// Step 0: Don't let your Battlesnake move back in on it's own neck
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of body piece directly behind your head (your "neck")
	if myNeck.X < myHead.X {
		possibleMoves["left"] = false
	} else if myNeck.X > myHead.X {
		possibleMoves["right"] = false
	} else if myNeck.Y < myHead.Y {
		possibleMoves["down"] = false
	} else if myNeck.Y > myHead.Y {
		possibleMoves["up"] = false
	}

	// Step 1 - Don't hit walls.
	// Use information in GameState to prevent your Battlesnake from moving beyond the boundaries of the board.
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height

	leftEnd := 0
	rightEnd := boardWidth - 1
	bottomEnd := 0
	topEnd := boardHeight - 1

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

	// Step 2 - Don't hit yourself.
	// Use information in GameState to prevent your Battlesnake from colliding with itself.
	myBody := state.You.Body

	// Check the 4 positions arround the head
	leftHead := Coord{X: myHead.X - 1, Y: myHead.Y}
	rightHead := Coord{X: myHead.X + 1, Y: myHead.Y}
	bottomHead := Coord{X: myHead.X, Y: myHead.Y - 1}
	topHead := Coord{X: myHead.X, Y: myHead.Y + 1}

	for _, bodyPart := range myBody {
		if leftHead == bodyPart {
			possibleMoves["left"] = false
		}
		if rightHead == bodyPart {
			possibleMoves["right"] = false
		}
		if bottomHead == bodyPart {
			possibleMoves["down"] = false
		}
		if topHead == bodyPart {
			possibleMoves["up"] = false
		}
	}

	// TODO: Step 3 - Don't collide with others.
	// Use information in GameState to prevent your Battlesnake from colliding with others.

	// TODO: Step 4 - Find food.
	// Use information in GameState to seek out and find food.

	// Finally, choose a move from the available safe moves.
	// TODO: Step 5 - Select a move to make based on strategy, rather than random.
	var nextMove string

	safeMoves := []string{}
	for move, isSafe := range possibleMoves {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		nextMove = "down"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		nextMove = safeMoves[rand.Intn(len(safeMoves))]
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
