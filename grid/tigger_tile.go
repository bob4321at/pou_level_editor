package grid

import (
	"main/utils"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TriggerTile struct {
	Tile      *int
	Signal    string
	Visible   bool
	Direction int
}

type TriggerTileJson struct {
	Pos       utils.Vec2
	Signal    int
	Visible   bool
	Direction int
}

var TriggerTileImg, _, _ = ebitenutil.NewImageFromFile("./art/trigger_tile.png")
var ButtonTileImgs []*ebiten.Image
var SelectedTriggerTile *TriggerTile

func init() {

	imga, _, err := ebitenutil.NewImageFromFile("./art/buttonup.png")
	if err != nil {
		panic(err)
	}
	imgb, _, err := ebitenutil.NewImageFromFile("./art/buttonright.png")
	if err != nil {
		panic(err)
	}
	imgc, _, err := ebitenutil.NewImageFromFile("./art/buttondown.png")
	if err != nil {
		panic(err)
	}
	imgd, _, err := ebitenutil.NewImageFromFile("./art/buttonleft.png")
	if err != nil {
		panic(err)
	}

	ButtonTileImgs = []*ebiten.Image{
		imga,
		imgb,
		imgc,
		imgd,
	}
}

func (tile *TriggerTile) Serialize(chunk_x, chunk_y, tile_x, tile_y int) TriggerTileJson {
	new_trigger_tile := TriggerTileJson{}

	signal, err := strconv.Atoi(tile.Signal)
	if err != nil {
		panic(err)
	}

	new_trigger_tile = TriggerTileJson{utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}, signal, tile.Visible, tile.Direction}

	return new_trigger_tile
}

func (trigger_tile *TriggerTileJson) Deserialize(tile *int) TriggerTile {
	new_tile := TriggerTile{}

	new_tile.Signal = strconv.Itoa(trigger_tile.Signal)
	new_tile.Tile = tile
	new_tile.Visible = trigger_tile.Visible
	new_tile.Direction = trigger_tile.Direction

	return new_tile
}

func (level *Level) ManageTriggerTiles() {
	for i, trigger_tile := range level.TriggerTile {
		if *trigger_tile.Tile != -5 {
			utils.RemoveArrayElement(i, &level.TriggerTile)
			*trigger_tile.Tile = 0
		}
	}
}

func (level *Level) PlaceTriggerTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -5

					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					for i := range level.TriggerTile {
						t := &level.TriggerTile[i]
						if t.Tile == tile {
							can_add = false
							SelectedTriggerTile = t
						}
					}

					if can_add {
						if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -5 {
							level.TriggerTile = append(level.TriggerTile, TriggerTile{&level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)], "0", false, 0})
							SelectedTriggerTile = &level.TriggerTile[len(level.TriggerTile)-1]
							SelectedEnemySpawner = nil
						}
					}

					level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
				}
			}
		}
	}
}

func (level *Level) SelectTriggerTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					if *tile == -5 {
						for i := range level.TriggerTile {
							t := &level.TriggerTile[i].Tile
							if *t == tile {
								can_add = false
								SelectedTriggerTile = &level.TriggerTile[i]
								SelectedEnemySpawner = nil
								SelectedBreakableTile = nil
							}
						}

						if can_add {
							if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -5 {
								level.TriggerTile = append(level.TriggerTile, TriggerTile{tile, "1", false, 0})
								SelectedTriggerTile = &level.TriggerTile[len(level.TriggerTile)-1]
							}
						}
					}
				}
			}
		}
	}
}
