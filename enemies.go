package main

import (
	"fmt"
	"image"
	"math"
)

var (
	RADIUS_SUMMON float64 = 800
	RADIUS_DESPAWN = 1500
	RADIUS_ITEM_DESPAWN = 2200
)

func (e Enemy) Rect() image.Rectangle {
	return image.Rect(e.Coords[0]-e.Sprite.Bounds().Dx()/2, e.Coords[1]-e.Sprite.Bounds().Dy()/2, e.Coords[0]+e.Sprite.Bounds().Dx()/2, e.Coords[1]+e.Sprite.Bounds().Dy()/2)
}

// Creates an enemy based of the info. Coords are empty.
func (e EnemySummonResourcePack) Summon(ticks int) Enemy {
	hp := e.Health.Resolve(ticks)
	return Enemy{
		Speed:      e.Speed.Resolve(ticks),
		HP:         hp,
		MaxHP:      hp,
		Reward: Reward{
			XP:   e.RewardXP.Resolve(ticks),
			Gold: e.RewardGold.Resolve(ticks),
			ItemChance: e.RewardItemChance.ResolveFloat(ticks),
			ItemRarity: e.RewardItemRarity,
		},

		Shape:      e.ShapeInfo,
		Sprite: 	e.Sprite,
	}
}

func (g *Game) SummonInfo() Enemy {
	eType := g.Resources.Enemies[RandomInt(0, len(g.Resources.Enemies)-1)]
	enemy := eType.Summon(g.Time)
	angle := RandomFloat(0, 2*math.Pi)
	enemy.Coords = CenterCoords
	enemy.Coords.Add(AngleToCoords(angle, RADIUS_SUMMON))
	return enemy
}

func (g *Game) SummonOne() {
	enemy := g.SummonInfo()
	g.Enemies = append(g.Enemies, enemy)
}

func (g *Game) Summon() {
	amt := RandomInt(1, g.Resources.SpawnAmount.Resolve(g.Time))
	fmt.Println("Summon: ", amt)
	for i := 0; i<amt; i++ {
		g.SummonOne()
	}
}

func (m *Enemy) Move() Vector {
	if m.Coords.DistanceToCenter() <= 5 {
		return Vector{}
	}
	
	diff := Vector{
		CenterCoords[0]-m.Coords[0],
		CenterCoords[1]-m.Coords[1],
	}

	ang := math.Atan(float64(diff[0])/float64(diff[1]))

	dx := math.Cos(ang)
	dy := math.Sin(ang)
	s := true

	if diff[1] == 0 {
		s = false
		dx = 0
		dy = 1
		if diff[0] < 0 {
			dy = -1
		}
	}
	if diff[0] == 0 {
		s = false
		dy = 0
		dx = 1
		if diff[1] < 0 {
			dx = -1
		}
	}

	
	if diff[1] < 0 && s {
		dx *= -1
		dy *= -1
	}

	dy *= float64(m.Speed)
	dx *= float64(m.Speed)
	
	return Vector{
		int(math.Round(dy)),
		int(math.Round(dx)),
	}
}
