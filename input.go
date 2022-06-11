package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Direction int8

const (
	DirUp Direction = iota
	DirRight
	DirDown
	DirLeft
)

var DirEnum = map[Direction]Vector{
	DirUp:    {0, -1},
	DirRight: {1, 0},
	DirDown:  {0, 1},
	DirLeft:  {-1, 0},
}

// usefull for debugging
var DirStr = map[Direction]string{
	DirUp:    "up",
	DirRight: "right",
	DirDown:  "down",
	DirLeft:  "left",
}

func (g *Game) ProtagMove(dirs Vector) {
	sM := float64(g.Protag.Speed)

	if dirs[0] != 0 && dirs[1] != 0 {
		sM *= 1/math.Sqrt2
	}

	dirs[0] *= int(math.Round(sM))
	dirs[1] *= int(math.Round(sM))

	g.BGOffset.Add(dirs)
	
	g.BGOffset[0] = Overflow(g.BGOffset[0], g.Resources.BGSize)
	g.BGOffset[1] = Overflow(g.BGOffset[1], g.Resources.BGSize)

	newMonsters := []Enemy{}

	for _, m := range g.Enemies {
		m.Coords.Add(dirs)
		newMonsters = append(newMonsters, m)
	}

	for id := range g.Protag.Abilities {
		g.Protag.Abilities[id].Move(dirs[0], dirs[1])
	}

	g.Enemies = newMonsters
}

func keyHeld(key ebiten.Key) bool {
	if !ebiten.IsKeyPressed(key) {
		return false
	}
	d := inpututil.KeyPressDuration(key)
	return d%5 == 0
}

// returns true if at least 1 key of keys is pressed
func oneKeyPressed(keys []ebiten.Key) bool {
	for _, key := range keys {
		if ebiten.IsKeyPressed(key) {
			return true
		}
	}
	return false
}
