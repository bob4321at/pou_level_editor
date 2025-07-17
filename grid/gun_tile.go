package grid

import (
	"main/utils"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GunTile struct {
	Tile          *int
	GunId         string
	SendSignal    string
	ReceiveSignal string
}

type GunTileJson struct {
	Pos           utils.Vec2
	GunId         int
	SendSignal    int
	ReceiveSignal int
}

var GunTileImg, _, _ = ebitenutil.NewImageFromFile("./art/gun_tile.png")
var SelectedGunTile *GunTile

func (tile *GunTile) Serialize(chunk_x, chunk_y, tile_x, tile_y int) GunTileJson {
	new_gun_tile := GunTileJson{}

	gun_id, err := strconv.Atoi(tile.GunId)
	if err != nil {
		panic(err)
	}

	send_signal, err := strconv.Atoi(tile.SendSignal)
	if err != nil {
		panic(err)
	}

	receive_signal, err := strconv.Atoi(tile.ReceiveSignal)
	if err != nil {
		panic(err)
	}

	new_gun_tile = GunTileJson{utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}, gun_id, send_signal, receive_signal}

	return new_gun_tile
}

func (gun_tile *GunTileJson) Deserialize(tile *int) GunTile {
	new_tile := GunTile{}

	new_tile.GunId = strconv.Itoa(gun_tile.GunId)
	new_tile.SendSignal = strconv.Itoa(gun_tile.SendSignal)
	new_tile.ReceiveSignal = strconv.Itoa(gun_tile.ReceiveSignal)
	new_tile.Tile = tile

	return new_tile
}

func (level *Level) ManageGunTiles() {
	for i, gun_tile := range level.GunTiles {
		if *gun_tile.Tile != -7 {
			utils.RemoveArrayElement(i, &level.GunTiles)
			*gun_tile.Tile = 0
		}
	}
}

func (level *Level) PlaceGunTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -7

					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					for i := range level.GunTiles {
						t := &level.GunTiles[i]
						if t.Tile == tile {
							can_add = false
							SelectedGunTile = t
						}
					}

					if can_add {
						if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -7 {
							level.GunTiles = append(level.GunTiles, GunTile{&level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)], "0", "0", "0"})
							SelectedGunTile = &level.GunTiles[len(level.GunTiles)-1]
						}
					}

					level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
				}
			}
		}
	}
}

func (level *Level) SelectGunTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					if *tile == -7 {
						for i := range level.GunTiles {
							t := &level.GunTiles[i].Tile
							if *t == tile {
								can_add = false
								SelectedGunTile = &level.GunTiles[i]
							}
						}

						if can_add {
							if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -4 {
								level.GunTiles = append(level.GunTiles, GunTile{tile, "0", "0", "0"})
								SelectedGunTile = &level.GunTiles[len(level.GunTiles)-1]
							}
						}
					}
				}
			}
		}
	}
}
