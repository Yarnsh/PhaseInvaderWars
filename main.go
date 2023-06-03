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

    player_control input.InputActionHandler
    ui_anims map[string]animation.Animation
)

type MetaGame struct {
    // test until we get VN stuff going
    tac_game TacticalGame
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

    result.tac_game = NewTacticalGame(tac_map_1, player_control)

    result.tac_game.AddUnit("infantry", 7, 13, 0)
    result.tac_game.AddUnit("tank", 7, 16, 0)
    result.tac_game.AddUnit("antitank", 8, 15, 0)

    result.tac_game.AddUnit("infantry", 23, 13, 1)
    result.tac_game.AddUnit("tank", 25, 16, 2)
    result.tac_game.AddUnit("antitank", 28, 15, 3)

    return result
}
func (g *MetaGame) Update() error {
    err := g.tac_game.Update()
    finished, _ := g.tac_game.GetResult()
    if finished { // test until we get VN stuff going
        fmt.Println("Game Over!")
        os.Exit(0)
    }
    return err
}
func (g *MetaGame) Draw(screen *ebiten.Image) {
    g.tac_game.Draw(screen)
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
    InitTacMapData()
    InitUnitData()

    game := NewMetaGame()
    if err := ebiten.RunGame(&game); err != nil {
        log.Fatal(err)
    }
}
