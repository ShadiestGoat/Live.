package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"time"
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
	// if max-min <= 0 {
	// 	tmp := min
	// 	min = max
	// 	max = tmp
	// }
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
	*g = Game{
		Resources: &r,
		BGOffset:  Vector{},
		IsPaused:  false,
		Time:      0,
		Protag:    Protag{
			MaxHP:    100,
			HP:       100,
			RegenFreq: 120,

			IVTicks: 0,
			MaxIVTicks: 35,

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

					CurExecutionTick: 0,
					DspExecutionTick: 20,
					MaxExecutionTick: 45,
					
					Sprite:           slashTexture,
				},
			},
			// TODO: Add abilities!
		},
		Enemies:   []Enemy{},
	}
}