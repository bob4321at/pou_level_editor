package grid

import (
	"fmt"
	"main/utils"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SpringTile struct {
	Tile      *int
	Power     string
	Direction int
}

type SpringTileJson struct {
	Pos       utils.Vec2
	Power     float64
	Direction int
}

var SpringTileImg []*ebiten.Image
var SelectedSpringTile *SpringTile

func init() {

	imga, _, err := ebitenutil.NewImageFromFile("./art/springup.png")
	if err != nil {
		panic(err)
	}
	imgb, _, err := ebitenutil.NewImageFromFile("./art/springright.png")
	if err != nil {
		panic(err)
	}
	imgc, _, err := ebitenutil.NewImageFromFile("./art/springdown.png")
	if err != nil {
		panic(err)
	}
	imgd, _, err := ebitenutil.NewImageFromFile("./art/springleft.png")
	if err != nil {
		panic(err)
	}

	SpringTileImg = []*ebiten.Image{
		imga,
		imgb,
		imgc,
		imgd,
	}
}

func (tile *SpringTile) Serialize(chunk_x, chunk_y, tile_x, tile_y int) SpringTileJson {
	new_spring_tile := SpringTileJson{}

	power, err := strconv.ParseFloat(tile.Power, 64)
	if err != nil {
		panic(err)
	}

	new_spring_tile = SpringTileJson{utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}, power, tile.Direction}

	return new_spring_tile
}

func (spring_tile *SpringTileJson) Deserialize(tile *int) SpringTile {
	new_tile := SpringTile{}

	new_tile.Direction = spring_tile.Direction
	new_tile.Power = fmt.Sprintf("%f", spring_tile.Power)
	new_tile.Tile = tile

	return new_tile
}

func (level *Level) ManageSpringTiles() {
	for i, spike_tile := range level.SpringTiles {
		if *spike_tile.Tile != -9 {
			utils.RemoveArrayElement(i, &level.SpringTiles)
			*spike_tile.Tile = 0
		}
	}
}

func (level *Level) PlaceSpringTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -9

					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					for i := range level.SpringTiles {
						t := &level.SpringTiles[i]
						if t.Tile == tile {
							can_add = false
						}
					}

					if can_add {
						if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -9 {
							level.SpringTiles = append(level.SpringTiles, SpringTile{&level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)], "1", 0})
						}
					}

					level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
				}
			}
		}
	}
}
func (level *Level) SelectSpringTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					if *tile == -9 {
						for i := range level.SpringTiles {
							t := &level.SpringTiles[i].Tile
							if *t == tile {
								can_add = false
								SelectedSpringTile = &level.SpringTiles[i]
								SelectedEnemySpawner = nil
								SelectedBreakableTile = nil
							}
						}

						if can_add {
							if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -9 {
								level.SpringTiles = append(level.SpringTiles, SpringTile{tile, "1", 0})
								SelectedSpringTile = &level.SpringTiles[len(level.SpringTiles)-1]
							}
						}
					}
				}
			}
		}
	}
}
