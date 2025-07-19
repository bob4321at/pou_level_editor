package grid

import (
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SpikeTile struct {
	Tile      *int
	Direction int
}

type SpikeTileJson struct {
	Pos       utils.Vec2
	Direction int
}

var SpikeTileImg []*ebiten.Image

func init() {
	imga, _, err := ebitenutil.NewImageFromFile("./art/spikeup.png")
	if err != nil {
		panic(err)
	}
	imgb, _, err := ebitenutil.NewImageFromFile("./art/spikeright.png")
	if err != nil {
		panic(err)
	}
	imgc, _, err := ebitenutil.NewImageFromFile("./art/spikedown.png")
	if err != nil {
		panic(err)
	}
	imgd, _, err := ebitenutil.NewImageFromFile("./art/spikeleft.png")
	if err != nil {
		panic(err)
	}

	SpikeTileImg = []*ebiten.Image{
		imga,
		imgb,
		imgc,
		imgd,
	}

}

func (tile *SpikeTile) Serialize(chunk_x, chunk_y, tile_x, tile_y int) SpikeTileJson {
	new_spike_tile := SpikeTileJson{}

	new_spike_tile = SpikeTileJson{utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}, tile.Direction}

	return new_spike_tile
}

func (spike_tile *SpikeTileJson) Deserialize(tile *int) SpikeTile {
	new_tile := SpikeTile{}

	new_tile.Direction = spike_tile.Direction
	new_tile.Tile = tile

	return new_tile
}

func (level *Level) ManageSpikeTiles() {
	for i, spike_tile := range level.SpikeTiles {
		if *spike_tile.Tile != -8 {
			utils.RemoveArrayElement(i, &level.SpikeTiles)
			*spike_tile.Tile = 0
		}
	}
}

func (level *Level) PlaceSpikeTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -8

					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					for i := range level.SpikeTiles {
						t := &level.SpikeTiles[i]
						if t.Tile == tile {
							can_add = false
						}
					}

					if can_add {
						if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -8 {
							level.SpikeTiles = append(level.SpikeTiles, SpikeTile{&level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)], 0})
						}
					}

					level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
				}
			}
		}
	}
}
