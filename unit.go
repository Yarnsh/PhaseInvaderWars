package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/Yarnsh/hippo/animation"
    "strconv"
)

var (
	unit_animations map[string]animation.Animation

	base_attacks = [][]float64{
		{0.55, 0.1, 0.85},
		{0.8, 0.55, 0.40},
		{0.15, 0.9, 0.30},
	}
)

func InitUnitData() {
	unit_animations, _ = animation.LoadAnimationMap("assets/units.json")
}

type Unit struct {
	which string
	x, y int
	strength float64
	army int
	actions int
}

func NewUnit(which string, x, y, army int) Unit {
	return Unit{
		which: which,
		x: x,
		y: y,
		army: army,
		strength: 1.0,
	}
}

func (u Unit) Draw(target *ebiten.Image, anim string, time float64) {
	done := ""

	if u.army == 0 && u.actions <= 0 {
		done = "_done"
	}

	unit_animations[u.which + "_" + anim + "_" + strconv.Itoa(u.army) + done].Draw(target, float64(u.x) * tileSizeF + (tileSizeF / 2.0), float64(u.y) * tileSizeF + tileSizeF, 1.0, time)

	// TODO: draw strength number if < 1.0
}

func (u Unit) CalculateDamage(o Unit, tacmap TacticalMap) (float64, float64) {
	ut := nameToType(u.which)
	ot := nameToType(o.which)
	base_attack := base_attacks[ut][ot]
	base_counterattack := base_attacks[ot][ut]

	u_str := u.strength
	o_str := o.strength

	o_str -= base_attack * u_str * tacmap.tiles[o.x][o.y].defense
	if o_str < 0.0 {
		o_str = 0.0
	}

	u_str -= base_counterattack * o_str * tacmap.tiles[u.x][u.y].defense
	if u_str < 0.0 {
		u_str = 0.0
	}

	return u_str, o_str
}

func (u Unit) GetMoves() int {
	switch u.which {
	case "infantry":
		return 3
	case "tank":
		return 6
	case "antitank":
		return 3
	}
	return 0
}

func nameToType(name string) int {
	switch name {
	case "infantry":
		return 0
	case "tank":
		return 1
	case "antitank":
		return 2
	}
	return 0
}
