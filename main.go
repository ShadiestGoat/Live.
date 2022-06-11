package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_DIMENSIONS, SCREEN_DIMENSIONS
}

func main() {
	game := Game{}
	game.Restart()

	ebiten.SetWindowSize(SCREEN_DIMENSIONS, SCREEN_DIMENSIONS)
	ebiten.SetWindowTitle("Live.")
	ebiten.SetWindowFloating(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&game); err != nil {
		panic(err)
	}
}
