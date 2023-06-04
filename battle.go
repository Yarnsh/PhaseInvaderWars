package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/Yarnsh/hippo/animation"
    "math"
    "strconv"
)

type BattleGame struct {
	attacker_type, defender_type string
	attacker_army, defender_army int
	start_astr, start_dstr float64
	end_astr, end_dstr float64
	attacker_bg, defender_bg animation.Animation

	attacker_image, defender_image, portflip_image *ebiten.Image

	time float64
}

func NewBattleGame(a, d Unit, m TacticalMap) BattleGame {
	r := BattleGame{}

	r.attacker_type = a.which
	r.defender_type = d.which
	r.attacker_army = a.army
	r.defender_army = d.army

	r.attacker_bg = m.tiles[a.x][a.y].battle_visual
	r.defender_bg = m.tiles[d.x][d.y].battle_visual

	r.start_astr = a.strength
	r.start_dstr = d.strength
	r.end_astr, r.end_dstr = a.CalculateDamage(d, m)

	r.attacker_image = ebiten.NewImage(400, 600)
	r.defender_image = ebiten.NewImage(400, 600)
	r.portflip_image = ebiten.NewImage(800, 600)

	return r
}

func (g *BattleGame) Update() error {
    g.time += 1.0/60.0

    //0.5 second run in
    //1.5 second attack
    //1.5 second defense
    //0.5 second delay before end

    return nil
}
func (g *BattleGame) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	g.attacker_bg.Draw(g.attacker_image, 200, 600, 1.0, g.time)
	g.defender_bg.Draw(g.defender_image, 200, 600, 1.0, g.time)

	// TODO: draw unit animations
	// attacker things
	if g.start_astr > 0.8 {
		ktime := 10.0
		if g.end_astr <= 0.8 {
			ktime = 3.0
		}
		g.DrawAttacker(200.0, 300.0, 0.0, ktime)
	}
	if g.start_astr > 0.6 {
		ktime := 10.0
		if g.end_astr <= 0.6 {
			ktime = 2.6
		}
		g.DrawAttacker(150.0, 350.0, 0.2, ktime)
	}
	if g.start_astr > 0.4 {
		ktime := 10.0
		if g.end_astr <= 0.4 {
			ktime = 2.2
		}
		g.DrawAttacker(250.0, 400.0, 0.4, ktime)
	}
	if g.start_astr > 0.2 {
		ktime := 10.0
		if g.end_astr <= 0.2 {
			ktime = 3.2
		}
		g.DrawAttacker(230.0, 450.0, 0.5, ktime)
	}
	if g.start_astr > 0.0 {
		ktime := 10.0
		if g.end_astr <= 0.0 {
			ktime = 2.2
		}
		g.DrawAttacker(170.0, 500.0, 0.7, ktime)
	}

	// defender things
	if g.start_dstr > 0.8 {
		ktime := 10.0
		if g.end_dstr <= 0.8 {
			ktime = 1.5
		}
		g.DrawDefender(200.0, 300.0, 0.0, ktime)
	}
	if g.start_dstr > 0.6 {
		ktime := 10.0
		if g.end_dstr <= 0.6 {
			ktime = 1.1
		}
		g.DrawDefender(150.0, 350.0, 0.2, ktime)
	}
	if g.start_dstr > 0.4 {
		ktime := 10.0
		if g.end_dstr <= 0.4 {
			ktime = 0.7
		}
		g.DrawDefender(250.0, 400.0, 0.4, ktime)
	}
	if g.start_dstr > 0.2 {
		ktime := 10.0
		if g.end_dstr <= 0.2 {
			ktime = 1.2
		}
		g.DrawDefender(230.0, 450.0, 0.5, ktime)
	}
	if g.start_dstr > 0.0 {
		ktime := 10.0
		if g.end_dstr <= 0.0 {
			ktime = 0.9
		}
		g.DrawDefender(170.0, 500.0, 0.7, ktime)
	}

    op.GeoM.Reset()
    screen.DrawImage(g.attacker_image, op)
	op.GeoM.Scale(-1.0, 1.0)
	op.GeoM.Translate(800, 0)
    screen.DrawImage(g.defender_image, op)
	ui_anims["battle_ui"].Draw(screen, 400, 600, 1.0, g.time)

	atk, def := g.GetCurrentStrs()

	if atk >= 10 {
		basic_font[16 + 1].Draw(screen, 306.0, 95.0, 3.0, 0.0)
		basic_font[16].Draw(screen, 330.0, 95.0, 3.0, 0.0)
	} else {
		basic_font[16 + atk].Draw(screen, 330.0, 95.0, 3.0, 0.0)
	}

	if def >= 10 {
		basic_font[16 + 1].Draw(screen, 470.0, 95.0, 3.0, 0.0)
		basic_font[16].Draw(screen, 494.0, 95.0, 3.0, 0.0)
	} else {
		basic_font[16 + def].Draw(screen, 470.0, 95.0, 3.0, 0.0)
	}

	aport := 2
	dport := 2
	if atk > def {
		aport = 3
		dport = 1
	} else if def > atk {
		aport = 1
		dport = 3
	}
	if def == 0 {
		aport = 4
		dport = 0
	} else if atk == 0 {
		aport = 0
		dport = 4
	}

	ui_anims["port_" + strconv.Itoa(g.defender_army) + "_" + strconv.Itoa(dport)].Draw(g.portflip_image, 88.0, 161.0, 1.0, g.time)
    op.GeoM.Reset()
	op.GeoM.Scale(-1.0, 1.0)
	op.GeoM.Translate(800, 0)
	screen.DrawImage(g.portflip_image, op)

	ui_anims["port_" + strconv.Itoa(g.attacker_army) + "_" + strconv.Itoa(aport)].Draw(screen, 88.0, 161.0, 1.0, g.time)
}

func (g BattleGame) DrawAttacker(xoff, yoff, toff, ktime float64) {
	aas := strconv.Itoa(g.attacker_army)

	if g.time <= 0.5 {
		// running in
		battle_anims[aas + "_" + g.attacker_type + "_run"].Draw(g.attacker_image, xoff, yoff, 1.0, g.time)
	} else if g.time > ktime {
		battle_anims[aas + "_" + g.attacker_type + "_die"].Draw(g.attacker_image, xoff, yoff, 1.0, g.time - ktime)
	} else if g.time - toff > 0.5 && g.time - toff < 1.5 {
		//attack
		battle_anims[aas + "_" + g.attacker_type + "_shoot"].Draw(g.attacker_image, xoff, yoff, 1.0, g.time - toff - 0.5)
	} else {
		// idle
		battle_anims[aas + "_" + g.attacker_type + "_idle"].Draw(g.attacker_image, xoff, yoff, 1.0, g.time)
	}
}

func (g BattleGame) DrawDefender(xoff, yoff, toff, ktime float64) {
	das := strconv.Itoa(g.defender_army)

	if g.time > ktime {
		battle_anims[das + "_" + g.defender_type + "_die"].Draw(g.defender_image, xoff, yoff, 1.0, g.time - ktime)
	} else if g.time - toff > 2.0 && g.time - toff < 3.0 {
		//attack
		battle_anims[das + "_" + g.defender_type + "_shoot"].Draw(g.defender_image, xoff, yoff, 1.0, g.time - toff - 2.0)
	} else {
		// idle
		battle_anims[das + "_" + g.defender_type + "_idle"].Draw(g.defender_image, xoff, yoff, 1.0, g.time)
	}
}

func (g BattleGame) GetCurrentStrs() (int, int) {
	atk := g.start_astr
	def := g.start_dstr

	if g.time > 2.0 {
		if g.time > 2.0 + 1.5 {
			atk = g.end_astr
		} else {
			atktime := (2.0 - g.time) / -1.5
			atk = g.start_astr + ((g.end_astr - g.start_astr) * atktime)
		}
	}

	if g.time > 0.5 {
		if g.time > 0.5 + 1.5 {
			def = g.end_dstr
		} else {
			deftime := (0.5 - g.time) / -1.5
			def = g.start_dstr + ((g.end_dstr - g.start_dstr) * deftime)
		}
	}

	iatk := int(math.Floor(atk * 10.0))
	if iatk <= 0 && atk > 0  {
		iatk = 1
	}

	idef := int(math.Floor(def * 10.0))
	if idef <= 0 && def > 0 {
		idef = 1
	}

	return iatk, idef
}

func (g *BattleGame) Layout(outsideWidth, outsideHeight int) (int, int) {
    return outsideWidth, outsideHeight
}

func (g BattleGame) GetResult() (bool, bool) {
	return g.time > 4.0, false
}
