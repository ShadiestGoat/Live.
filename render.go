package main

import (
	"image/color"
	"io"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

const SCREEN_DIMENSIONS = 1024
const DEBUG_W = 8

type shadowInfo struct {
	Sat    float64
	Bright float64
	Offset float64
}

var ShadowSl = []shadowInfo{
	{
		Sat:    0.5,
		Bright: 0.5,
		Offset: 4,
	},
	{
		Sat:    0.7,
		Bright: 0.7,
		Offset: 2,
	},
}

var debugLevels = []DebugCoords{}

type DebugCoords struct {
	IsY bool
	Color color.Color
	Level int
}

var face font.Face

var bestTime = 0

func init() {
	f := loadFile("font.ttf")
	b, err := io.ReadAll(f)
	PanicIfErr(err)
	sf, err := sfnt.Parse(b)
	PanicIfErr(err)
	tt, err := opentype.NewFace(sf, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	face = tt
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	
	op.GeoM.Translate(float64(g.BGOffset[0])-float64(g.Resources.BGSize), float64(g.BGOffset[1])-float64(g.Resources.BGSize))
	screen.DrawImage(g.Resources.BG, op)

	for _, enemy := range g.Enemies {
		opMons := &ebiten.DrawImageOptions{}
		b := enemy.Box(*g)
		width := float64(enemy.Sprite(*g).Bounds().Dx())
		opMons.GeoM.Translate(float64(b.Min.X), float64(b.Min.Y))
		
		screen.DrawImage(enemy.Sprite(*g), opMons)
		
		DrawHP(screen, width, float64(b.Min.X), float64(b.Min.Y)-10, enemy.HP(), enemy.HPMax(*g))
	}

	for _, info := range ShadowSl {
		opShadow := &ebiten.DrawImageOptions{}
		opShadow.GeoM.Translate(float64(ProtagBox.Min.X)+float64(curMovementVector[0])*info.Offset, 
								float64(ProtagBox.Min.Y)+float64(curMovementVector[1])*info.Offset,
		)
		opShadow.ColorM.ChangeHSV(0, info.Sat, info.Bright)
		screen.DrawImage(g.Resources.Protag, opShadow)
	}

	opProt := &ebiten.DrawImageOptions{}
	opProt.GeoM.Translate(float64(ProtagBox.Min.X), float64(ProtagBox.Min.Y))
	screen.DrawImage(g.Resources.Protag, opProt)

	for _, ability := range g.Protag.Abilities {
		ability.Draw(screen, *g)
	}
	
	timeBound := TimeSTR(g.Time)
	size := text.BoundString(face, timeBound)
	text.Draw(screen, timeBound, face, SCREEN_DIMENSIONS/2-size.Dx()/2, 50, color.RGBA{255,255,255,255})

	if bestTime != 0 {
		text.Draw(screen, TimeSTR(bestTime), face, 25, SCREEN_DIMENSIONS-25, color.RGBA{255,255,255,255})
	}

	opIcon := &ebiten.DrawImageOptions{}
	opIcon.GeoM.Translate(16, 16)
	if g.Protag.Abilities[AU_SLASH].Cooldown() != 0 {
		opIcon.ColorM.ChangeHSV(0, 0.5, 0.5)
		circ := vector.Path{}
		circ.MoveTo(32,32)
		circ.Arc(32, 32, 18, 6*math.Pi/4, 6*math.Pi/4+2*math.Pi*float32(g.Protag.Abilities[AU_SLASH].Cooldown())/float32(g.Protag.Abilities[AU_SLASH].MaxCooldown()), vector.Clockwise)
		vs, is := circ.AppendVerticesAndIndicesForFilling(nil, nil)
		op := &ebiten.DrawTrianglesOptions{FillRule: ebiten.EvenOdd}
		screen.DrawTriangles(vs, is, emptyImage, op)
	}

	screen.DrawImage(slashIcon, opIcon)

	DrawHP(screen, float64(ProtagBox.Dx()), float64(ProtagBox.Min.X), float64(ProtagBox.Max.Y)+10, g.Protag.HP(), g.Protag.HPMax(*g))

	if g.IsPaused {
		ebitenutil.DrawRect(screen, 0, 0, SCREEN_DIMENSIONS, SCREEN_DIMENSIONS, color.RGBA{
			0,
			0,
			0,
			191,
		})
	}
}

var slashIcon *ebiten.Image

// no idea what this does but the example said to use it so

var (
	emptyImage = ebiten.NewImage(1, 1)
	// emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	slashIcon = loadTexture(loadFile("slashIcon.png"))
	emptyImage.Fill(color.White)
}