package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type ActionUpgradeID int

const (
	AU_NONE ActionUpgradeID = iota
	AU_SLASH
)

type IncreasingScaleAbility struct {
	_Cooldown int
	_MaxCooldown int

	Attack int

	XP int
	Level int

	CurExecutionTick float64
	// When to begin the 'death' animation of this sprite
	DspExecutionTick float64
	MaxExecutionTick float64

	ActiveDir float64
	X int
	Y int
	Y2 int

	Sprite *ebiten.Image
}

func (a *IncreasingScaleAbility) Move(dx int, dy int) {
	a.X += dx
	a.Y += dy
	a.Y2 += dy
}

func (a IncreasingScaleAbility) Cooldown() int {
	return a._Cooldown
}

func (a IncreasingScaleAbility) MaxCooldown() int {
	return a._MaxCooldown
}

func (a *IncreasingScaleAbility) Start(g *Game) {
	a.ActiveDir = -float64(lastMovementVector[0])
	if a.ActiveDir == 0 {
		a.ActiveDir = -1
	}
	a.X = CenterCoords[0] + (int(a.ActiveDir*(math.Round(float64(g.Resources.Protag.Bounds().Dx())/2) + 5)))
	
	y := CenterCoords[1]
	y2 := y - int(math.Round(float64(g.Resources.Protag.Bounds().Dy())/2))
	y += int(math.Round(float64(g.Resources.Protag.Bounds().Dy())/2))
	
	a.Y = y
	a.Y2 = y2

	a.CurExecutionTick++

	a._Cooldown = a._MaxCooldown
}

func (a IncreasingScaleAbility) ActiveTimeLeft() int {
	if a.CurExecutionTick == 0 {
		return 0
	}
	return int(a.MaxExecutionTick - a.CurExecutionTick)
}

func (a IncreasingScaleAbility) IsWeapon() bool {
	return true
}

func (a IncreasingScaleAbility) Box() image.Rectangle {
	x := a.X
	if a.CurExecutionTick >= a.DspExecutionTick {
		x += int(a.ActiveDir*math.Round((1-math.Abs(a.ScaleX()))*float64(a.Sprite.Bounds().Dx())))
	}
	x1 := x
	x += int(math.Round(float64(a.Sprite.Bounds().Dx())*a.ScaleX()))

	return image.Rect(
		x1,
		a.Y,
		x,
		a.Y2,
	)
}

var aaaaaaa int

func (a IncreasingScaleAbility) ScaleX() float64 {
	scaleX := a.CurExecutionTick/a.DspExecutionTick
	if a.CurExecutionTick >= a.DspExecutionTick {
		scaleX = 1-(a.CurExecutionTick-a.DspExecutionTick)/(a.MaxExecutionTick-a.DspExecutionTick)
	}

	if scaleX > 1 {
		scaleX = 1
	}
	scaleX *= float64(a.ActiveDir)
	return math.Round(scaleX*100)/100
}

func (a *IncreasingScaleAbility) Update() {
	if a.CurExecutionTick == 0 {
		if a._Cooldown != 0 {
			a._Cooldown--
		}
		return
	}
	if a.CurExecutionTick == a.MaxExecutionTick {
		a.CurExecutionTick = 0
		return
	}
	a.CurExecutionTick++
}

func (a *IncreasingScaleAbility) Draw(screen *ebiten.Image, g Game) {
	if a.CurExecutionTick == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	scaleX := a.ScaleX()
	t2 := 0.0
	if a.CurExecutionTick >= a.DspExecutionTick {
		t2 = (1-math.Abs(scaleX))
		t2*=float64(a.Sprite.Bounds().Dx())
		t2 = math.Round(t2)*a.ActiveDir
	}
	op.GeoM.Scale(scaleX, 1)
	op.GeoM.Translate(float64(a.X)+t2, float64(a.Y2))
	screen.DrawImage(a.Sprite, op)
}

var slashTexture *ebiten.Image

func init() {
	slashTexture = loadTexture(loadFile("slash.png"))
}

func (a IncreasingScaleAbility) Damage() int {
	return a.Attack
}