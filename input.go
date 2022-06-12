package main

import (
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

func (g *Game) ProtagMove(velocity Vector) {
	g.BGOffset.Add(velocity)
	
	g.BGOffset[0] = Overflow(g.BGOffset[0], g.Resources.BGSize)
	g.BGOffset[1] = Overflow(g.BGOffset[1], g.Resources.BGSize)

	newMonsters := []Enemy{}

	for _, m := range g.Enemies {
		m.Move(velocity, g)
		newMonsters = append(newMonsters, m)
	}

	for id := range g.Protag.Abilities {
		g.Protag.Abilities[id].Move(velocity[0], velocity[1])
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
