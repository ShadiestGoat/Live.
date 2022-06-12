package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Goal struct {
	_Coords Vector
	_HP int

	_IVFrames int
	IVFramesMax int
}

func (g Goal) Coords() Vector {
	return g._Coords
}

func (g *Goal) Move(v Vector, game *Game) bool {
	g._Coords.Add(v)
	return true
}

func (g Goal) Velocity(game Game) Vector {
	if g._Coords.DistanceToCenter() > RADIUS_DESPAWN {
		return Vector{}
	}
	return GetVelocityToCenter(&g, game.Resources.GoalInfo.Speed)
}

func (g Goal) HP() int {
	return g._HP
}

func (g Goal) HPMax(goal Game) int {
	return goal.Resources.GoalInfo.HP
}

func (g Goal) Sprite(game Game) *ebiten.Image {
	return game.Resources.Goal
}

func (g Goal) Box(game Game) image.Rectangle {
	return BasicBox(g._Coords, game.Resources.Goal)
}

func (g *Goal) setHP(newHP int) {
	g._HP = newHP
}

func (g Goal) IVFrames() int {
	return g._IVFrames
}

func (g *Goal) DecreaseIVFRames(game Game) {
	if g._IVFrames != 0 {
		g._IVFrames--
	}
}

func (g *Goal) ResetIVFrames() {
	g._IVFrames = g.IVFramesMax
}
