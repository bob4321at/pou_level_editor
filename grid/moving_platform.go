package grid

import (
	"main/utils"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type MovingPlatformTile struct {
	Tile       *int
	Signal     string
	Track      string
	TrackIndex string
	Loop       bool
}

type MovingPlatformTileJson struct {
	Pos        utils.Vec2
	Signal     int
	Track      int
	TrackIndex int
	Loop       bool
}

var MovingPlatformTileImg, _, _ = ebitenutil.NewImageFromFile("./art/moving_platform_tile.png")
var SelectedMovingPlatformTile *MovingPlatformTile

func (tile *MovingPlatformTile) Serialize(chunk_x, chunk_y, tile_x, tile_y int) MovingPlatformTileJson {
	new_moving_platform_tile := MovingPlatformTileJson{}

	signal, err := strconv.Atoi(tile.Signal)
	if err != nil {
		panic(err)
	}

	track, err := strconv.Atoi(tile.Track)
	if err != nil {
		panic(err)
	}

	track_index, err := strconv.Atoi(tile.TrackIndex)
	if err != nil {
		panic(err)
	}

	new_moving_platform_tile = MovingPlatformTileJson{utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}, signal, track, track_index, tile.Loop}

	return new_moving_platform_tile
}

func (moving_platform_tile *MovingPlatformTileJson) Deserialize(tile *int) MovingPlatformTile {
	new_tile := MovingPlatformTile{}

	new_tile.Signal = strconv.Itoa(moving_platform_tile.Signal)
	new_tile.Tile = tile
	new_tile.Track = strconv.Itoa(moving_platform_tile.Track)
	new_tile.TrackIndex = strconv.Itoa(moving_platform_tile.TrackIndex)
	new_tile.Loop = moving_platform_tile.Loop

	return new_tile
}

func (level *Level) ManageMovingPlatformTiles() {
	for i, moving_platform_tiles := range level.MovingPlatformTiles {
		if *moving_platform_tiles.Tile != -11 {
			utils.RemoveArrayElement(i, &level.MovingPlatformTiles)
			*moving_platform_tiles.Tile = 0
		}
	}
}

func (level *Level) PlaceMovingPlatformTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -11

					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					for i := range level.MovingPlatformTiles {
						t := &level.MovingPlatformTiles[i]
						if t.Tile == tile {
							can_add = false
							SelectedMovingPlatformTile = t
						}
					}

					if can_add {
						if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -11 {
							level.MovingPlatformTiles = append(level.MovingPlatformTiles, MovingPlatformTile{&level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)], "0", "0", "0", false})
							SelectedMovingPlatformTile = &level.MovingPlatformTiles[len(level.MovingPlatformTiles)-1]
							SelectedEnemySpawner = nil
						}
					}

					level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
				}
			}
		}
	}
}

func (level *Level) SelectMovingTile(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					if *tile == -11 {
						for i := range level.MovingPlatformTiles {
							t := &level.MovingPlatformTiles[i].Tile
							if *t == tile {
								can_add = false
								SelectedMovingPlatformTile = &level.MovingPlatformTiles[i]
							}
						}

						if can_add {
							if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -11 {
								level.MovingPlatformTiles = append(level.MovingPlatformTiles, MovingPlatformTile{tile, "1", "0", "0", false})
								SelectedMovingPlatformTile = &level.MovingPlatformTiles[len(level.TriggerTile)-1]
							}
						}
					}
				}
			}
		}
	}
}
