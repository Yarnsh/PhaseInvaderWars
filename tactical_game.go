package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/Yarnsh/hippo/input"
    "github.com/Yarnsh/hippo/utils"
    "math"
    "strconv"
    "fmt"
)

const (
    CURSOR_DELAY = 0.25
    CURSOR_DELAY2 = 0.05
    UNIT_MOVE_DELAY = 0.1

    MODE_CURSOR = 0
    MODE_MOVE_SELECT = 1
    MODE_FACTORY_SELECT = 2
    MODE_UNIT_MOVING = 3
    MODE_ATTACK_SELECT = 4
    MODE_END_TURN_MENU = 5
    MODE_AI_PLAYING = 6
    MODE_BATTLE = 7
    MODE_BATTLE_AI = 8
)

type TacticalGame struct {
    time float64
    player_input input.InputActionHandler
    mode int

    ai TacticalAI
    ai_building bool
    ai_next_unit int
    ai_best_attack utils.IntPair
    ai_need_to_attack bool
    p2_army int

    selected_unit Unit
    selected_unit_idx int
    selected_unit_player int
    unit_move_path []utils.IntPair
    unit_move_start float64
    selected_factory Tile

    tacmap TacticalMap
    movemap TacticalMap
    p1_units []Unit
    p2_units []Unit

    p1_money, p2_money int

    cx, cy int

    up_pressed, down_pressed, left_pressed, right_pressed float64

    menu_selection int

    current_battle BattleGame

    game_over bool
    we_won bool
}

func NewTacticalGame(tacmap TacticalMap, player_input input.InputActionHandler, p2_army int, ai TacticalAI) TacticalGame {
    result := TacticalGame{
        tacmap: tacmap,
        player_input: player_input,
        p2_army: p2_army,
        ai: ai,
    }

    return result
}

func (g *TacticalGame) StartAITurn() {
    g.ai_building = true
    g.ai_next_unit = 0
    g.ai_need_to_attack = false
}

func (g *TacticalGame) AIUpdate() bool { // returns if the AI is done
    if g.ai_building {
        next_unit := g.ai.GetNextBuild()
        if next_unit == "" || costOfUnit(next_unit) > g.p2_money {
            g.ai_building = false
            return false
        }

        for _, facpos := range g.tacmap.p2_factories {
            // check for empty factories
            fac_free := true
            for _, unit := range g.p1_units {
                if unit.x == facpos.X && unit.y == facpos.Y {
                    fac_free = false
                    break
                }
            }
            if !fac_free {
                continue
            }
            for _, unit := range g.p2_units {
                if unit.x == facpos.X && unit.y == facpos.Y {
                    fac_free = false
                    break
                }
            }
            if !fac_free {
                continue
            }

            g.AddUnit(next_unit, facpos.X, facpos.Y, g.p2_army)
            g.p2_money -= costOfUnit(next_unit)
            g.ai.IncrementNextBuild()
            return false
        }

        g.ai_building = false
        return false
    } else if g.ai_need_to_attack {
        for idx, unit := range g.p1_units {
            if unit.x == g.ai_best_attack.X && unit.y == g.ai_best_attack.Y {
                astr, dstr := g.p2_units[g.ai_next_unit].CalculateDamage(unit, g.tacmap)
                fmt.Println(astr, dstr)
                // switch to attack animation mode here
                g.mode = MODE_BATTLE_AI
                g.current_battle = NewBattleGame(g.p2_units[g.ai_next_unit], unit, g.tacmap)

                g.p2_units[g.ai_next_unit].strength = astr
                if astr <= 0.0 {
                    g.p2_units = append(g.p2_units[:g.ai_next_unit], g.p2_units[g.ai_next_unit+1:]...)
                    g.ai_next_unit -= 1
                }
                g.p1_units[idx].strength = dstr
                if dstr <= 0.0 {
                    g.p1_units = append(g.p1_units[:idx], g.p1_units[idx+1:]...)
                }
            }
        }

        g.ai_need_to_attack = false
        g.ai_next_unit += 1
        if len(g.p2_units) <= g.ai_next_unit {
            fmt.Println("doneatl ", len(g.p2_units), g.ai_next_unit)
            return true
        }

        return false

    } else {
        fmt.Println("thinkin of move for ", g.ai_next_unit)
        if len(g.p2_units) == 0 || len(g.p2_units) <= g.ai_next_unit {
            fmt.Println("done ", len(g.p2_units), g.ai_next_unit)
            return true
        }

        x := g.p2_units[g.ai_next_unit].x
        y := g.p2_units[g.ai_next_unit].y

        best_path, best_attack := g.ai.GetBestMove(g.p2_units[g.ai_next_unit], *g)

        if len(best_path) > 0 {
            g.unit_move_path = best_path
            g.mode = MODE_UNIT_MOVING
            g.unit_move_start = g.time
        }
        g.selected_unit = g.p2_units[g.ai_next_unit]
        g.selected_unit_idx = g.ai_next_unit
        g.selected_unit_player = 1
        g.ai_best_attack = best_attack
        g.ai_need_to_attack = (x != best_attack.X || y != best_attack.Y)

        if !g.ai_need_to_attack {
            g.ai_next_unit += 1
            fmt.Println("donend ", len(g.p2_units), g.ai_next_unit)
            if len(g.p2_units) <= g.ai_next_unit {
                return true
            }
        }

        return false
    }

    return true
}

func (g *TacticalGame) Update() error {
    g.movemap = g.tacmap.GetMovableMap(*g)

    g.time += 1.0/60.0

    if g.game_over {
        // some kinda end card probably
        return nil
    }

    // check for win/loss
    for _, unit := range g.p1_units {
        if unit.x == g.tacmap.p2_hq.X && unit.y == g.tacmap.p2_hq.Y {
            g.game_over = true
            g.we_won = true
        }
    }
    for _, unit := range g.p2_units {
        if unit.x == g.tacmap.p1_hq.X && unit.y == g.tacmap.p1_hq.Y {
            g.game_over = true
            g.we_won = false
        }
    }

    if g.mode == MODE_AI_PLAYING {
        if g.AIUpdate() {
            g.RefreshP1()
            g.mode = MODE_CURSOR
        }
        return nil
    }

    if g.mode == MODE_CURSOR || g.mode == MODE_MOVE_SELECT || g.mode == MODE_ATTACK_SELECT {
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
    }

    if g.mode != MODE_CURSOR && g.mode != MODE_UNIT_MOVING {
        if g.player_input.IsActionJustPressed("cancel") {
            if g.mode == MODE_ATTACK_SELECT {
                g.p1_units[g.selected_unit_idx].actions -= 1
            }
            g.mode = MODE_CURSOR
            return nil
        }
    }

    // everything after this should return

    if g.mode == MODE_END_TURN_MENU {
        if g.player_input.IsActionJustPressed("up") || g.player_input.IsActionJustPressed("down") {
            g.menu_selection += 1
            if g.menu_selection >= 2 {
                g.menu_selection = 0
            }
        }
        if g.player_input.IsActionJustPressed("accept") {
            if g.menu_selection == 0 {
                g.mode = MODE_AI_PLAYING
                g.RefreshP2()
            } else {
                g.mode = MODE_CURSOR
            }
        }
        return nil
    }

    if g.mode == MODE_FACTORY_SELECT {
        if g.player_input.IsActionJustPressed("up"){
            g.menu_selection -= 1
            if g.menu_selection < 0 {
                g.menu_selection = 2
            }
        }
        if g.player_input.IsActionJustPressed("down"){
            g.menu_selection += 1
            if g.menu_selection > 2 {
                g.menu_selection = 0
            }
        }
        if g.player_input.IsActionJustPressed("accept") {
            if g.menu_selection == 0 {
                if g.p1_money >= 1000 {
                    g.AddUnit("infantry", g.cx, g.cy, 0)
                    g.p1_money -= 1000
                    g.mode = MODE_CURSOR
                }
            } else if g.menu_selection == 1 {
                if g.p1_money >= 7000 {
                    g.AddUnit("tank", g.cx, g.cy, 0)
                    g.p1_money -= 7000
                    g.mode = MODE_CURSOR
                }
            } else if g.menu_selection == 2 {
                if g.p1_money >= 3000 {
                    g.AddUnit("antitank", g.cx, g.cy, 0)
                    g.p1_money -= 3000
                    g.mode = MODE_CURSOR
                }
            }
        }
        return nil
    }

    if g.mode == MODE_CURSOR {
        if g.player_input.IsActionJustPressed("accept") {
            for idx, unit := range g.p1_units {
                if unit.x == g.cx && unit.y == g.cy {
                    if unit.actions > 0 {
                        if unit.actions > 1 {
                            g.mode = MODE_MOVE_SELECT
                        } else {
                            g.mode = MODE_ATTACK_SELECT
                        }
                        g.selected_unit = unit
                        g.selected_unit_idx = idx
                        g.selected_unit_player = 0
                    }
                    return nil
                }
            }

            for _, fac_pos := range g.tacmap.p1_factories {
                if fac_pos.X == g.cx && fac_pos.Y == g.cy {
                    g.mode = MODE_FACTORY_SELECT
                    g.menu_selection = 0
                }
            }

            return nil
        } else if g.player_input.IsActionJustPressed("cancel") {
            g.mode = MODE_END_TURN_MENU
            g.menu_selection = 0
            return nil
        }
        return nil
    }

    if g.mode == MODE_MOVE_SELECT {
        if g.player_input.IsActionJustPressed("accept") {
            moves := g.movemap.GetMoves(g.selected_unit)
            for _, node := range moves.nodes {
                if node.pos.X == g.cx && node.pos.Y == g.cy {
                    g.unit_move_path = append(node.shortest_path, node.pos)[1:]
                    g.mode = MODE_UNIT_MOVING
                    g.unit_move_start = g.time
                    if g.selected_unit_player == 0 {
                        g.p1_units[g.selected_unit_idx].actions -= 1
                    } else {
                        g.p2_units[g.selected_unit_idx].actions -= 1 // this really should never run
                    }
                }
            }
        }
        return nil
    }

    if g.mode == MODE_ATTACK_SELECT {
        if g.player_input.IsActionJustPressed("accept") {
            moves := SearchNodes{nodes: []SearchNode{SearchNode{pos: utils.IntPair{X: g.selected_unit.x, Y: g.selected_unit.y}}}} // oh I hate this
            attacks := getAdjacentNodes(moves)
            for _, node := range attacks.nodes {
                if g.cx == node.pos.X && g.cy == node.pos.Y {
                    for idx, unit := range g.p2_units {
                        if unit.x == node.pos.X && unit.y == node.pos.Y {
                            astr, dstr := g.selected_unit.CalculateDamage(unit, g.tacmap)
                            fmt.Println(astr, dstr)
                            // switch to attack animation mode here
                            g.mode = MODE_BATTLE
                            g.current_battle = NewBattleGame(g.p1_units[g.selected_unit_idx], unit, g.tacmap)

                            g.p1_units[g.selected_unit_idx].strength = astr
                            if astr <= 0.0 {
                                g.p1_units = append(g.p1_units[:g.selected_unit_idx], g.p1_units[g.selected_unit_idx+1:]...)
                            }
                            g.p2_units[idx].strength = dstr
                            if dstr <= 0.0 {
                                g.p2_units = append(g.p2_units[:idx], g.p2_units[idx+1:]...)
                            }

                            g.p1_units[g.selected_unit_idx].actions -= 1

                            return nil
                        }
                    }
                }
            }
        }
        return nil
    }

    if g.mode == MODE_UNIT_MOVING {
        if len(g.unit_move_path) == 0 {
            g.mode = MODE_ATTACK_SELECT
            return nil
        }

        if g.unit_move_start + UNIT_MOVE_DELAY <= g.time {
            if g.selected_unit_player == 0 {
                g.p1_units[g.selected_unit_idx].x = g.unit_move_path[0].X
                g.p1_units[g.selected_unit_idx].y = g.unit_move_path[0].Y
                g.selected_unit = g.p1_units[g.selected_unit_idx]
            } else {
                g.p2_units[g.selected_unit_idx].x = g.unit_move_path[0].X
                g.p2_units[g.selected_unit_idx].y = g.unit_move_path[0].Y
                g.selected_unit = g.p2_units[g.selected_unit_idx]
            }
            g.unit_move_path = g.unit_move_path[1:]
            g.unit_move_start = g.time
        }

        if len(g.unit_move_path) == 0 {
            if g.selected_unit_player == 0 {
                g.mode = MODE_ATTACK_SELECT
            } else {
                g.mode = MODE_AI_PLAYING
            }
            return nil
        }
        
        return nil
    }

    if g.mode == MODE_BATTLE || g.mode == MODE_BATTLE_AI {
        g.current_battle.Update()
        done_battle, _ := g.current_battle.GetResult()
        if done_battle {
            if g.mode == MODE_BATTLE {
                g.mode = MODE_CURSOR
            } else {
                g.mode = MODE_AI_PLAYING
            }
        }
    }

    return nil
}
func (g *TacticalGame) Draw(screen *ebiten.Image) {
    if g.mode == MODE_BATTLE || g.mode == MODE_BATTLE_AI {
        g.current_battle.Draw(screen)
        return
    }

    g.tacmap.Draw(screen)

    if g.mode == MODE_UNIT_MOVING {
        for idx, unit := range g.p1_units { // TODO: draw moving things differently
            if len(g.unit_move_path) > 0 && g.selected_unit_player == 0 && idx == g.selected_unit_idx {
                if unit.x > g.unit_move_path[0].X {
                    unit.Draw(screen, "left", g.time)
                } else if unit.x < g.unit_move_path[0].X {
                    unit.Draw(screen, "right", g.time)
                } else if unit.y < g.unit_move_path[0].Y {
                    unit.Draw(screen, "down", g.time)
                } else {
                    unit.Draw(screen, "up", g.time)
                }
            } else {
                unit.Draw(screen, "idle", g.time)
                drawUnitStr(screen, unit)
            }
        }
        for idx, unit := range g.p2_units { // TODO: draw moving things differently
            if len(g.unit_move_path) > 0 && g.selected_unit_player == 1 && idx == g.selected_unit_idx {
                if unit.x > g.unit_move_path[0].X {
                    unit.Draw(screen, "left", g.time)
                } else if unit.x < g.unit_move_path[0].X {
                    unit.Draw(screen, "right", g.time)
                } else if unit.y < g.unit_move_path[0].Y {
                    unit.Draw(screen, "down", g.time)
                } else {
                    unit.Draw(screen, "up", g.time)
                }
            } else {
                unit.Draw(screen, "idle", g.time)
                drawUnitStr(screen, unit)
            }
        }
    } else {
        for _, unit := range g.p1_units { // TODO: draw moving things differently
            unit.Draw(screen, "idle", g.time)
            drawUnitStr(screen, unit)
        }
        for _, unit := range g.p2_units { // TODO: draw moving things differently
            unit.Draw(screen, "idle", g.time)
            drawUnitStr(screen, unit)
        }
    }

    // figure out if the cursor can click a thing right now
    hover := false
    if g.mode == MODE_CURSOR {
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
    }

    if g.mode != MODE_AI_PLAYING {
        if hover {
            ui_anims["cursor_hover"].Draw(screen, float64(g.cx) * tileSizeF + (tileSizeF / 2.0), float64(g.cy) * tileSizeF + tileSizeF, 1.0, g.time)
        } else {
            ui_anims["cursor"].Draw(screen, float64(g.cx) * tileSizeF + (tileSizeF / 2.0), float64(g.cy) * tileSizeF + tileSizeF, 1.0, g.time)
        }
    }

    if g.mode == MODE_MOVE_SELECT {
        moves := g.movemap.GetMoves(g.selected_unit)
        //attacks := getAdjacentNodes(moves)
        for _, node := range moves.nodes {
            ui_anims["walk_tile"].Draw(screen, float64(node.pos.X) * tileSizeF + (tileSizeF / 2.0), float64(node.pos.Y) * tileSizeF + tileSizeF, 1.0, g.time)
        }
        //for _, node := range attacks.nodes {
        //    ui_anims["attack_tile"].Draw(screen, float64(node.pos.X) * tileSizeF + (tileSizeF / 2.0), float64(node.pos.Y) * tileSizeF + tileSizeF, 1.0, g.time)
        //}
    }

    if g.mode == MODE_ATTACK_SELECT {
        moves := SearchNodes{nodes: []SearchNode{SearchNode{pos: utils.IntPair{X: g.selected_unit.x, Y: g.selected_unit.y}}}} // oh I hate this
        attacks := getAdjacentNodes(moves)
        for _, node := range attacks.nodes {
            ui_anims["attack_tile"].Draw(screen, float64(node.pos.X) * tileSizeF + (tileSizeF / 2.0), float64(node.pos.Y) * tileSizeF + tileSizeF, 1.0, g.time)
        }
    }

    if g.mode == MODE_END_TURN_MENU {
        drawEndTurnPrompt(screen, g.menu_selection)
    }

    if g.mode == MODE_FACTORY_SELECT {
        drawFactoryPrompt(screen, g.menu_selection, g.p1_money)
    }
}

func drawUnitStr(screen *ebiten.Image, unit Unit) {
    if unit.strength < 1.0 {
        num := int(math.Floor(unit.strength * 10.0))
        if num < 1 {
            num = 1
        }
        basic_font[16 + num].Draw(screen, float64(unit.x) * 20.0 + 16.0, float64(unit.y) * 20.0 + 20.0, 1.0, 0.0)
    }
}

func (g *TacticalGame) Layout(outsideWidth, outsideHeight int) (int, int) {
    return outsideWidth, outsideHeight
}

func (g *TacticalGame) RefreshP1() {
    for idx, _ := range g.p1_units {
        g.p1_units[idx].actions = 2
    }
    g.p1_money += 1500
}

func (g *TacticalGame) RefreshP2() {
    for idx, _ := range g.p2_units {
        g.p2_units[idx].actions = 2
    }
    g.p2_money += 1500
    g.StartAITurn()
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
    return g.game_over, g.we_won
}

func drawEndTurnPrompt(screen *ebiten.Image, selection int) {
    drawCenteredTextLine(screen, "END TURN?", 300 - 30)
    if selection == 0 {
        drawCenteredTextLine(screen, ">Yeah<", 300)
    } else {
        drawCenteredTextLine(screen, "Yeah", 300)
    }

    if selection == 1 {
        drawCenteredTextLine(screen, ">Nah<", 300 + 30)
    } else {
        drawCenteredTextLine(screen, "Nah", 300 + 30)
    }
}

func drawFactoryPrompt(screen *ebiten.Image, selection int, money int) {
    drawCenteredTextLine(screen, "BUILD UNIT", 300 - 40)
    if selection == 0 {
        drawCenteredTextLine(screen, ">Infantry $1000<", 300)
    } else {
        drawCenteredTextLine(screen, "Infantry $1000", 300)
    }

    if selection == 1 {
        drawCenteredTextLine(screen, ">Tank $7000<", 300 + 30)
    } else {
        drawCenteredTextLine(screen, "Tank $7000", 300 + 30)
    }

    if selection == 2 {
        drawCenteredTextLine(screen, ">Anti-tank $3000<", 300 + 60)
    } else {
        drawCenteredTextLine(screen, "Anti-tank $3000", 300 + 60)
    }

    drawCenteredTextLine(screen, "HAVE: $" + strconv.Itoa(money), 300 + 120)
}

func drawCenteredTextLine(screen *ebiten.Image, text string, height float64) {
    chars := len(text)
    left_offset := 400 - float64((chars * 16) / 2)
    for column, char := range []rune(text) {
        // For some reason we are a row off, hence "- 32", might be a hippo issue
        basic_font[char - 32].Draw(
            screen,
            left_offset + float64((column * 16)),
            height,
            2.0,
            0.0) // not pretty to hard code all this but whatever
    }
}
