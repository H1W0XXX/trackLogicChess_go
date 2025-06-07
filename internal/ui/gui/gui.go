package gui

import (
	"trackLogicChess/internal/assets"
	"trackLogicChess/internal/game"
)

var (
	marbleA = assets.LoadPNG("images/marbleA.png")
	marbleB = assets.LoadPNG("images/marbleB.png")
)

type GUI struct {
	anim animator
}

func NewApp(gs *game.GameState, ai bool) *App {
	return &App{
		state: gs,
		useAI: ai,
		imgA:  marbleA,
		imgB:  marbleB,
	}
}
