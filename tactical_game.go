package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/Yarnsh/hippo/input"
)

const (
    CURSOR_DELAY = 0.25
    CURSOR_DELAY2 = 0.05
)

type TacticalGame struct {
    time float64
    player_input input.InputActionHandler

    tacmap TacticalMap
    p1_units []Unit
    p2_units []Unit

    cx, cy int

    up_pressed, down_pressed, left_pressed, right_pressed float64
}

func NewTacticalGame(tacmap TacticalMap, player_input input.InputActionHandler) TacticalGame {
    // TODO
    result := TacticalGame{
        tacmap: tacmap,
        player_input: player_input,
    }

    return result
}
func (g *TacticalGame) Update() error {
    // TODO
    g.time += 1.0/60.0

    if g.player_input.IsActionJustPressed("left") {
        g.left_pressed = g.time + CURSOR_DELAY
        g.MoveCursor(-1, 0)
    }
    if g.player_input.IsActionJustPressed("right") {
        g.right_pressed = g.time + CURSOR_DELAY
        g.MoveCursor(1, 0)
    }
    if g.player_input.IsActionJustPressed("up") {
        g.up_pressed = g.time + CURSOR_DELAY
        g.MoveCursor(0, -1)
    }
    if g.player_input.IsActionJustPressed("down") {
        g.down_pressed = g.time + CURSOR_DELAY
        g.MoveCursor(0, 1)
    }

    if g.player_input.ActionPressedDuration("left") > 0 && g.time >= g.left_pressed {
        g.left_pressed = g.time + CURSOR_DELAY2
        g.MoveCursor(-1, 0)
    }
    if g.player_input.ActionPressedDuration("right") > 0 && g.time >= g.right_pressed {
        g.right_pressed = g.time + CURSOR_DELAY2
        g.MoveCursor(1, 0)
    }
    if g.player_input.ActionPressedDuration("up") > 0 && g.time >= g.up_pressed {
        g.up_pressed = g.time + CURSOR_DELAY2
        g.MoveCursor(0, -1)
    }
    if g.player_input.ActionPressedDuration("down") > 0 && g.time >= g.down_pressed {
        g.down_pressed = g.time + CURSOR_DELAY2
        g.MoveCursor(0, 1)
    }

    return nil
}
func (g *TacticalGame) Draw(screen *ebiten.Image) {
    // TODO
    g.tacmap.Draw(screen)

    for _, unit := range g.p1_units { // TODO: draw moving things differently
        unit.Draw(screen, "idle", g.time)
    }
    for _, unit := range g.p2_units { // TODO: draw moving things differently
        unit.Draw(screen, "idle", g.time)
    }

    // figure out if the cursor can click a thing right now
    hover := false
    found := false
    for _, unit := range g.p1_units {
        if unit.x == g.cx && unit.y == g.cy {
                found = true
            if unit.actions > 0 {
                hover = true
            }
            break
        }
    }
    if !found && g.tacmap.tiles[g.cx][g.cy].p1_fac {
        hover = true
    }

    if hover {
        ui_anims["cursor_hover"].Draw(screen, float64(g.cx) * tileSizeF + (tileSizeF / 2.0), float64(g.cy) * tileSizeF + tileSizeF, 1.0, g.time)
    } else {
        ui_anims["cursor"].Draw(screen, float64(g.cx) * tileSizeF + (tileSizeF / 2.0), float64(g.cy) * tileSizeF + tileSizeF, 1.0, g.time)
    }
}
func (g *TacticalGame) Layout(outsideWidth, outsideHeight int) (int, int) {
    return outsideWidth, outsideHeight
}

func (g *TacticalGame) AddUnit(which string, x, y, army int) {
    new := NewUnit(which, x, y, army)
    if army == 0 {
        g.p1_units = append(g.p1_units, new)
    } else {
        g.p2_units = append(g.p2_units, new)
    }
}

func (g *TacticalGame) MoveCursor(dx, dy int) {
    g.cx += dx
    if g.cx < 0 {
        g.cx = 0
    } else if g.cx >= tacMapWidth {
        g.cx = tacMapWidth-1
    }
    g.cy += dy
    if g.cy < 0 {
        g.cy = 0
    } else if g.cy >= tacMapHeight {
        g.cy = tacMapHeight-1
    }
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
