package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/Yarnsh/hippo/input"
    "github.com/Yarnsh/hippo/animation"
    "github.com/Yarnsh/hippo/audio"
    "log"
    "fmt"
    "os"
    "embed"
)

const (
    screenWidth  = 800
    screenHeight = 600

    tacMapWidth = 40
    tacMapHeight = 30
)

var (
    //go:embed assets/*
    EmbeddedFileSystem embed.FS

    ui_anims map[string]animation.Animation
    basic_font []animation.Animation
    lumi_font []animation.Animation
    jelly_font []animation.Animation
    dizzy_font []animation.Animation
    ember_font []animation.Animation
    battle_anims map[string]animation.Animation
)

type MetaGame struct {
    camp CampaignGame
}

func NewMetaGame() MetaGame {
    result := MetaGame{}

    player_control := input.NewInputActionHandler()
    player_control.RegisterKeyboardAction("left", ebiten.KeyA)
    player_control.RegisterKeyboardAction("right", ebiten.KeyD)
    player_control.RegisterKeyboardAction("up", ebiten.KeyW)
    player_control.RegisterKeyboardAction("down", ebiten.KeyS)
    player_control.RegisterKeyboardAction("accept", ebiten.KeyN)
    player_control.RegisterKeyboardAction("cancel", ebiten.KeyM)

    player_control.RegisterKeyboardAction("cheat0", ebiten.KeyF1)
    player_control.RegisterKeyboardAction("cheat1", ebiten.KeyF2)
    player_control.RegisterKeyboardAction("cheat2", ebiten.KeyF3)

    result.camp = NewCampaign(GAME_STAGES, player_control)

    return result
}
func (g *MetaGame) Update() error {
    err := g.camp.Update()
    finished, victory := g.camp.GetResult()
    if finished { // test until we get VN stuff going
        fmt.Println("Game Over!", victory)
        os.Exit(0)
    }
    return err
}
func (g *MetaGame) Draw(screen *ebiten.Image) {
    g.camp.Draw(screen)
}
func (g *MetaGame) Layout(outsideWidth, outsideHeight int) (int, int) {
    return outsideWidth, outsideHeight
}

func main() {
    animation.FileSystem = EmbeddedFileSystem
    audio.FileSystem = EmbeddedFileSystem
    
    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle("Invader Wars")

    ui_anims, _ = animation.LoadAnimationMap("assets/ui.json")
    var err error
    battle_anims, err = animation.LoadAnimationMap("assets/battle.json")
    if err != nil {
        panic(err)
    }
    basic_font, _ = animation.NewFontAnimationMap("assets/font.png", 8, 15, 32, 23)
    lumi_font, _ = animation.NewFontAnimationMap("assets/lumifont.png", 8, 15, 32, 23)
    jelly_font, _ = animation.NewFontAnimationMap("assets/jellyfont.png", 8, 15, 32, 23)
    dizzy_font, _ = animation.NewFontAnimationMap("assets/dizzyfont.png", 8, 15, 32, 23)
    ember_font, _ = animation.NewFontAnimationMap("assets/emberfont.png", 8, 15, 32, 23)
    InitTacMapData()
    InitUnitData()

    audio.Init()
    audio.LoadSoundPath("assets/sfx/marching.wav")
    audio.LoadSoundPath("assets/sfx/engine.wav")
    audio.LoadSoundPath("assets/sfx/car.wav")
    audio.LoadSoundPath("assets/sfx/rapidshoot.wav")
    audio.LoadSoundPath("assets/sfx/boom.wav")
    audio.LoadSoundPath("assets/sfx/death.wav")
    audio.LoadSoundPath("assets/sfx/build.wav")
    audio.LoadSoundPath("assets/sfx/accept.wav")
    audio.LoadSoundPath("assets/sfx/cancel.wav")
    audio.LoadSoundPath("assets/sfx/click.wav")

    game := NewMetaGame()
    if err := ebiten.RunGame(&game); err != nil {
        log.Fatal(err)
    }
}
