package main

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

type Enemy struct {
	Speed int
	HP int
	MaxHP int
	Reward Reward
	Shape EnemyShape
	Coords Vector
	Sprite *ebiten.Image
}

type Protag struct {
	MaxHP int
	HP int
	XP int
	Level int
	Speed int
	Luck int
	// This is what should be used for collision, not width/height! It's a circle radius.
	Reach int

	Coins int
	Abilities map[ActionUpgradeID]ActiveAbility
}

type  ActiveAbility interface {
	Cooldown() int
	MaxCooldown() int
	// If it has stat changes, it should be done. 
	Start(g *Game)
	// Returns in ticks. If 0, then it's not active
	// TODO: Remember to implement reversal of applied stat changes!
	ActiveTimeLeft() int
	
	// If it's a weapon, then like it has a collision box. Otherwise, no.
	IsWeapon() bool
	Box() image.Rectangle
	Draw(screen *ebiten.Image, g Game)
	Update()

	Move(dx int, dy int)
}