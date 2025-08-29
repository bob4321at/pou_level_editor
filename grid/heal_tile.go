package grid

import (
	"main/utils"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ItemTile struct {
	Tile          *int
	ItemId        string
	CatagoryId    string
	SendSignal    string
	ReceiveSignal string
}

type ItemTileJson struct {
	Pos           utils.Vec2
	ItemId        int
	CatagoryId    int
	SendSignal    int
	ReceiveSignal int
}

var ItemTileImg, _, _ = ebitenutil.NewImageFromFile("./art/item_tile.png")
var SelectedItemTile *ItemTile

func (tile *ItemTile) Serialize(chunk_x, chunk_y, tile_x, tile_y int) ItemTileJson {
	new_item := ItemTileJson{}

	item_id, err := strconv.Atoi(tile.ItemId)
	if err != nil {
		panic(err)
	}

	catagory_id, err := strconv.Atoi(tile.CatagoryId)
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

	new_item = ItemTileJson{utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}, item_id, catagory_id, send_signal, receive_signal}

	return new_item
}

func (gun_tile *ItemTileJson) Deserialize(tile *int) ItemTile {
	new_tile := ItemTile{}

	new_tile.ItemId = strconv.Itoa(gun_tile.ItemId)
	new_tile.CatagoryId = strconv.Itoa(gun_tile.CatagoryId)
	new_tile.SendSignal = strconv.Itoa(gun_tile.SendSignal)
	new_tile.ReceiveSignal = strconv.Itoa(gun_tile.ReceiveSignal)
	new_tile.Tile = tile

	return new_tile
}

func (level *Level) ManageItemTiles() {
	for i, item_tile := range level.ItemTiles {
		if *item_tile.Tile != -10 {
			utils.RemoveArrayElement(i, &level.ItemTiles)
			*item_tile.Tile = 0
		}
	}
}

func (level *Level) PlaceItemTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -10

					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					for i := range level.ItemTiles {
						t := &level.ItemTiles[i]
						if t.Tile == tile {
							can_add = false
							SelectedItemTile = t
						}
					}

					if can_add {
						if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -10 {
							level.ItemTiles = append(level.ItemTiles, ItemTile{&level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)], "0", "0", "0", "0"})
							SelectedItemTile = &level.ItemTiles[len(level.ItemTiles)-1]
						}
					}

					level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
				}
			}
		}
	}
}

func (level *Level) SelectItemTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					if *tile == -10 {
						for i := range level.ItemTiles {
							t := &level.ItemTiles[i].Tile
							if *t == tile {
								can_add = false
								SelectedItemTile = &level.ItemTiles[i]
							}
						}

						if can_add {
							if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -10 {
								level.ItemTiles = append(level.ItemTiles, ItemTile{tile, "0", "0", "0", "0"})
								SelectedItemTile = &level.ItemTiles[len(level.ItemTiles)-1]
							}
						}
					}
				}
			}
		}
	}
}
