package grid

import (
	"main/utils"
	"strconv"

	"github.com/bob4321at/textures"
)

type WaterTile struct {
	Tile          *int
	Top_Bottom    bool
	ReceiveSignal string
}

type WaterTileJson struct {
	Pos           utils.Vec2
	Top_Bottom    bool
	ReceiveSignal int
}

var Water_Tile_Imgs = map[bool]textures.Texture{
	false: *textures.NewTexture("./art/water_tile.png", ""),
	true:  *textures.NewTexture("./art/water_tile_bottom.png", ""),
}

var SelectedWatertile *WaterTile

func (tile *WaterTile) Serialize(chunk_x, chunk_y, tile_x, tile_y int) WaterTileJson {
	new_water_tile := WaterTileJson{}

	signal, err := strconv.Atoi(tile.ReceiveSignal)
	if err != nil {
		panic(err)
	}

	new_water_tile = WaterTileJson{utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}, tile.Top_Bottom, signal}

	return new_water_tile
}

func (water_tile *WaterTileJson) Deserialize(tile *int) WaterTile {
	new_tile := WaterTile{}

	new_tile.Top_Bottom = water_tile.Top_Bottom
	new_tile.Tile = tile
	new_tile.ReceiveSignal = strconv.Itoa(water_tile.ReceiveSignal)

	return new_tile
}

func (level *Level) ManageWaterTiles() {
	for i, water_tile := range level.WaterTiles {
		if *water_tile.Tile != -12 {
			utils.RemoveArrayElement(i, &level.WaterTiles)
			*water_tile.Tile = 0
		}
	}
}

func (level *Level) PlaceWaterTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -12

					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					for i := range level.WaterTiles {
						t := &level.WaterTiles[i]
						if t.Tile == tile {
							can_add = false
						}
					}

					if can_add {
						if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -12 {
							level.WaterTiles = append(level.WaterTiles, WaterTile{&level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)], false, "0"})
						}
					}

					level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
				}
			}
		}
	}
}
func (level *Level) SelectWaterTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					if *tile == -12 {
						for i := range level.WaterTiles {
							t := &level.WaterTiles[i].Tile
							if *t == tile {
								can_add = false
								SelectedWatertile = &level.WaterTiles[i]
							}
						}

						if can_add {
							if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -12 {
								level.WaterTiles = append(level.WaterTiles, WaterTile{tile, false, "0"})
								SelectedWatertile = &level.WaterTiles[len(level.WaterTiles)-1]
							}
						}
					}
				}
			}
		}
	}
}
