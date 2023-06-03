package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/Yarnsh/hippo/animation"
)

const (
	tileSize = 20
	tileSizeF = 20.0
)

type Tile struct {
	visual animation.Animation
	battle_visual animation.Animation
	defense float64
	move_cost int
	p1_hq, p2_hq, p1_fac, p2_fac bool
}

func (t Tile) Draw(target *ebiten.Image, x, y int) {
	t.visual.Draw(target, float64(x) * tileSizeF + (tileSizeF / 2.0), float64(y) * tileSizeF + tileSizeF, 1.0, 0.0)
}
