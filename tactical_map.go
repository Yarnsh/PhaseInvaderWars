package main

import (
    "github.com/hajimehoshi/ebiten/v2"
)

type TacticalMap struct {
	tiles [][]Tile
}

func CloneMap(tm TacticalMap) TacticalMap {
	result := TacticalMap{}
	for x := 0; x < tacMapWidth; x++ {
		for y := 0; y < tacMapHeight; y++ {
			result.tiles[x][y] = tm.tiles[x][y]
		}
	}

	return result
}

func (t TacticalMap) Draw(target *ebiten.Image) {
	for x := 0; x < tacMapWidth; x++ {
		for y := 0; y < tacMapHeight; y++ {
			t.tiles[x][y].Draw(target, x, y)
		}
	}
}
