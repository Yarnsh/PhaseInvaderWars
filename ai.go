package main

import (
    "github.com/Yarnsh/hippo/utils"
    "math"
)

type TacticalAI struct {
	build_order []string
	build_repeat_idx int
	next_build int
}

func (ai TacticalAI) GetNextBuild() string {
	if len(ai.build_order) > 0 {
		return ai.build_order[ai.next_build]
	}
	return ""
}

func (ai *TacticalAI) IncrementNextBuild() {
	ai.next_build += 1
	if ai.next_build >= len(ai.build_order) {
		ai.next_build = ai.build_repeat_idx
	}
}

func costOfUnit(ut string) int {
	switch ut {
	case "infantry":
		return 1000
	case "tank":
		return 7000
	case "antitank":
		return 3000
	}
	return -1000
}

func (ai TacticalAI) GetBestMove(unit Unit, game TacticalGame) ([]utils.IntPair, utils.IntPair) {
	moves := game.movemap.GetMoves(unit)

	best_eval := -999999999.9
	best_path := []utils.IntPair{}
	best_attack := utils.IntPair{X: unit.x, Y: unit.y}

	for _, move := range moves.nodes {
		m := SearchNodes{nodes: []SearchNode{move}}
        attacks_n := getAdjacentNodes(m)
        attacks := []utils.IntPair{}
        for _, an := range attacks_n.nodes {
        	attacks = append(attacks, utils.IntPair{X: an.pos.X, Y: an.pos.Y})
        }
        attacks = append(attacks, utils.IntPair{X: unit.x, Y: unit.y})

        for _, attack := range attacks {
	        eval := ai.Evaluate(unit, move, attack, game)
	        if eval > best_eval {
	        	best_eval = eval
	        	best_path = append(move.shortest_path, move.pos)[1:]
	        	best_attack = attack
	        }
    	}
	}

	return best_path, best_attack
}

func (ai TacticalAI) Evaluate(us Unit, move SearchNode, attack utils.IntPair, game TacticalGame) float64 {
	// TODO: lots could be improved here, but we just do whatever to get it working for now
	eval_result := 0.0

	// position is better if its defensible
	eval_result += 1.0 / game.tacmap.tiles[move.pos.X][move.pos.Y].defense

	dist_to_hq := math.Abs(float64(game.tacmap.p1_hq.X - move.pos.X)) + math.Abs(float64(game.tacmap.p1_hq.Y - move.pos.Y))
	eval_result -= dist_to_hq * 0.85

	sastr := us.strength
	for _, unit := range game.p1_units {
        if unit.x == attack.X && unit.y == attack.Y {
        	sdstr := unit.strength
            astr, dstr := us.CalculateDamage(unit, game.tacmap)

            eval_result -= (sastr - astr) * 2.5
            eval_result += (sdstr - dstr) * 4.0
        }
    }

	return eval_result
}
