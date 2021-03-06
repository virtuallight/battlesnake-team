package main

type ContentType int

const FOOD_BONUS = 10

const (
	Empty  ContentType = 0
	Food   ContentType = 2 // TODO: figure out a better value for this in relation to health
	Body   ContentType = -1
	Head   ContentType = -2
	Wall   ContentType = -3
	Hazard ContentType = -2
)

type Tile struct {
	content ContentType
}

func (t *Tile) value() int {
	return int(t.content)
}

func (t *Tile) isSafeMove() bool {
	return t.content == Empty || t.content == Food
}

type GameBoardExtended [][]Tile

// --

type GameState struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

type Game struct {
	ID      string  `json:"id"`
	Ruleset Ruleset `json:"ruleset"`
	Timeout int32   `json:"timeout"`
}

type Ruleset struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`

	// Used in non-standard game modes
	Hazards []Coord `json:"hazards"`
}

type Battlesnake struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Health  int32   `json:"health"`
	Body    []Coord `json:"body"`
	Head    Coord   `json:"head"`
	Length  int32   `json:"length"`
	Latency string  `json:"latency"`

	// Used in non-standard game modes
	Shout string `json:"shout"`
	Squad string `json:"squad"`
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}
