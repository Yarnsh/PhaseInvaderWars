package main

import (
    "github.com/Yarnsh/hippo/animation"
    "github.com/Yarnsh/hippo/input"
    "github.com/Yarnsh/hippo/utils"
    "strings"
    "io/fs"
	"path/filepath"
)

var (
	TP_ Tile
	TF_ Tile
	TM_ Tile
	TO_ Tile
	TH0 Tile
	TF0 Tile
	TH1 Tile
	TF1 Tile
	TH2 Tile
	TF2 Tile
	TH3 Tile
	TF3 Tile

	tac_map_1 TacticalMap
	tac_map_2 TacticalMap
	tac_map_3 TacticalMap
)

func InitTacMapData() {
	tiles, _ := animation.LoadAnimationMap("assets/tiles.json")
	backgrounds, _ := animation.LoadAnimationMap("assets/backgrounds.json")

	TP_ = Tile{
		visual: tiles["plains"],
		battle_visual: backgrounds["plains"],
		defense: 1.0,
		move_cost: 1,
	}

	TF_ = Tile{
		visual: tiles["forest"],
		battle_visual: backgrounds["forest"],
		defense: 0.8,
		move_cost: 2,
	}

	TM_ = Tile{
		visual: tiles["mountain"],
		battle_visual: backgrounds["mountain"],
		defense: 0.5,
		move_cost: 3,
	}

	TO_ = Tile{
		visual: tiles["ocean"],
		battle_visual: backgrounds["ocean"],
		defense: 1.0,
		move_cost: -1,
	}

	TH0 = Tile{
		visual: tiles["hq0"],
		battle_visual: backgrounds["hq"],
		defense: 0.5,
		move_cost: 1,
		p1_hq: true,
	}

	TF0 = Tile{
		visual: tiles["factory0"],
		battle_visual: backgrounds["factory"],
		defense: 1.0,
		move_cost: 1,
		p1_fac: true,
	}

	TH1 = Tile{
		visual: tiles["hq1"],
		battle_visual: backgrounds["hq"],
		defense: 0.5,
		move_cost: 1,
		p2_hq: true,
	}

	TF1 = Tile{
		visual: tiles["factory1"],
		battle_visual: backgrounds["factory"],
		defense: 1.0,
		move_cost: 1,
		p2_fac: true,
	}

	TH2 = Tile{
		visual: tiles["hq2"],
		battle_visual: backgrounds["hq"],
		defense: 0.5,
		move_cost: 1,
		p2_hq: true,
	}

	TF2 = Tile{
		visual: tiles["factory2"],
		battle_visual: backgrounds["factory"],
		defense: 1.0,
		move_cost: 1,
		p2_fac: true,
	}

	TH3 = Tile{
		visual: tiles["hq2"],
		battle_visual: backgrounds["hq"],
		defense: 0.5,
		move_cost: 1,
		p2_hq: true,
	}

	TF3 = Tile{
		visual: tiles["factory3"],
		battle_visual: backgrounds["factory"],
		defense: 1.0,
		move_cost: 1,
		p2_fac: true,
	}

	tac_map_1 = csvToTacticalMap("assets/maps/tac_map_1.csv")
	tac_map_2 = csvToTacticalMap("assets/maps/tac_map_2.csv")
	tac_map_3 = csvToTacticalMap("assets/maps/tac_map_3.csv")
}

func csvToTacticalMap(file_path string) TacticalMap {
	dat, _ := fs.ReadFile(EmbeddedFileSystem, filepath.ToSlash(file_path))
    sdat := string(dat)
    split_dat := strings.Split(sdat, ",")

    map_slice := make([][]Tile, tacMapWidth)
	for i := range map_slice {
	    map_slice[i] = make([]Tile, tacMapHeight)
	}
    result := TacticalMap{
    	tiles: map_slice,
    }

    for x := 0; x < tacMapWidth; x++ {
	    for y := 0; y < tacMapHeight; y++ {
	    	switch split_dat[(y * tacMapWidth) + x] {
	    	case "0":
	    		result.tiles[x][y] = TO_
	    	case "2":
	    		result.tiles[x][y] = TP_
    		case "3":
    			result.tiles[x][y] = TF_
    		case "4":
    			result.tiles[x][y] = TM_

    		case "0.4":
    			result.tiles[x][y] = TH0
    			result.p1_hq = utils.IntPair{X: x, Y: y}
    		case "1.4":
    			result.tiles[x][y] = TF0
    			result.p1_factories = append(result.p1_factories, utils.IntPair{X: x, Y: y})

    		case "0.5":
    			result.tiles[x][y] = TH1
    			result.p2_hq = utils.IntPair{X: x, Y: y}
    		case "1.5":
    			result.tiles[x][y] = TF1
    			result.p2_factories = append(result.p2_factories, utils.IntPair{X: x, Y: y})

    		case "0.6":
    			result.tiles[x][y] = TH2
    			result.p2_hq = utils.IntPair{X: x, Y: y}
    		case "1.6":
    			result.tiles[x][y] = TF2
    			result.p2_factories = append(result.p2_factories, utils.IntPair{X: x, Y: y})

    		case "0.7":
    			result.tiles[x][y] = TH3
    			result.p2_hq = utils.IntPair{X: x, Y: y}
    		case "1.7":
    			result.tiles[x][y] = TF3
    			result.p2_factories = append(result.p2_factories, utils.IntPair{X: x, Y: y})
	    	}
		}
    }

    return result
}

func CreateTacticalGame(level int, player_control input.InputActionHandler) TacticalGame {
	if level == 0 {
		tacai := NewTacticalAI([]string{"infantry","infantry","infantry","antitank","tank"}, 2)

	    result := NewTacticalGame(tac_map_1, player_control, 1, tacai)

	    result.AddUnit("infantry", 7, 13, 0)
	    result.AddUnit("tank", 14, 16, 0)
	    result.AddUnit("antitank", 8, 15, 0)

	    result.AddUnit("infantry", 23, 13, 1)
	    result.AddUnit("tank", 24, 16, 1)
	    result.AddUnit("antitank", 28, 15, 1)

	    result.RefreshP1()
	    result.RefreshP2()

	    return result
	}
	if level == 1 {
		tacai := NewTacticalAI([]string{"tank","infantry","infantry","infantry","antitank","antitank","antitank"}, 0)

	    result := NewTacticalGame(tac_map_2, player_control, 2, tacai)

	    result.AddUnit("infantry", 19, 25, 0)
	    result.AddUnit("infantry", 17, 25, 0)
	    result.AddUnit("infantry", 21, 25, 0)

	    result.AddUnit("infantry", 17, 7, 2)
	    result.AddUnit("tank", 19, 7, 2)
	    result.AddUnit("antitank", 21, 7, 2)

	    result.RefreshP1()
	    result.RefreshP2()

	    return result
	}
	if level == 2 {
		tacai := NewTacticalAI([]string{"antitank","antitank","infantry","infantry","infantry","antitank","tank"}, 4)

	    result := NewTacticalGame(tac_map_3, player_control, 3, tacai)

	    result.AddUnit("tank", 22, 26, 0)

	    result.AddUnit("infantry", 16, 9, 3)
	    result.AddUnit("tank", 0, 29, 3)

	    result.RefreshP1()
	    result.RefreshP2()

	    return result
	}

	return TacticalGame{}
}
