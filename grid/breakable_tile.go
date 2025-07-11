package grid

import (
	"main/utils"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type BreakableTile struct {
	Tile   *int
	Signal string
}

type BreakableTileJson struct {
	Pos    utils.Vec2
	Signal int
}

var BreakableTIleImg, _, _ = ebitenutil.NewImageFromFile("./art/breakable_tile.png")
var SelectedBreakableTile *BreakableTile

func (tile *BreakableTile) Serialize(chunk_x, chunk_y, tile_x, tile_y int) BreakableTileJson {
	new_breakable_tile := BreakableTileJson{}

	signal, err := strconv.Atoi(tile.Signal)
	if err != nil {
		panic(err)
	}

	new_breakable_tile = BreakableTileJson{utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}, signal}

	return new_breakable_tile
}

func (breakable_tile *BreakableTileJson) Deserialize(tile *int) BreakableTile {
	new_tile := BreakableTile{}

	new_tile.Signal = strconv.Itoa(breakable_tile.Signal)
	new_tile.Tile = tile

	return new_tile
}

func (level *Level) ManageBreakableTiles() {
	for i, breakable_tile := range level.BreakableTile {
		if *breakable_tile.Tile != -4 {
			utils.RemoveArrayElement(i, &level.BreakableTile)
			*breakable_tile.Tile = 0
		}
	}
}

func (level *Level) PlaceBreakableTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -4

					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					for i := range level.BreakableTile {
						t := &level.BreakableTile[i]
						if t.Tile == tile {
							can_add = false
							SelectedBreakableTile = t
						}
					}

					if can_add {
						if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -4 {
							level.BreakableTile = append(level.BreakableTile, BreakableTile{&level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)], "0"})
							SelectedBreakableTile = &level.BreakableTile[len(level.BreakableTile)-1]
							SelectedEnemySpawner = nil
						}
					}

					level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
				}
			}
		}
	}
}

func (level *Level) SelectBreakableTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					if *tile == -4 {
						for i := range level.BreakableTile {
							t := &level.BreakableTile[i].Tile
							if *t == tile {
								can_add = false
								SelectedBreakableTile = &level.BreakableTile[i]
								SelectedEnemySpawner = nil
								SelectedTriggerTile = nil
							}
						}

						if can_add {
							if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -4 {
								level.BreakableTile = append(level.BreakableTile, BreakableTile{tile, "1"})
								SelectedBreakableTile = &level.BreakableTile[len(level.BreakableTile)-1]
							}
						}
					}
				}
			}
		}
	}
}
