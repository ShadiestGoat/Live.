package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var CenterCoords = Vector{
	SCREEN_DIMENSIONS/2,
	SCREEN_DIMENSIONS/2,
}

var ProtagBox image.Rectangle

func ResolveVector(dirs map[Direction]bool) Vector {
	offset := Vector{}

	for dir, p := range dirs {
		if !p {
			continue
		}
		offset.Sub(DirEnum[dir])
	}
	
	return offset
}

func RandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func RandomFloat(min float64, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Add to current vector
func (v *Vector) Add(v2 Vector) {
	v[0] += v2[0]
	v[1] += v2[1]
}

func (v *Vector) Sub(v2 Vector) {
	v[0] -= v2[0]
	v[1] -= v2[1]
}

func Overflow(original int, max int) int {
	neg := original < 0

	original = int(math.Abs(float64(original)))

	for original >= max {
		original -= max
	}

	if neg {
		original = -original
	}

	return original
}

func (v Vector) DistanceToCenter() int {
	return v.DistanceTo(CenterCoords)
}

func (v Vector) DistanceTo(v2 Vector) int {
	dx := v[0]-v2[0]
	dy := v[1]-v2[1]
	return int(math.Round(math.Sqrt(float64(dx*dx+dy*dy))))
}

// returns 1 or -1, depending on if i > 0 or not
func GetSign(i int) int {
	return int(GetSignFloat(float64(i)))
}

// returns 1 or -1, depending on if i > 0 or not
func GetSignFloat(i float64) float64 {
	if i > 0 {
		return 1
	} else {
		return -1
	}
}


func AngleToCoords(angle float64, distance float64) Vector {
	return Vector{
		int(math.Round(math.Cos(angle)*distance)),
		int(math.Round(math.Sin(angle)*distance)),
	}
}

func (m ManifestScaling) Resolve(tick int) int {
	return int(math.Round(m.ResolveFloat(tick)))
}

func (m ManifestScaling) ResolveFloat(tick int) float64 {
	v := float64(m.Base)
	if !m.Scale {
		return v
	}
	
	times := math.Floor(float64(tick)/float64(m.Interval))

	switch m.Opeartion {
	case SO_ADD:
		v += m.Factor*times
	case SO_DIV:
		v /= math.Pow(m.Factor, times)
	case SO_MUL:
		v *= math.Pow(m.Factor, times)
	case SO_SUB:
		v -= m.Factor*times
	}

	if m.Max != 0 {
		if v > float64(m.Max) {
			v = float64(m.Max)
		}
	}

	if m.Min != 0 {
		if v < float64(m.Min) {
			v = float64(m.Min)
		}
	}

	return v
}



func RectCollision(col1 image.Rectangle, col2 image.Rectangle) bool {
	return rectCol(col1, col2) || rectCol(col2, col1)
}

func rectCol(col1 image.Rectangle, col2 image.Rectangle) bool {
	xMinCol := col1.Min.X <= col2.Min.X && col2.Min.X <= col1.Max.X
	xMaxCol := col1.Min.X <= col2.Max.X && col2.Max.X <= col1.Max.X
	
	yMinCol := col1.Min.Y <= col2.Min.Y && col2.Min.Y <= col1.Max.Y
	yMaxCol := col1.Min.Y <= col2.Max.Y && col2.Max.Y <= col1.Max.Y
	
	return (xMinCol || xMaxCol) && (yMinCol || yMaxCol)
}

func SlowTPS() {	
	time.Sleep(200 * time.Millisecond)
}

func TimeSTR(ticks int) string {
	secondsPassed := math.Floor(float64(ticks)/60)
	minsPassed := int(math.Floor(secondsPassed / 60))
	return fmt.Sprintf("%02d:%02d", minsPassed, int(secondsPassed)%60)
}

func (g *Game) Restart() {
	if g.Time > bestTime {
		bestTime = g.Time
	}
	r := LoadResourcePack("hell")
	goalCoords := AngleToCoords(RandomFloat(0, 2*math.Pi), 400)

	*g = Game{
		Resources: &r,
		BGOffset:  Vector{},
		IsPaused:  false,
		Time:      0,
		Protag:    Protag{
			_MaxHP:    100,
			_HP:       100,

			RegenFreq: 120,
			MaxRegenFreq: 120,

			_IVFrames: 0,
			MaxIVFrames: 35,

			XP:       0,
			Level:    0,
			Speed:    8,
			Luck:     0,
			Reach:    0,
			Coins: 0,
			Abilities: map[ActionUpgradeID]ActiveAbility{
				AU_SLASH: &IncreasingScaleAbility{
					_Cooldown:        0,
					_MaxCooldown:     2400,

					XP:               0,
					Level:            0,

					Attack: 40,
					
					CurExecutionTick: 0,
					DspExecutionTick: 20,
					MaxExecutionTick: 45,
					
					Sprite:           slashTexture,
				},
			},
			// TODO: Add abilities!
		},
		Enemies:   []Enemy{},
		Goal: Goal{
			_Coords:      goalCoords,
			_HP:          r.GoalInfo.HP,
			_IVFrames:    0,
			IVFramesMax: 40,
		},
	}
}

func DrawHP(dst *ebiten.Image, width, x, y float64, hp, hpmax int) {
	if hp != hpmax {
		perc := float64(hp)/float64(hpmax)
		greenW := math.Floor(perc*width)

		ebitenutil.DrawRect(dst, x, y, greenW, 8, color.RGBA{
			R: 39,
			G: 206,
			B: 0,
			A: 122,
		})

		ebitenutil.DrawRect(dst, x+greenW, y, width-greenW, 8, color.RGBA{
			R: 71,
			G: 71,
			B: 71,
			A: 122,
		})
	}
}

// A safe way to change a creature's hp. Returns true if its still alive.
func (g *Game) ChangeHP(hpChange int, a Alive) bool {
	newHP := a.HP() + hpChange
	if newHP <= 0 {
		return false
	}
	if newHP > a.HPMax(*g) {
		newHP = a.HPMax(*g)
	}
	a.setHP(newHP)
	return true
}

// Returns true if 'a' is still *alive* 
func (g *Game) Hurt(attack int, a Alive) bool {
	if a.IVFrames() != 0 {
		return true
	}
	ret := g.ChangeHP(-attack, a)
	if ret {
		a.ResetIVFrames()
	}
	return ret
}

func GetVelocityToCenter(a Alive, speed int) Vector {
	curCoords := a.Coords()
	diff := Vector{
		CenterCoords[0]-curCoords[0],
		CenterCoords[1]-curCoords[1],
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

	dy *= float64(speed)
	dx *= float64(speed)
	
	return Vector{
		int(math.Round(dy)),
		int(math.Round(dx)),
	}
}

func BasicBox(coords Vector, sprite *ebiten.Image) image.Rectangle {
	return image.Rect(coords[0]-sprite.Bounds().Dx()/2, coords[1]-sprite.Bounds().Dy()/2, coords[0]+sprite.Bounds().Dx()/2, coords[1]+sprite.Bounds().Dy()/2)
}