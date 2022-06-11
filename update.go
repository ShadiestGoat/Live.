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

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		pauseCooldown--
		if pauseCooldown <= 0 {
			g.IsPaused = !g.IsPaused
			pauseCooldown = 10
		}
	}
	// FIXME: Debug only!
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.Restart()
		return nil
	}
	if g.IsPaused {
		return nil
	}
	g.Time++
	dirs := map[Direction]bool{}
	dirs[DirUp] = ebiten.IsKeyPressed(ebiten.KeyArrowUp)
	dirs[DirRight] = ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	dirs[DirDown] = ebiten.IsKeyPressed(ebiten.KeyArrowDown)
	dirs[DirLeft] = ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	
	offset := ResolveVector(dirs)
	curMovementVector = offset
	
	if offset[0] != 0 {
		lastMovementVector[0] = offset[0]
	}
	
	if offset[1] != 0 {
		lastMovementVector[1] = offset[1]
	}
	
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		if g.Protag.Abilities[AU_SLASH].Cooldown() <= 0 {
			g.Protag.Abilities[AU_SLASH].Start(g)
		}
	}

	g.ProtagMove(offset)

	// FIXME: Debug only!
	addSpeed := ResolveVector(map[Direction]bool{
		DirUp:   keyHeld(ebiten.KeyA),
		DirDown: keyHeld(ebiten.KeyS),
	})

	g.Protag.Speed += addSpeed[1]

	summonCooldown--
	
	if summonCooldown == 0 {
		summonCooldown = g.Resources.SpawnRate.Resolve(g.Time)
		
		// put a limit on the amount that can be summoned
		if len(g.Enemies) < 550 {
			g.Summon()
		}
	}

	newMonsters := []Enemy{}

	weaponBoxes := []image.Rectangle{}
	
	for i := range g.Protag.Abilities {
		g.Protag.Abilities[i].Update()
		if g.Protag.Abilities[i].ActiveTimeLeft() == 0 {
			continue
		}
		weaponBoxes = append(weaponBoxes, g.Protag.Abilities[i].Box())
	}

	for _, m := range g.Enemies {
		// despawn far away monsters
		if m.Coords.DistanceToCenter() > RADIUS_DESPAWN {
			continue
		}

		bounds := m.Rect()

		killed := false
		for _, weapon := range weaponBoxes {
			if RectCollision(weapon, bounds) {
				killed = true
				break
			}
		}
		if killed {
			continue
		}
		if RectCollision(bounds, ProtagBox) {
			if g.Time > bestTime {
				bestTime = g.Time
			}
			g.Restart()
			return nil
		}

		velocity := m.Move()
		m.Coords.Add(velocity)
		

		newMonsters = append(newMonsters, m)
	}

	g.Enemies = newMonsters

	if slowTps {
		SlowTPS()
	}

	return nil
}
