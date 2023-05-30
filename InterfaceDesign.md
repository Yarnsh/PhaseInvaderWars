# Interface Design
Designing the structure and methods I will hopefully need before the jam to speed things up.

## Global

### Methods
`CalculateCombatResult(attacker, defender Unit) (float64, float64)`
- Returns the strength of the attacker and defender after calculating how combat would go between them

## TacticalGame

### Methods
`NewTacticalGame(playerArmy Army, enemyArmy Army, enemyAI TacticalAI, map TacticalMap) TacticalGame`
- Constructor

`Update() err`
- Called once per frame
- All the logic of playing the tactical battle map
- Will call Update() on CombatGame if one is ongoing

`Draw(target *ebiten.Image)`
- Called once per frame
- Draw everything relating to the tactical battle to `target`, map, units, menus, etc.

`Layout(outsideWidth, outsideHeight int) (int, int)`
- Method to match the ebitengine game interface, no real use

`GetResult() (bool, bool)`
- Returns:
  - True if battle is over
  - True if player won, False if AI did

`GetMap() TacticalMap`
- Returns the map, unmodified

`GetMoveMap() TacticalMap`
- Returns:
  - The TacticalMap of this game, but any tile with a unit on it is marked as impassable

## TacticalMap

### Methods
`Draw(target *animation.DrawBuffer)`
- Draw all the tiles to the DrawBuffer (need to update the engine to have one of those)

## Tile

### Methods
`Draw(target *animation.DrawBuffer)`
- Draw this one tile into the DrawBuffer

`GetDefense() float64`
- Returns damage recieved multiplier for this tile type

`DrawBackground(target *ebiten.Image)`
- Draw the tile's combat game background to `target`
- Always drawn as if we are the attacker, caller is responsible for flipping and positioning it for defenders

`GetMoveCost() int`
- Returns how much it costs to move into this tile

`IsP1Factory() bool`
- Returns true is a factory owned by the human

`IsP2Factory() bool`
- Returns true is a factory owned by the AI

`IsP1HQ() bool`
- Returns true is the HQ owned by the human

`IsP2HQ() bool`
- Returns true is the HQ owned by the AI

## TacticalAI

### Methods
`BestMoveForUnit(game TacticalGame, unit Unit) (int, int, int, int)`
- Calculate where the AI would want to move this unit, and where it should attack
- Returns:
  - X coord of position to move to
  - Y coord of position to move to
  - X coord of position to attack (if same as move, don't attack)
  - Y coord of position to attack (if same as move, don't attack)

`evaluateGameState(game TacticalGame) float64`
- Calculate how good we think the position is for the AI player
- Since AI is only ever player 2 we don't need to pass what side is checking or anything like that
- Returns:
  - Rough numerical representation of how good the position is for player 2 (more is better)

`GetNextBuildUnit() int`
- Returns the next type of unit we would like to build
- Don't update the internal counter, since we may not build it depending on cost

`UpdateNextBuildUnit()`
- Update the internal counter used for `GetNextBuildUnit()`
- Intended to be called after we build a unit

## Army

### Methods
`DrawCO(target *ebiten.Image, mood int)`
- Draw the army's CO portrait with given mood to `target`
- Always drawn as if we are on the left, caller is responsible for flipping and positioning

## Unit

### Methods
`NewUnit(type, x, y int) Unit`
- Constructor
- `x` and `y` are the tile coordinates to make this unit at

`Draw(target *animation.DrawBuffer, time float64)`
- Draw this one unit into the DrawBuffer
- This should also draw the strength number if strength < 10

`DrawMoveRight(target *animation.DrawBuffer, startTime float64, time float64)`
- Draw this one unit moving into the DrawBuffer
- This should also draw the strength number if strength < 10
- Use startTime and time to figure out offset to draw at for illusion of motion

`DrawMoveLeft(target *animation.DrawBuffer, startTime float64, time float64)`
- Draw this one unit moving into the DrawBuffer
- This should also draw the strength number if strength < 10
- Use startTime and time to figure out offset to draw at for illusion of motion

`DrawMoveUp(target *animation.DrawBuffer, startTime float64, time float64)`
- Draw this one unit moving into the DrawBuffer
- This should also draw the strength number if strength < 10
- Use startTime and time to figure out offset to draw at for illusion of motion

`DrawMoveDown(target *animation.DrawBuffer, startTime float64, time float64)`
- Draw this one unit moving into the DrawBuffer
- This should also draw the strength number if strength < 10
- Use startTime and time to figure out offset to draw at for illusion of motion

`DrawCombatEnter(target *ebiten.Image, time float64, x, y float64)`
- Draw the combat game running in animation to the target

`DrawCombatIdle(target *ebiten.Image, time float64, x, y float64)`
- Draw the combat game idle animation to the target
- Draw as if we are on the left, caller is responsible for flipping and positioning

`DrawCombatShoot(target *ebiten.Image, time float64, x, y float64)`
- Draw the combat game shoot animation to the target
- Draw as if we are on the left, caller is responsible for flipping and positioning

`DrawCombatDie(target *ebiten.Image, startTime float64, time float64, x, y float64)`
- Draw the combat game die animation to the target
- Draw as if we are on the left, caller is responsible for flipping and positioning

`GetPossibleMoves(game TacticalGame) []utils.IntPair`
- Figure out where the unit can reach on the map
- Keep in mind we can't move through other units
- Returns:
  - List of X Y coordinates we could reach this turn, including the start position

`GetStrength() float64`
- Return the unit's current combat strength
- Used for combat calculations

`SetStrength(strength float64)`
- Sets the strength of the unit
- Used with the results of combat usually

## CombatGame

### Methods
`NewCombatGame(attacker Unit, defender Unit, attackerArmy Army, defenderArmy Army, attackerTile Tile, defenderTile Tile) CombatGame`
- Constructor
- Calculate the after combat unit strengths here

`Update() err`
- Called once per frame
- Mostly to move forward the clock, all the real logic will be in Draw

`Draw(target *ebiten.Image)`
- Called once per frame
- Draw UI, CO portraits, units, and their various animations as this plays out

`Layout(outsideWidth, outsideHeight int) (int, int)`
- Method to match the ebitengine game interface, no real use

`GetResult() (bool, bool)`
- Returns:
  - True if animation is done
  - Always false, just here to match the game interface

## VnGame

### Methods
`NewVnGame(stages []Stage) VnGame`
- Constructor

`Update() err`
- Called once per frame
- Handle time handling logic VN text boxes and such
- Handle starting and calling Update() for tactical battles if the stage has one
- Check for completion of tactical battles to change stages and such

`Draw(target *ebiten.Image)`
- Called once per frame
- Call Draw() for tactical game if one is happening
- Otherwise do all the text and UI drawing of the current VN stage

`Layout(outsideWidth, outsideHeight int) (int, int)`
- Method to match the ebitengine game interface, no real use

`GetResult() (bool, bool)`
- Returns:
  - True if we are done with the last stage
  - Always false, just here to match the game interface

## Stage

### Methods
`IsTacticalBattle() bool`
- Returns true if this stage should start a tactical game

`CreateTacticalBattle() TacticalGame`
- Returns a TacticalGame with the settings this stage should run
- If `IsTacticalBattle()` returns false for this stage, this will return a default TacticalGame object

`GetTextScrollTime() float64`
- Returns the amount of time it should take from starting this stage to display all the text in the text box
- Used to know when user input should go to next step vs skipping text

`Draw(target *ebiten.Image, time float64)`
- Draws the VN stuff of this stage, including portraits and text
- `time` is time in seconds since the start of the stage
