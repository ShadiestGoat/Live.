package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Protag struct {
	_MaxHP int
	_HP int

	RegenFreq int
	MaxRegenFreq int

	_IVFrames int
	MaxIVFrames int

	XP int
	Level int
	Speed int
	Luck int
	// This is what should be used for collision, not width/height! It's a circle radius.
	Reach int

	Coins int
	Abilities map[ActionUpgradeID]ActiveAbility
}

func (p Protag) Coords() Vector {
	return CenterCoords
}

func (p Protag) Velocity(g Game) Vector {
	dirs := map[Direction]bool{}
	dirs[DirUp] = oneKeyPressed([]ebiten.Key{ebiten.KeyArrowUp, ebiten.KeyW})
	dirs[DirRight] = oneKeyPressed([]ebiten.Key{ebiten.KeyArrowRight, ebiten.KeyD})
	dirs[DirDown] = oneKeyPressed([]ebiten.Key{ebiten.KeyArrowDown, ebiten.KeyS})
	dirs[DirLeft] = oneKeyPressed([]ebiten.Key{ebiten.KeyArrowLeft, ebiten.KeyA})

	resolved := ResolveVector(dirs)

	sM := float64(p.Speed)

	if resolved[0] != 0 && resolved[1] != 0 {
		// 1/sqrt(2)
		sM *= 0.70710678
	}

	resolved[0] *= int(math.Round(sM))
	resolved[1] *= int(math.Round(sM))

	return resolved
}

func (p *Protag) Move(v Vector, g *Game) bool {
	g.ProtagMove(v)
	return true
}

func (p Protag) HP() int {
	return p._HP
}

func (p Protag) HPMax(g Game) int {
	return p._MaxHP
}

func (p *Protag) setHP(newHP int) {
	p._HP = newHP
}

func (p Protag) Sprite(g Game) *ebiten.Image {
	return g.Resources.Protag
}

func (p Protag) Box(g Game) image.Rectangle {
	return ProtagBox
}

func (p Protag) IVFrames() int {
	return p._IVFrames
}

func (p *Protag) DecreaseIVFRames(g Game) {
	if p._IVFrames != 0 {
		p._IVFrames--
	}
}

func (p *Protag) ResetIVFrames() {
	p._IVFrames = p.MaxIVFrames
}
