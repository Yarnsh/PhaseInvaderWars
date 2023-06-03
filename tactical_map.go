package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/Yarnsh/hippo/utils"
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

func (t TacticalMap) GetMoves(u Unit) SearchNodes {
	nodes := SearchNodes{}
	start_node := SearchNode{
		pos: utils.IntPair{X: u.x, Y: u.y},
		moves_left: u.GetMoves(),
	}

	nposs, ncosts := start_node.GetNeighborsPositions(t.tiles)

	for idx, npos := range nposs {
		mleft := u.GetMoves() - ncosts[idx]
		if mleft < 0 {
			continue
		}
		should_search, replace_idx := shouldSearchNeighbor(npos, mleft, nodes)
		if should_search {
			new := SearchNode{
				pos: npos,
				moves_left: mleft,
				shortest_path: append(start_node.shortest_path, start_node.pos),
			}
			if replace_idx != -1 {
				nodes.nodes[replace_idx] = new
			} else {
				nodes.nodes = append(nodes.nodes, new)
			}

			t.getMovesHelper(mleft, new, &nodes)
		}
	}

	return nodes
}

func (t TacticalMap) getMovesHelper(moves int, us SearchNode, nodes *SearchNodes) {
	if moves <= 0 {
		return
	}

	nposs, ncosts := us.GetNeighborsPositions(t.tiles)

	for idx, npos := range nposs {
		mleft := moves - ncosts[idx]
		if mleft < 0 {
			continue
		}
		should_search, replace_idx := shouldSearchNeighbor(npos, mleft, *nodes)
		if should_search {
			new := SearchNode{
				pos: npos,
				moves_left: mleft,
				shortest_path: append(us.shortest_path, us.pos),
			}
			if replace_idx != -1 {
				nodes.nodes[replace_idx] = new
			} else {
				nodes.nodes = append(nodes.nodes, new)
			}

			t.getMovesHelper(mleft, new, nodes)
		}
	}
}

func shouldSearchNeighbor(npos utils.IntPair, moves_left int, nodes SearchNodes) (bool, int) { // returns if we should add to search tree, and what idx to replace if any
	for idx, node := range nodes.nodes {
		if node.pos.X == npos.X && node.pos.Y == npos.Y {
			if moves_left > node.moves_left {
				return true, idx
			} else {
				return false, -1
			}
		}
	}
	return true, -1
}

type SearchNode struct {
	pos utils.IntPair
	moves_left int
	shortest_path []utils.IntPair
}

func (n SearchNode) GetNeighborsPositions(tiles [][]Tile) ([]utils.IntPair, []int) {
	result := []utils.IntPair{}
	costs := []int{}
	if n.pos.X > 0 {
		if tiles[n.pos.X - 1][n.pos.Y].move_cost > 0 {
			result = append(result, utils.IntPair{X: n.pos.X - 1, Y: n.pos.Y})
			costs = append(costs, tiles[n.pos.X - 1][n.pos.Y].move_cost)
		}
	}
	if n.pos.X < tacMapWidth - 1 {
		if tiles[n.pos.X + 1][n.pos.Y].move_cost > 0 {
			result = append(result, utils.IntPair{X: n.pos.X + 1, Y: n.pos.Y})
			costs = append(costs, tiles[n.pos.X + 1][n.pos.Y].move_cost)
		}
	}
	if n.pos.Y > 0 {
		if tiles[n.pos.X][n.pos.Y - 1].move_cost > 0 {
			result = append(result, utils.IntPair{X: n.pos.X, Y: n.pos.Y - 1})
			costs = append(costs, tiles[n.pos.X][n.pos.Y - 1].move_cost)
		}
	}
	if n.pos.Y < tacMapHeight - 1 {
		if tiles[n.pos.X][n.pos.Y + 1].move_cost > 0 {
			result = append(result, utils.IntPair{X: n.pos.X, Y: n.pos.Y + 1})
			costs = append(costs, tiles[n.pos.X][n.pos.Y + 1].move_cost)
		}
	}

	return result, costs
}

type SearchNodes struct {
	nodes []SearchNode
}
