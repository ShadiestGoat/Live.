package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	RADIUS_SUMMON float64 = 800
	RADIUS_DESPAWN = 1500
	RADIUS_ITEM_DESPAWN = 2200
)

type Enemy struct {
	Speed int
	_HP int
	_MaxHP int
	Attack int
	Reward Reward
	Shape EnemyShape
	_Coords Vector
	_Sprite *ebiten.Image

	_IVFrames int
	MaxIVFrames int
}

// Creates an enemy based of the info. Coords are empty.
func (e EnemySummonResourcePack) Summon(ticks int) Enemy {
	hp := e.Health.Resolve(ticks)
	return Enemy{
		Speed:      e.Speed.Resolve(ticks),
		_HP:         hp,
		_MaxHP:      hp,
		Reward: Reward{
			XP:   e.RewardXP.Resolve(ticks),
			Gold: e.RewardGold.Resolve(ticks),
			ItemChance: e.RewardItemChance.ResolveFloat(ticks),
			ItemRarity: e.RewardItemRarity,
		},
		Attack: e.Damage.Resolve(ticks),

		Shape:      e.ShapeInfo,
		_Sprite: 	e.Sprite,

		MaxIVFrames: 45,
	}
}

func (g *Game) SummonInfo() Enemy {
	eType := g.Resources.Enemies[RandomInt(0, len(g.Resources.Enemies)-1)]
	enemy := eType.Summon(g.Time)
	angle := RandomFloat(0, 2*math.Pi)
	enemy._Coords = CenterCoords
	enemy._Coords.Add(AngleToCoords(angle, RADIUS_SUMMON))
	return enemy
}

func (g *Game) SummonOne() {
	enemy := g.SummonInfo()
	g.Enemies = append(g.Enemies, enemy)
}

func (g *Game) Summon() {
	maxSpawnAmt := g.Resources.SpawnAmount.Resolve(g.Time)
	amt := RandomInt(int(math.Ceil(0.25*float64(maxSpawnAmt))), maxSpawnAmt)
	for i := 0; i<amt; i++ {
		g.SummonOne()
	}
}

func (e Enemy) Coords() Vector {
	return e._Coords
}

func (e *Enemy) Move(v Vector, g *Game) bool {
	e._Coords.Add(v)
	return e._Coords.DistanceToCenter() < RADIUS_DESPAWN
}

func (m *Enemy) Velocity(game Game) Vector {
	if m._Coords.DistanceToCenter() <= 5 {
		return Vector{}
	}
	return GetVelocityToCenter(m, m.Speed)
}

func (e Enemy) HP() int {
	return e._HP
}

func (e Enemy) HPMax(g Game) int {
	return e._MaxHP
}

func (e Enemy) Sprite(g Game) *ebiten.Image {
	return e._Sprite
}

func (e Enemy) Box(g Game) image.Rectangle {
	return BasicBox(e._Coords, e._Sprite)
}

func (e *Enemy) setHP(newHP int) {
	e._HP = newHP
}

func (e Enemy) IVFrames() int {
	return e._IVFrames
}

func (e *Enemy) DecreaseIVFRames(g Game) {
	if e._IVFrames != 0 {
		e._IVFrames--
	}
}

func (e *Enemy) ResetIVFrames() {
	e._IVFrames = e.MaxIVFrames
}
