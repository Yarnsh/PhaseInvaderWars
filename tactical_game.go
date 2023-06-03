package main

import (
    "github.com/hajimehoshi/ebiten/v2"
)

type TacticalGame struct {
    tacmap TacticalMap
}

func NewTacticalGame(tacmap TacticalMap) TacticalGame {
    // TODO
    result := TacticalGame{
        tacmap: tacmap,
    }

    return result
}
func (g *TacticalGame) Update() error {
    // TODO

    return nil
}
func (g *TacticalGame) Draw(screen *ebiten.Image) {
    // TODO
    g.tacmap.Draw(screen)
}
func (g *TacticalGame) Layout(outsideWidth, outsideHeight int) (int, int) {
    return outsideWidth, outsideHeight
}

func (g TacticalGame) GetResult() (bool, bool) {
    // TODO
    return false, false
}

func (g TacticalGame) GetMap() TacticalMap {
    return g.tacmap
}

func (g TacticalGame) GetMoveMap() TacticalMap {
    // TODO
    return g.tacmap
}
