package grid

import (
	"main/utils"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type EnemySpawner struct {
	Tile          *int
	Enemies       []string
	SendSignal    string
	ReceiveSignal string
}

type EnemySpawnerJson struct {
	Pos           utils.Vec2
	Enemies       []int
	SendSignal    int
	ReceiveSignal int
}

var EnemySpawnerTile, _, _ = ebitenutil.NewImageFromFile("./art/enemy_spawner.png")
var SelectedEnemySpawner *EnemySpawner

func (spawner *EnemySpawner) Serialize(chunk_x, chunk_y, tile_x, tile_y int) EnemySpawnerJson {
	new_spawner := EnemySpawnerJson{}

	var converted_enemies []int

	for _, enemy := range spawner.Enemies {
		num, err := strconv.Atoi(enemy)
		if err != nil {
			panic(err)
		}
		converted_enemies = append(converted_enemies, num)
	}

	var send_signal int

	if spawner.SendSignal != "" {
		var err error
		send_signal, err = strconv.Atoi(spawner.SendSignal)
		if err != nil {
			panic(err)
		}
	}

	var receive_signal int

	if spawner.ReceiveSignal != "" {
		var err error
		receive_signal, err = strconv.Atoi(spawner.ReceiveSignal)
		if err != nil {
			panic(err)
		}
	}

	new_spawner = EnemySpawnerJson{utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}, converted_enemies, send_signal, receive_signal}

	return new_spawner
}

func (spawner *EnemySpawnerJson) Deserialize(tile *int) EnemySpawner {
	new_spawner := EnemySpawner{}
	new_spawner.SendSignal = strconv.Itoa(spawner.SendSignal)
	new_spawner.ReceiveSignal = strconv.Itoa(spawner.ReceiveSignal)
	new_spawner.Tile = tile

	for _, enemy := range spawner.Enemies {
		new_spawner.Enemies = append(new_spawner.Enemies, strconv.Itoa(enemy))
	}

	return new_spawner
}

func (level *Level) ManageEnemySpawners() {
	for i, enemy_spawner := range level.Enemy_Spawner {
		if *enemy_spawner.Tile != -3 {
			utils.RemoveArrayElement(i, &level.Enemy_Spawner)
			*enemy_spawner.Tile = 0
		}
	}
}

func (level *Level) PlaceEnemySpawner(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)
	if chunk_x < len(level.Level_In_Matrix[0]) {
		if chunk_y < len(level.Level_In_Matrix) {
			if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
				if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
					level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -3

					can_add := true

					tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

					for i := range level.Enemy_Spawner {
						t := &level.Enemy_Spawner[i]
						if t.Tile == tile {
							can_add = false
							SelectedEnemySpawner = t
						}
					}

					if can_add {
						if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -3 {
							level.Enemy_Spawner = append(level.Enemy_Spawner, EnemySpawner{&level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)], []string{}, "", ""})
							SelectedEnemySpawner = &level.Enemy_Spawner[len(level.Enemy_Spawner)-1]
							SelectedBreakableTile = nil
						}
					}

					level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
				}
			}
		}
	}
}

func (level *Level) SelectEnemySpawner(world_cord_x, world_cord_y float64) {
	chunk_x := int((world_cord_x / 32) / 32)
	chunk_y := int((world_cord_y / 32) / 32)

	can_add := true

	tile := &level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)]

	if *tile == -3 {
		for i := range level.Enemy_Spawner {
			t := &level.Enemy_Spawner[i]
			if t.Tile == tile {
				can_add = false
				SelectedEnemySpawner = t
				SelectedBreakableTile = nil
				SelectedTriggerTile = nil
			}
		}

		if can_add {
			if level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] == -3 {
				level.Enemy_Spawner = append(level.Enemy_Spawner, EnemySpawner{tile, []string{}, "", ""})
				SelectedEnemySpawner = &level.Enemy_Spawner[len(level.Enemy_Spawner)-1]
			}
		}
	}
}
