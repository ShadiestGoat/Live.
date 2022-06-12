package main

// TODO: restructure alive things to use interface, so that theres a universal way to check for inv frames & a universal way to damage

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Could also be in place of coordinates
// The X values is the first value, the second value is Y
// This is done for better readability.
type Vector [2]int

type Game struct {
	Resources *ResourcePack
	BGOffset Vector
	
	IsPaused bool
	Time int

	Protag Protag
	Enemies []Enemy

	Goal Goal
}

type ItemRarity string

const (
	// Can only get common items
	IR_COMMON ItemRarity = "COMMON"
	// Can get rare (30%) common (70%)
	IR_RARE ItemRarity = "RARE"
	// Can get rare+ (20%), rare (45%) and common (35%)
	IR_RARER ItemRarity = "RARE+"
	// Can get legendary (65%), rare+ (17%), rare (16%) and common (2%)
	IR_LEGEND ItemRarity = "LEGENDARY"
)

type Reward struct {
	XP int `json:"xp"`
	Gold int `json:"gold"`
	ItemChance float64 `json:"itemChance"`
	ItemRarity ItemRarity `json:"itemRarity"`
}


type ActiveAbility interface {
	Cooldown() int
	MaxCooldown() int
	// If it has stat changes, it should be done. 
	Start(g *Game)
	// Returns in ticks. If 0, then it's not active
	// TODO: Remember to implement reversal of applied stat changes!
	ActiveTimeLeft() int
	
	// If it's a weapon, then like it has a collision box. Otherwise, no.
	IsWeapon() bool
	// Returns attack if it's a weapon, otherwise, 0
	Damage() int

	Box() image.Rectangle
	Draw(screen *ebiten.Image, g Game)
	Update()

	Move(dx int, dy int)
}

type Alive interface {
	Coords() Vector

	// returns true if it should still be alive (ie. not despawned)
	Move(v Vector, g *Game) bool

	Velocity(g Game) Vector

	HP() int
	HPMax(g Game) int

	Sprite(g Game) *ebiten.Image
	Box(g Game) image.Rectangle

	// unsafe method to set this creature's hp. For a safe method see tools.go ChangeHP()
	setHP(newHP int)

	IVFrames() int
	DecreaseIVFRames(g Game)
	ResetIVFrames()
}
