package main

import (
	"embed"
	"encoding/json"
	"image"
	"io"

	_ "image/png"

	"io/fs"
	// "path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScalingOperation string

const (
	SO_ADD ScalingOperation = "+"
	SO_SUB ScalingOperation = "-"
	SO_MUL ScalingOperation = "*"
	SO_DIV ScalingOperation = "/"
)

type ManifestScaling struct {
	// The base number
	Base int `json:"base"`
	// Wheather to scale at all
	Scale bool `json:"scale"`
	// The operation used to scale
	Opeartion ScalingOperation `json:"operation"`
	// The other side of the operation
	Factor float64 `json:"factor"`
	// Apply every __ ticks
	Interval int `json:"interval"`
	
	Max int `json:"max"`
	Min int `json:"min"`
}

type EnemyShape struct {
	Width int `json:"width"`
	Height int `json:"height"`
	IsCircle int `json:"circle"`
}

type EnemySummonBase struct {
	ID string `json:"id"`
	ShapeInfo EnemyShape `json:"shape"`

	Health ManifestScaling `json:"hp"`
	Speed  ManifestScaling `json:"speed"`
	Damage ManifestScaling `json:"dmg"`

	RewardXP ManifestScaling `json:"rewardXP"`
	RewardGold ManifestScaling `json:"rewardGold"`
	RewardItemChance ManifestScaling `json:"rewardItemChance"`
	RewardItemRarity ItemRarity `json:"rewardItemRarity"`
}

type EnemySummonResourcePack struct {
	Death []*ebiten.Image
	Sprite *ebiten.Image

	EnemySummonBase
}

type ManifestEnemySummonBase struct {
	// The top left coordinate of the texture
	Location Vector `json:"location"`
	// Has to be max. 30 for fps, then when 30%n = 0.
	DeathAnimation []Vector `json:"deathAnimation"`
	EnemySummonBase
}

type ResourcePack struct {
	BG *ebiten.Image
	BGSize int
	
	Protag *ebiten.Image
	Goal *ebiten.Image

	Enemies []EnemySummonResourcePack

	SpawnRate ManifestScaling
	SpawnAmount ManifestScaling

	GoalInfo GoalSummonInfo
}

type GoalSummonInfo struct {
	HP int `json:"HP"`
	Speed int `json:"speed"`
}

type Manifest struct {
	Enemies []ManifestEnemySummonBase `json:"enemies"`
	SpawnRate ManifestScaling `json:"spawnRate"`
	SpawnAmount ManifestScaling `json:"spawnAmount"`

	BGLocs []Vector `json:"backgroundLocations"`
	BGSize int `json:"backgroundSize"`

	Goal GoalSummonInfo `json:"goal"`
}

//go:embed resources
var textureFiles embed.FS

func loadFile(name string) fs.File {
	// even windows has '/'
	t, err := textureFiles.Open("resources/" + name)
	PanicIfErr(err)
	return t
}

func loadTexture(f fs.File) *ebiten.Image {
	img, _, err := image.Decode(f)
	PanicIfErr(err)
	return ebiten.NewImageFromImage(img)
}

func LoadResourcePack(world string) ResourcePack {
	textureFile := loadFile(world + ".png")
	atlas := loadTexture(textureFile)
	manifestFile := loadFile(world + ".manifest.json")
	manBytes, err := io.ReadAll(manifestFile)
	PanicIfErr(err)
	manifest := Manifest{}
	json.Unmarshal(manBytes, &manifest)
	bg := ebiten.NewImage(SCREEN_DIMENSIONS+manifest.BGSize*2, SCREEN_DIMENSIONS+manifest.BGSize*2)

	tiles := []*ebiten.Image{}

	for _, coords := range manifest.BGLocs {
		tiles = append(tiles, ebiten.NewImageFromImage(
			atlas.SubImage(
				image.Rect(
					coords[0], coords[1],
					coords[0] + manifest.BGSize, coords[1] + manifest.BGSize,
				),
			),
		))
	}
	
	for x := 0; x < SCREEN_DIMENSIONS + 2*manifest.BGSize; x += manifest.BGSize {
		for y := 0; y < SCREEN_DIMENSIONS + 2*manifest.BGSize; y += manifest.BGSize {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			tile := tiles[RandomInt(0, len(tiles)-1)]
			bg.DrawImage(tile, op)
		}
	}

	protag := loadTexture(loadFile("protag.png"))
	goal := loadTexture(loadFile("goal.png"))

	enemies := []EnemySummonResourcePack{}

	for _, e := range manifest.Enemies {
		sprite := ebiten.NewImageFromImage(
			atlas.SubImage(
				image.Rect(
					e.Location[0],
					e.Location[1],
					e.Location[0] + e.ShapeInfo.Width,
					e.Location[1] + e.ShapeInfo.Height,
				),
			),
		)
		
		death := []*ebiten.Image{}

		for _, dimg := range e.DeathAnimation {
			death = append(death, ebiten.NewImageFromImage(
				atlas.SubImage(
					image.Rect(
						dimg[0],
						dimg[1],
						dimg[0] + e.ShapeInfo.Width,
						dimg[1] + e.ShapeInfo.Height,
					),
				),
			))
		}

		enemies = append(enemies, EnemySummonResourcePack{
			Death:           death,
			Sprite:          sprite,
			EnemySummonBase: e.EnemySummonBase,
		})
	}

	ProtagBox = image.Rect(CenterCoords[0]-protag.Bounds().Dx()/2, CenterCoords[1]-protag.Bounds().Dy()/2, CenterCoords[0]+protag.Bounds().Dx()/2, CenterCoords[1]+protag.Bounds().Dy()/2)

	return ResourcePack{
		BG:        bg,
		BGSize: manifest.BGSize,

		Protag:    protag,
		Goal:      goal,
		Enemies:   enemies,

		SpawnRate: manifest.SpawnRate,
		SpawnAmount: manifest.SpawnAmount,
		
		GoalInfo: manifest.Goal,
	}
}
