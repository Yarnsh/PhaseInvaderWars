package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/Yarnsh/hippo/input"
    "github.com/Yarnsh/hippo/audio"
    "strconv"
)

const (
	CHARS_PER_LINE = 40
	TEXT_BOX_WIDTH = 700.0
	TEXT_BOX_HEIGHT = 130.0
	TEXT_BOX_Y_OFFSET = -30.0
)

type CampaignGame struct {
	player input.InputActionHandler
	time float64
	stage_start_time float64

	stages []Stage
	current_stage int
	current_battle TacticalGame
	flip_img *ebiten.Image
	text_drawn int
}

func NewCampaign(stages []Stage, player input.InputActionHandler) CampaignGame {
	return CampaignGame{
		player: player,
		flip_img: ebiten.NewImage(800, 600),
		stages: stages,
	}
}

func (g *CampaignGame) Update() error {
	g.time += 1.0/60.0

	if g.player.IsActionJustPressed("cheat0") {
		g.JumpToStage(0)
	}
	if g.player.IsActionJustPressed("cheat1") {
		g.JumpToStage(5)
	}
	if g.player.IsActionJustPressed("cheat2") {
		g.JumpToStage(8)
	}

	if g.stages[g.current_stage].is_battle {
		g.current_battle.Update()
		done, win := g.current_battle.GetResult()
		if done && win {
			g.NextStage()
		} else if done {
			g.JumpToStage(g.stages[g.current_stage].lose_stage)
		}
	} else {
		// handle the VN stuff
		stage_done_time := float64(len(g.stages[g.current_stage].text)) * g.stages[g.current_stage].text_speed
		if stage_done_time < g.time - g.stage_start_time {
			// we are on a text box, so wait for input
			if g.player.IsActionJustPressed("accept") {
				audio.Play("assets/sfx/accept.wav")
				g.NextStage()
			}
		} else {
			// animation ongoing with a text box, on input we should jump forward in time to the end of the text draw time
			if g.player.IsActionJustPressed("accept") {
				audio.Play("assets/sfx/accept.wav")
				g.time = g.stage_start_time + stage_done_time
			}
		}
	}

    return nil
}
func (g *CampaignGame) Draw(screen *ebiten.Image) {
	if g.stages[g.current_stage].is_battle {
		g.current_battle.Draw(screen)
	} else {
		ui_anims["vnbg"].Draw(screen, 400, 600, 1.0, 0.0)
		ui_anims["port_" + strconv.Itoa(g.stages[g.current_stage].p1) + "_" + strconv.Itoa(g.stages[g.current_stage].e1)].Draw(g.flip_img, 628, 400, 2.0, 0.0)
		ui_anims["port_" + strconv.Itoa(g.stages[g.current_stage].p2) + "_" + strconv.Itoa(g.stages[g.current_stage].e2)].Draw(screen, 628, 400, 2.0, 0.0)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(-1.0, 1.0)
		op.GeoM.Translate(800, 0)
		screen.DrawImage(g.flip_img, op)

		font := basic_font
		if g.stages[g.current_stage].font == 1 {
			font = lumi_font
		}
		if g.stages[g.current_stage].font == 2 {
			font = jelly_font
		}
		if g.stages[g.current_stage].font == 3 {
			font = dizzy_font
		}
		if g.stages[g.current_stage].font == 4 {
			font = ember_font
		}

		line := 1.0
		column := 0
		chars_to_show := charactersToShow(g.stages[g.current_stage].text, g.stages[g.current_stage].text_speed, g.time - g.stage_start_time)
		for _, char := range []rune(g.stages[g.current_stage].text[:chars_to_show]) {
			if column > CHARS_PER_LINE {
				column = 0
				line += 1.0
			}
			// For some reason we are a row off, hence "- 32", might be a hippo issue
			font[char - 32].Draw(
				screen,
				float64((column * 16) + 8 + 20) + (800 / 2.0) - (TEXT_BOX_WIDTH / 2.0),
				600 + TEXT_BOX_Y_OFFSET + (-TEXT_BOX_HEIGHT) + 20.0 + (30.0 * line),
				2.0,
				0.0) // not pretty to hard code all this but whatever
			column += 1
		}

		if chars_to_show > g.text_drawn {
			g.text_drawn = chars_to_show
			audio.Play("assets/sfx/click.wav")
		}
	}
}
func (g *CampaignGame) Layout(outsideWidth, outsideHeight int) (int, int) {
    return outsideWidth, outsideHeight
}

func (g CampaignGame) GetResult() (bool, bool) {
	return false, false
}

func (g *CampaignGame) NextStage() {
	if g.current_stage + 1 >= len(g.stages) {
		return
	}
	g.current_stage += 1
	if g.stages[g.current_stage].is_battle {
		g.current_battle = CreateTacticalGame(g.stages[g.current_stage].level, g.player)
	}
	g.text_drawn = 0
	g.stage_start_time = g.time
}

func (g *CampaignGame) JumpToStage(stage int) {
	g.current_stage = stage
	if g.stages[g.current_stage].is_battle {
		g.current_battle = CreateTacticalGame(g.stages[g.current_stage].level, g.player)
	}
	g.text_drawn = 0
	g.stage_start_time = g.time
}

func charactersToShow(text string, text_speed, time float64) int {
	chars := len(text)
	totaltime := float64(chars) * text_speed
	if totaltime <= time {
		return chars
	}
	return int((time / totaltime) * float64(chars))
}

type Stage struct {
	text string
	text_speed float64
	is_battle bool
	level int
	p1, e1, p2, e2 int
	lose_stage int
	font int
}
