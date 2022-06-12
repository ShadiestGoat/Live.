package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var curMovementVector = Vector{}
var lastMovementVector = Vector{}

var summonCooldown = 60
var pauseCooldown = 10

var slowTps = false

var Debug string

func DebugUpdate(g *Game) bool {
	addSpeed := ResolveVector(map[Direction]bool{
		DirDown: keyHeld(ebiten.KeyO),
		DirUp:   keyHeld(ebiten.KeyP),
	})

	g.Protag.Speed += addSpeed[1]
	return false
}

func (g *Game) Update() error {
	if Debug == "t" {
		if DebugUpdate(g) {
			return nil
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.Restart()
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		pauseCooldown--
		if pauseCooldown <= 0 {
			g.IsPaused = !g.IsPaused
			pauseCooldown = 10
		}
	}
	if g.IsPaused {
		return nil
	}
	g.Time++

	protagVelocity := g.Protag.Velocity(*g)
	offset := g.Protag.GetOffset()

	curMovementVector = offset
	
	if offset[0] != 0 {
		lastMovementVector[0] = offset[0]
	}
	
	if offset[1] != 0 {
		lastMovementVector[1] = offset[1]
	}
	
	if oneKeyPressed([]ebiten.Key{ebiten.KeyZ, ebiten.KeyE}) {
		if g.Protag.Abilities[AU_SLASH].Cooldown() <= 0 {
			g.Protag.Abilities[AU_SLASH].Start(g)
		}
	}
	
	g.Protag.DecreaseIVFRames(*g)

	g.Protag.RegenFreq--
	if g.Protag.RegenFreq == 0 {
		g.ChangeHP(1, &g.Protag)
		g.Protag.RegenFreq = g.Protag.MaxRegenFreq
	}

	g.ProtagMove(protagVelocity)

	summonCooldown--
	
	if summonCooldown == 0 {
		summonCooldown = g.Resources.SpawnRate.Resolve(g.Time)
		
		// put a limit on the amount that can be summoned
		if len(g.Enemies) < 600 {
			g.Summon()
		}
	}

	newMonsters := []Enemy{}

	weaponBoxes := []image.Rectangle{}
	dmgWeapons := []int{}
	
	for i := range g.Protag.Abilities {
		g.Protag.Abilities[i].Update()
		if g.Protag.Abilities[i].ActiveTimeLeft() == 0 {
			continue
		}
		if !g.Protag.Abilities[i].IsWeapon() {
			continue
		}
		weaponBoxes = append(weaponBoxes, g.Protag.Abilities[i].Box())
		dmgWeapons = append(dmgWeapons, g.Protag.Abilities[i].Damage())
	}

	for _, m := range g.Enemies {
		m.DecreaseIVFRames(*g)
		bounds := m.Box(*g)

		killed := false

		if m.IVFrames() == 0 {
			for i, weapon := range weaponBoxes {
				if RectCollision(weapon, bounds) {
					if !g.Hurt(dmgWeapons[i], &m) {
						killed = true
					}
					break
				}
			}
		}

		if killed {
			continue
		}
		if g.Protag.IVFrames() == 0 {
			if RectCollision(bounds, ProtagBox) {
				if !g.Hurt(m.Attack, &g.Protag) {
					g.Restart()
					return nil
				}
			}
		}

		velocity := m.Velocity(*g)
		if !m.Move(velocity, g) {
			continue
		}

		newMonsters = append(newMonsters, m)
	}

	g.Enemies = newMonsters

	if slowTps {
		SlowTPS()
	}

	return nil
}
