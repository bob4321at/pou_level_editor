package grid

import (
	"main/utils"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type FloodTile struct {
	Tile          *int
	Speed         string
	SendSignal    string
	ReceiveSignal string
}

type FloodTileJson struct {
	Pos           utils.Vec2
	Speed         int
	SendSignal    int
	ReceiveSignal int
}

var FloodTileImg, _, _ = ebitenutil.NewImageFromFile("./art/flood_tile.png")

var SelectedFloodtile *FloodTile

func (tile *FloodTile) Serialize(chunk_x, chunk_y, tile_x, tile_y int) FloodTileJson {
	new_water_tile := FloodTileJson{}

	speed, err := strconv.Atoi(tile.Speed)
	if err != nil {
		panic(err)
	}

	receive_signal, err := strconv.Atoi(tile.ReceiveSignal)
	if err != nil {
		panic(err)
	}

	send_signal, err := strconv.Atoi(tile.SendSignal)
	if err != nil {
		panic(err)
	}

	new_water_tile = FloodTileJson{utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}, speed, send_signal, receive_signal}

	return new_water_tile
}

func (flood_tile *FloodTileJson) Deserialize(tile *int) FloodTile {
	new_tile := FloodTile{}

	new_tile.Tile = tile
	new_tile.Speed = strconv.Itoa(flood_tile.Speed)
	new_tile.SendSignal = strconv.Itoa(flood_tile.SendSignal)
	new_tile.ReceiveSignal = strconv.Itoa(flood_tile.ReceiveSignal)

	return new_tile
}

func (level *Level) ManageFloodTiles() {
	for i, flood_tiles := range level.FloodTiles {
		if *flood_tiles.Tile != -13 {
			utils.RemoveArrayElement(i, &level.FloodTiles)
			*flood_tiles.Tile = 0
		}
	}
}

func (level *Level) PlaceFloodTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -13

					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					for i := range level.FloodTiles {
						t := &level.FloodTiles[i]
						if t.Tile == tile {
							can_add = false
						}
					}

					if can_add {
						if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -13 {
							level.FloodTiles = append(level.FloodTiles, FloodTile{&level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)], "0", "0", "0"})
						}
					}

					level.Level_In_Matrix[chunk_y][chunk_x].Changed = 2
				}
			}
		}
	}
}
func (level *Level) SelectFloodTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					if *tile == -13 {
						for i := range level.FloodTiles {
							t := &level.FloodTiles[i].Tile
							if *t == tile {
								can_add = false
								SelectedFloodtile = &level.FloodTiles[i]
							}
						}

						if can_add {
							if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -13 {
								level.FloodTiles = append(level.FloodTiles, FloodTile{tile, "0", "0", "0"})
								SelectedFloodtile = &level.FloodTiles[len(level.FloodTiles)-1]
							}
						}
					}
				}
			}
		}
	}
}
