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
)

type MetaGame struct {
    // test until we get VN stuff going
    tac_game TacticalGame
}

func NewMetaGame() MetaGame {
    result := MetaGame{}

    player_control := input.NewInputActionHandler()
    player_control.RegisterKeyboardAction("Left", ebiten.KeyLeft)
    player_control.RegisterKeyboardAction("Right", ebiten.KeyRight)
    player_control.RegisterKeyboardAction("Rotate", ebiten.KeyUp)
    player_control.RegisterKeyboardAction("Fall", ebiten.KeyDown)
    player_control.RegisterKeyboardAction("Accept", ebiten.KeyZ)
    player_control.RegisterKeyboardAction("Cancel", ebiten.KeyX)

    result.tac_game = NewTacticalGame(tac_map_1)

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

    InitTacMapData()

    game := NewMetaGame()
    if err := ebiten.RunGame(&game); err != nil {
        log.Fatal(err)
    }
}
