## Quick Reference of Tasks
- Modify SFX and music tracks to fit
- Create sprites
- Create animation files for everything
- Code everything
- Write terrible dialogue

# Pre-jam design
Overall idea is an Advance Wars lite, Lumi as player CO, new invaders as enemy COs.
Therefore 3 levels, optimally with short VN like segments between them for a bit of character.

There will be many things in this document that are chosen for speed and ease of creation rather than quality, as this is a 3 day game jam and I am only one man.

## General Details
Screen size will be 800x600, unchangeable.
Saving, menus, options of any kind are all stretch goals. UI takes way too much time for how much it adds.

## Tactical Battle Map
- Top Priority
This is the "game" part of this game so this is top priority to create. Everything else is flavor in comparison.

### Map Itself
A grid of tiles that matches the size of the screen (to avoid coding screen scroll)
4 tile types come to mind, Plains, Forest, Mountain, Water.

Tile size should be 20x20 (ignoring overlap bits), 40x30 map size given the screen size.
Hippo engine should be updated to draw with a given depth, so we can layer the overlap of tiles and units correctly.

Plains:
- Basic tile with nothing special, basically a green square

Forest:
- 20% defense bonus in combat for unit standing on it
- Double move cost through it (for everything except infantry)
- Tile should overlap with tile above it slightly to give sense of tree height

Mountain:
- 50% defense bonus in combat for unit standing on it
- Triple move cost through it (for everything except infantry)
- Double move cost for infantry

Water:
- No units should be able to move through this
- (Stretch goal) Pleasant wave animation

Map should also have an HQ tile and a factory tile.

(Stretch goal) Capturable city tile, that provides places to heal and increases income.

HQ:
- 50% defense and default move speed
- As in Advance Wars, capturing this means victory
- (Stretch goal) capture mechanic similar to Advance Wars that uses unit health as how much you capture the building per turn, out of a total of 20
- If no capture mechanic, instant capture when you step on it would be fine.

Factory:
- 30% defense and default move speed
- Place to spawn new units
- (Stretch goal) if we have capturable cities, making the factories capturable would be good too, otherwise having to bring units to occupy them is more interesting.

(Stretch goal) City:
- 30% defense and default move speed
- Capturable
- If you own it, it will increase your funds per turn by $1000

All tiles that can be fought on should have a background for the combat animation, 400x600 to take up half the screen.

### Units
The pawns whos lives you will be throwing away.

3 basic units per army, Infantry, Tank, Anti-tank. Simple rock-paper-scissors dynamic. The 3 units for each army should be visually different, but mechanically the same.

Infantry:
- Damage dealing
  - 55% vs Infantry
  - 10% vs Tank
  - 85% vs Anti-Tank
- 3 movement
- $1000

Tank:
- Damage dealing
  - 80% vs Infantry
  - 55% vs Tank
  - 40% vs Anti-Tank
- 6 movement
- $7000

Anti-tank:
- Damage dealing
  - 15% vs Infantry
  - 90% vs Tank
  - 30% vs Anti-Tank
- 3 movement
- $3000

Health of all units will be stored as a float, max 10. This value rounded up to the nearest int will be the displayed unit strength.

(Stretch goal) A special unit for each army

Lumi:
- Jet plane, expensive and ignores movement costs of terrain, can go over water
Jelly:
- Very tanky jellyfish alien thing, moderate damage output to all, slow
Dizzy:
- No god damn idea
Ember:
- Flamer thrower infantry, cheap-ish and very very strong against infantry and anti-tank, useless against tanks

### Interactivity

Mode 1: Cursor

- Have a cursor that you can move around to select different tiles on the map.
- `select` and `cancel` actions, basically A and B on controller or Z and X on keyboard probably.
- Arrows or d-pad to move the cursor around.
- Cursor should have a slightly different sprite when it is over something selectable (owned unit or factory)
- Cursor movement should be instant as soon as you press a direction
- Holding direction should move the cursor every 0.2 seconds
If selecting an empty tile, or one with an enemy unit/building on it, switch to mode 3 (menu), options:
  - End Turn
  - Cancel
(Stretch goal) If selecting an enemy unit:
- Switch to an undesigned mode, where the move+1 range of the unit is displayed as red overlay, to show attack range
If selecting tile with a player controlled unit that has more than 0 `action_points`:
- Switch to mode 2 (unit control)
If selecting tile with player controlled factory:
- Switch to mode 3 (menu), options:
  - One option per unit type, showing unit name and cost, gray out if not enough funds
  - (stretch goal) special unit, with its cost, grayed out if cant afford
  - Cancel

---

Mode 2: Unit Control

A single player controlled unit will be selected when in this mode. Cursor will still be visible in this mode.

The `cancel` action should bring you back to Mode 1

Each unit will have an `action_points` value of 2 at the beginning of a turn.

If `action_points` == 2:
- Draw a blue overlay on every tile reachable by this unit this turn, calculated from its movement points against the terrain around us
- Terrain needs a dikstras like algorythm to find reachable tiles and lowest cost path to them, basically running pathfinding to nothing from the unit
- The `select` action on any blue tile will cause the unit to move there and reduce `action_points` to 1
- There should be an animation of the unit moving to the new location, during which the player will have no control
- A red overlay should be drawn on any non-blue tile adjacent to a blue tile, to show theoretical attack range
- If `action` is pressed on the unit again, move to the `action_points` == 1 case without actually changing the value
If `action_points` == 1:
- Draw a red overlay on tiles adjacent to selected unit
- If `select` is pressed on a red tile that has an enemy unit on it, initiate combat with that unit and reduce `action_points` to zero
If `action_points` == 0:
- Return to Mode 1

---

Mode 3: Menu

A simple rectangular menu with a vertical list of options should appear centered on the screen. There should be a cursor showing the currently hovered item.
Pressing `select` will do a custom action depending on the selected option, pressing `cancel` should bring you back to Mode 1.

Possible options:
- "Cancel", should bring you back to Mode 1
- "End Turn", end the player turn and let the AI do theirs
- Unit listing, if greyed out should do nothing, otherwise spawn a new unit with 0 `action_points` of the selected type on the selected tile, reduce funds by unit cost.

### Enemy AI
This section is very TODO. General idea should be to have an evaluation function, and test the map state after every possible move of a single unit,
 then do the best rated move.
Repeat this for every unit the AI controls.

Next, for every unoccupied factory the AI controls look at the next entry in a predefined "build order" for this level. If there are enough funds spawn that unit and
 increment the build order idx. The build order should also have a predefined "loop idx", if the idx goes past the end we loop to that point and not the start.

General ideas for evaluation function:
- Basic material advantage, fund cost of your active units vs theirs
- Distance to enemy HQ (closer is better)
- Distance of enemy to our HQ (closer is worse)
- Friendly units on a more defensive tile will improve evaluation
- Friends being within attack range of an enemy unit is bad (scaled on max amount of damage they could do some unit in range)

Testing this for effectiveness would take forever, so as long as it looks even remotely non-stupid, accept it.

During an enemy turn movement and battle animations should play, but the cursor and tile overlays and such have no reason to be shown.

### Asset Summary
Sprites
- 4 sprites for map tiles
- 8 sprites for HQs and factories
  - 1 HQ sprite per army
  - 1 factory sprite per army
- 160 unit sprites total
  - 10 sprites per unit per army
    - 2 x 4 for 2 frame walk cycle in each direction
    - 2 for 2 frame idle cycle
- 2 sprites for cursor states
- 1 sprite for menu cursor
- 1 font sheet for text

---

SFX
- "Marching" sfx for infantry movement
- "Engine" sfx for tank movement
- "Rolling" or "Car" like sfx for anti-tank movement
- "Select" bloop sfx
- "Cancel" bloop sfx
- "Move Cursor" light click sfx
- "Build Unit" mechanical rustling sfx

---

Music
- 3 battle themes, one per map
- Victory theme (short)
- Losing theme (short)

---
Totals

- 176 sprites
- 7 sfx
- 5 music

## Combat Animation
- Second priority, the prettiest part and gives the barebones combat a bit of flavor.

At start takes unit types, unit strengths, and tile types of both sides to render the battle and also to calculate the effects.
Effect calculation is needed to render properly, but for the tactical map it should be calculated on its own, as it is useful to check
 without an animation for AI stuff. No user input does anything during these which simplifies things.

Attacking unit rendered on the left side, defending unit on the right.
A half screen static background of the attacker tile on the left, and one of the defender tile on the right.
Top corners will have CO portraits, start in neutral expression, after damage is done switch to sad for whoever is lower strength, and happy to the other.
If a unit is destroyed outright, set the survivor to very happy and loser to very sad. 1 instance of a unit sprite drawn per 2 strength (rounding up to a multiple of 2).
Start with attacker unit sprites running in from off screen, defender units idling where they are.
Attacker will play shooting animation first, defender second (unless they are destroyed already).
Play bullet impact animation on other side while firing.

When strength goes down due to combat play death animation for a unit if the amount of units on screen should be reduced.

---

Damage calculation for a single side:

`(Base Damage vs Given Unit) * (Unit Strength) * (Defender Tile Defense)`

First run this calculation for the attacker and reduce defender strength by the result.
Secondly run this for the defender with their new strength (if they are alive) and reduce attacker strength by the result.

### Asset Summary
Sprites
- 16 idle frames
  - 1 frame idle animation per army per unit
- 32 shooting frames
  - 2 frame firing animation per army per unit
- 16 death frames
  - 1 frame death animation per army per unit
- 32 running frames
  - 2 frame running animation per army per unit
- 32 shot impact frames
  - 2 frame impact animation per army per unit
  - Can likely reuse some of these between units
- 20 CO portraits
  - 5 expressions per CO
- 6 backgrounds
  - 1 for each type of tile
- 1 UI border to hold portrait and such
- Reuse font image from tactical map

---
SFX
- Reuse "Marching" "Engine" and "Car" or whatever from tactical map for units walking in
- 16 "Shooting" sfx
  - 1 per unit per army
  - Can probably reuse between units somewhat
  - No impact sfx, it will be thought to be included in the shoot noise
- 16 "Death" sfx
  - 1 per unit per army
  - Can probably reuse between units somewhat

---
Music
- Let whatever track is playing in the tactical map continue playing

---
Totals
- 155 sprites
- 32 sfx
- 0 music

## VN Bits
- Least priority, drop if low on time

Small talking segments between missions. Probably just lumi talking to enemy CO of the next mission.

Text box at the bottom to contain the speech, with a small label spot at it's top left to hold name of current speaker.
CO portrait holders on left and right for both participants. Reuse the CO portraits from the combat animation bit here.

The things they say don't really need to make any sense at all, I don't know the invaders very well outside Lumi anyway. Just grab random things from their Twitter probably.

Ideas for topics of discussion:
- Jelly would just make "awawawawawa" noises and Lumi would act baby crazy, then they fight for some reason.
- Ember would say she burned down Lumi's second kitchen, this would enrage Lumi.
- Dizzy would act disinterested and deny Lumi's advances, creating a self defense situation.

### Asset Summary
Sprites
- 1 Full screen UI, including portrait holders and text box all in one
- Reuse font image from tactical map
- Reuse CO portrats from combat animations

---
SFX
- At most, a light noise of some kind to play as the letters get scrolled in, optional.

---
Music
- 1 BGM track sould suffice for all the conversations.

---
Writing
- 4-8 text boxes worth of text per opponent

---
Totals
- 1 sprite
- 1 sfx
- 1 music track
- 24 text boxes max
