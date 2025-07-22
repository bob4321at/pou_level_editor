package grid

import (
	"encoding/json"
	"image/color"
	"main/camera"
	"main/shader"
	"main/utils"
	"math"
	"os"

	"github.com/bob4321at/textures"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Chunk struct {
	Tiles     [][]int
	ShaderImg textures.Texture
	Changed   bool
}

var PlayerTile, _, _ = ebitenutil.NewImageFromFile("./art/player_tile.png")
var SockTile, _, _ = ebitenutil.NewImageFromFile("./art/sock.png")

var Right_Mouse_Just_Pressed int = 0

var R float64 = 1
var G float64 = 0.4
var B float64 = 0.4

var RR float64 = 1
var GG float64 = 0.5
var BB float64 = 0.5

var Background_Red float64 = 1
var Background_Green float64 = 0.7
var Background_Blue float64 = 0.7

func (chunk *Chunk) GenCache(chunk_x, chunk_y int) {
	for y, row := range chunk.Tiles {
		for x, tile := range row {
			if tile > 0 {
				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x)*32, float64(y)*32)

				Current_Level.TileSetWithShader[tile-1].Draw(chunk.ShaderImg.Img, &op)
			}
			if tile == -1 {
				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x)*32, float64(y)*32)
				op.Blend = ebiten.BlendClear

				chunk.ShaderImg.Img.DrawImage(ebiten.NewImage(32, 32), &op)

				chunk.Tiles[y][x] = 0
			}
			if tile == -2 {
				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x)*32, float64(y)*32)

				chunk.ShaderImg.Img.DrawImage(PlayerTile, &op)
			}
			if tile == -3 {
				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x)*32, float64(y)*32)

				chunk.ShaderImg.Img.DrawImage(EnemySpawnerTile, &op)
			}
			if tile == -4 {
				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x)*32, float64(y)*32)

				chunk.ShaderImg.Img.DrawImage(BreakableTIleImg, &op)
			}
			if tile == -5 {
				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x)*32, float64(y)*32)

				chunk.ShaderImg.Img.DrawImage(TriggerTileImg, &op)
			}
			if tile == -6 {
				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x)*32, float64(y)*32)

				chunk.ShaderImg.Img.DrawImage(SockTile, &op)
			}
			if tile == -7 {
				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x)*32, float64(y)*32)

				chunk.ShaderImg.Img.DrawImage(GunTileImg, &op)
			}
			if tile == -8 {
				for i := range Current_Level.SpikeTiles {
					spike := Current_Level.SpikeTiles[i].Tile
					tile := &Current_Level.Level_In_Matrix[chunk_y][chunk_x].Tiles[y][x]
					if spike == tile {
						op := ebiten.DrawImageOptions{}
						op.GeoM.Translate(float64(x)*32, float64(y)*32)
						op.Blend = ebiten.BlendClear

						chunk.ShaderImg.Img.DrawImage(ebiten.NewImage(32, 32), &op)

						op.GeoM.Reset()
						op.Blend = ebiten.BlendCopy
						op.GeoM.Translate(float64(x)*32, float64(y)*32)

						chunk.ShaderImg.Img.DrawImage(SpikeTileImg[Current_Level.SpikeTiles[i].Direction], &op)
					}
				}
			}
		}
	}

	chunk.Changed = false
}

type Level struct {
	Size              utils.Vec2
	TileSet_Img       *ebiten.Image
	TileSet           []*ebiten.Image
	TileSetWithShader []textures.RenderableTexture

	Level_In_Matrix [][]Chunk

	Chunks_Created bool

	Enemy_Spawner []EnemySpawner
	BreakableTile []BreakableTile
	TriggerTile   []TriggerTile
	GunTiles      []GunTile
	SpikeTiles    []SpikeTile
}

type Tile struct {
	ID  int
	Pos utils.Vec2
}

type LevelJson struct {
	Player_Spawn  utils.Vec2
	End           utils.Vec2
	Tiles         []Tile
	Enemies       []EnemySpawnerJson
	BreakableTile []BreakableTileJson
	TriggerTile   []TriggerTileJson
	GunTiles      []GunTileJson
	SpikeTiles    []SpikeTileJson

	TileBorderColor color.RGBA
	TileColor       color.RGBA
	BackgroundColor color.RGBA
}

var Grid_Img *ebiten.Image

func init() {
	temp_img, _, err := ebitenutil.NewImageFromFile("./art/grid.png")
	if err != nil {
		panic(err)
	}
	Grid_Img = temp_img
}

func (level *Level) Update() {

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && Right_Mouse_Just_Pressed == 1 {
		Right_Mouse_Just_Pressed = 2
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && Right_Mouse_Just_Pressed == 0 {
		Right_Mouse_Just_Pressed = 1
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		Right_Mouse_Just_Pressed = 0
	}

	world_cord_x := utils.Mouse_X + camera.Camera.Pos.X
	world_cord_y := utils.Mouse_Y + camera.Camera.Pos.Y

	level.ManageEnemySpawners()
	level.ManageBreakableTiles()
	level.ManageTriggerTiles()
	level.ManageGunTiles()
	level.ManageSpikeTiles()

	if level.Chunks_Created {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			chunk_x := int((world_cord_x / 32) / 32)
			chunk_y := int((world_cord_y / 32) / 32)
			if chunk_x < len(level.Level_In_Matrix[0]) {
				if chunk_y < len(level.Level_In_Matrix) {
					if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
						if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
							level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = 1
							level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
						}
					}
				}
			}
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
			chunk_x := int((world_cord_x / 32) / 32)
			chunk_y := int((world_cord_y / 32) / 32)
			if chunk_x < len(level.Level_In_Matrix[0]) {
				if chunk_y < len(level.Level_In_Matrix) {
					if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
						if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
							level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -1
							level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
						}
					}
				}
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyP) {
			chunk_x := int((world_cord_x / 32) / 32)
			chunk_y := int((world_cord_y / 32) / 32)
			if chunk_x < len(level.Level_In_Matrix[0]) {
				if chunk_y < len(level.Level_In_Matrix) {
					if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
						if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
							level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -2
							level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
						}
					}
				}
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			chunk_x := int((world_cord_x / 32) / 32)
			chunk_y := int((world_cord_y / 32) / 32)
			if chunk_x < len(level.Level_In_Matrix[0]) {
				if chunk_y < len(level.Level_In_Matrix) {
					if (chunk_x*32*32) < int(world_cord_x) && (chunk_x*32*32)+(32*32) > int(world_cord_x) {
						if (chunk_y*32*32) < int(world_cord_y) && (chunk_y*32*32)+(32*32) > int(world_cord_y) {
							level.Level_In_Matrix[chunk_y][chunk_x].Tiles[(int(world_cord_y)/32)-(chunk_y*32)][(int(world_cord_x)/32)-(chunk_x*32)] = -6
							level.Level_In_Matrix[chunk_y][chunk_x].Changed = true
						}
					}
				}
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyB) {
			level.PlaceBreakableTile(world_cord_x, world_cord_y)
		}
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			level.PlaceEnemySpawner(world_cord_x, world_cord_y)
		}
		if ebiten.IsKeyPressed(ebiten.KeyT) {
			level.PlaceTriggerTile(world_cord_x, world_cord_y)
		}
		if ebiten.IsKeyPressed(ebiten.KeyG) {
			level.PlaceGunTile(world_cord_x, world_cord_y)
		}
		if ebiten.IsKeyPressed(ebiten.KeyC) {
			level.PlaceSpikeTile(world_cord_x, world_cord_y)
		}
		if Right_Mouse_Just_Pressed == 1 {
			SelectedBreakableTile = nil
			SelectedEnemySpawner = nil
			SelectedTriggerTile = nil
			SelectedGunTile = nil
			SelectedSpikeTile = nil
			level.SelectTriggerTile(world_cord_x, world_cord_y)
			level.SelectEnemySpawner(world_cord_x, world_cord_y)
			level.SelectBreakableTile(world_cord_x, world_cord_y)
			level.SelectGunTile(world_cord_x, world_cord_y)
			level.SelectSpikeTile(world_cord_x, world_cord_y)
		}
	}

	level.NeighbourCheck()

	for chunk_y, chuck_rows := range level.Level_In_Matrix {
		for chunk_x, _ := range chuck_rows {
			if level.Level_In_Matrix[chunk_y][chunk_x].Changed {
				level.Level_In_Matrix[chunk_y][chunk_x].GenCache(chunk_x, chunk_y)
			}
		}
	}
}

func (level *Level) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}

	for y, _ := range level.Level_In_Matrix {
		for x := range level.Level_In_Matrix[y] {
			op.GeoM.Reset()
			op.GeoM.Translate((float64(x)+float64(x*1024))-(camera.Camera.Pos.X)-(camera.Camera.Start_Move_Pos.X-camera.Camera.Move_Pos.X), (float64(y)+float64(y*1024))-(camera.Camera.Pos.Y)-(camera.Camera.Start_Move_Pos.Y-camera.Camera.Move_Pos.Y))

			level.Level_In_Matrix[y][x].ShaderImg.SetUniforms(map[string]any{
				"R": R,
				"G": G,
				"B": B,

				"RR": RR,
				"GG": GG,
				"BB": BB,
			})
			level.Level_In_Matrix[y][x].ShaderImg.Draw(screen, &op)
			screen.DrawImage(Grid_Img, &op)
		}
	}

}

func NewLevel(width, height int, tileset_path string) (level Level) {
	level.Size = utils.Vec2{X: float64(width), Y: float64(height)}

	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/top_left.png", ""))
	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/top_center.png", ""))
	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/top_right.png", ""))

	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/middle_left.png", ""))
	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/middle_center.png", ""))
	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/middle_right.png", ""))

	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/bottom_left.png", ""))
	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/bottom_center.png", ""))
	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/bottom_right.png", ""))

	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/vertical_top.png", ""))
	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/vertical_middle.png", ""))
	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/vertical_bottom.png", ""))

	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/horizontal_left.png", ""))
	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/horizontal_center.png", ""))
	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/horizontal_right.png", ""))

	level.TileSetWithShader = append(level.TileSetWithShader, textures.NewTexture("./art/tileset/center.png", ""))

	level.Chunks_Created = false

	chuncks_y := level.Size.Y / 32
	chuncks_x := level.Size.X / 32

	if math.Mod(level.Size.X, 32) != 0 {
		chuncks_x = (level.Size.X + 32) / 32
	}
	if math.Mod(level.Size.Y, 32) != 0 {
		chuncks_y = (level.Size.Y + 32) / 32
	}
	for chunk_y := range int(chuncks_y) {
		level.Level_In_Matrix = append(level.Level_In_Matrix, []Chunk{})
		for _ = range int(chuncks_x) {
			var empty_Chunk [][]int
			empty_Chunk = make([][]int, 32)

			for i := range empty_Chunk {
				for _ = range 32 {
					empty_Chunk[i] = append(empty_Chunk[i], 0)
				}
			}

			level.Level_In_Matrix[chunk_y] = append(level.Level_In_Matrix[chunk_y], Chunk{empty_Chunk, *textures.NewTexture("./art/empty_chunk.png", shader.Chunk_Shader), false})
		}
	}

	level.Chunks_Created = true

	return level
}

func (level *Level) Save(name string) {
	tiles := LevelJson{}

	for chunk_y, chunk_row := range Current_Level.Level_In_Matrix {
		for chunk_x, chunk := range chunk_row {
			for tile_y, tile_row := range chunk.Tiles {
				for tile_x, tile := range tile_row {
					if tile != 0 {
						if tile > 0 {
							tiles.Tiles = append(tiles.Tiles, Tile{ID: tile, Pos: utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}})
						}
						if tile == -2 {
							tiles.Player_Spawn = utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}
						}
						if tile == -3 {
							for i := range level.Enemy_Spawner {
								if level.Enemy_Spawner[i].Tile == &Current_Level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] {
									spawner := &level.Enemy_Spawner[i]
									tiles.Enemies = append(tiles.Enemies, spawner.Serialize(chunk_x, chunk_y, tile_x, tile_y))
								}
							}
						}
						if tile == -4 {
							for i := range level.BreakableTile {
								if level.BreakableTile[i].Tile == &Current_Level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] {
									breakable_tile := &level.BreakableTile[i]
									tiles.BreakableTile = append(tiles.BreakableTile, breakable_tile.Serialize(chunk_x, chunk_y, tile_x, tile_y))
								}
							}
						}
						if tile == -5 {
							for i := range level.TriggerTile {
								if level.TriggerTile[i].Tile == &Current_Level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] {
									trigger_tile := &level.TriggerTile[i]
									tiles.TriggerTile = append(tiles.TriggerTile, trigger_tile.Serialize(chunk_x, chunk_y, tile_x, tile_y))
								}
							}
						}
						if tile == -6 {
							tiles.End = utils.Vec2{X: float64(chunk_x*1024) + float64(tile_x)*32, Y: float64(chunk_y*1024) + float64(tile_y)*32}
						}
						if tile == -7 {
							for i := range level.GunTiles {
								if level.GunTiles[i].Tile == &Current_Level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] {
									gun_tile := &level.GunTiles[i]
									tiles.GunTiles = append(tiles.GunTiles, gun_tile.Serialize(chunk_x, chunk_y, tile_x, tile_y))
								}
							}
						}
						if tile == -8 {
							for i := range level.SpikeTiles {
								if level.SpikeTiles[i].Tile == &Current_Level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] {
									spike_tile := &level.SpikeTiles[i]
									tiles.SpikeTiles = append(tiles.SpikeTiles, spike_tile.Serialize(chunk_x, chunk_y, tile_x, tile_y))
								}
							}
						}
					}
				}
			}
		}
	}

	tiles.TileBorderColor = color.RGBA{uint8(R * 255), uint8(G * 255), uint8(B * 255), 255}
	tiles.TileColor = color.RGBA{uint8(RR * 255), uint8(GG * 255), uint8(BB * 255), 255}
	tiles.BackgroundColor = color.RGBA{uint8(Background_Red * 255), uint8(Background_Green * 255), uint8(Background_Blue * 255), 255}

	data, err := json.Marshal(tiles)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("./levels/" + name)
	if err != nil {
		panic(err)
	}

	f.Write(data)
}

func LoadLevel(name string) {
	level := LevelJson{}

	level_file, err := os.ReadFile("./levels/" + name)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(level_file, &level); err != nil {
		panic(err)
	}

	Current_Level = NewLevel(300, 300, "./art/tile_set.png")

	for _, tile := range level.Tiles {
		Current_Level.Level_In_Matrix[int(tile.Pos.Y/1024)][int(tile.Pos.X/1024)].Tiles[int(math.Mod(tile.Pos.Y/32, 32))][int(math.Mod(tile.Pos.X/32, 32))] = 1
	}

	for _, enemy_spawner := range level.Enemies {
		Current_Level.Level_In_Matrix[int(enemy_spawner.Pos.Y/1024)][int(enemy_spawner.Pos.X/1024)].Tiles[int(math.Mod((enemy_spawner.Pos.Y/32), 32))][int(math.Mod(enemy_spawner.Pos.X/32, 32))] = -3
		tile := &Current_Level.Level_In_Matrix[int(enemy_spawner.Pos.Y/1024)][int(enemy_spawner.Pos.X/1024)].Tiles[int(math.Mod(enemy_spawner.Pos.Y/32, 32))][int(math.Mod(enemy_spawner.Pos.X/32, 32))]
		Current_Level.Enemy_Spawner = append(Current_Level.Enemy_Spawner, enemy_spawner.Deserialize(tile))
	}

	for _, breakable_tile := range level.BreakableTile {
		Current_Level.Level_In_Matrix[int(breakable_tile.Pos.Y/1024)][int(breakable_tile.Pos.X/1024)].Tiles[int(math.Mod((breakable_tile.Pos.Y/32), 32))][int(math.Mod(breakable_tile.Pos.X/32, 32))] = -4
		tile := &Current_Level.Level_In_Matrix[int(breakable_tile.Pos.Y/1024)][int(breakable_tile.Pos.X/1024)].Tiles[int(math.Mod(breakable_tile.Pos.Y/32, 32))][int(math.Mod(breakable_tile.Pos.X/32, 32))]
		Current_Level.BreakableTile = append(Current_Level.BreakableTile, breakable_tile.Deserialize(tile))
	}

	for _, trigger_tile := range level.TriggerTile {
		Current_Level.Level_In_Matrix[int(trigger_tile.Pos.Y/1024)][int(trigger_tile.Pos.X/1024)].Tiles[int(math.Mod((trigger_tile.Pos.Y/32), 32))][int(math.Mod(trigger_tile.Pos.X/32, 32))] = -5
		tile := &Current_Level.Level_In_Matrix[int(trigger_tile.Pos.Y/1024)][int(trigger_tile.Pos.X/1024)].Tiles[int(math.Mod(trigger_tile.Pos.Y/32, 32))][int(math.Mod(trigger_tile.Pos.X/32, 32))]
		Current_Level.TriggerTile = append(Current_Level.TriggerTile, trigger_tile.Deserialize(tile))
	}

	for _, gun_tile := range level.GunTiles {
		Current_Level.Level_In_Matrix[int(gun_tile.Pos.Y/1024)][int(gun_tile.Pos.X/1024)].Tiles[int(math.Mod((gun_tile.Pos.Y/32), 32))][int(math.Mod(gun_tile.Pos.X/32, 32))] = -7
		tile := &Current_Level.Level_In_Matrix[int(gun_tile.Pos.Y/1024)][int(gun_tile.Pos.X/1024)].Tiles[int(math.Mod(gun_tile.Pos.Y/32, 32))][int(math.Mod(gun_tile.Pos.X/32, 32))]
		Current_Level.GunTiles = append(Current_Level.GunTiles, gun_tile.Deserialize(tile))
	}

	for _, spike_tiles := range level.SpikeTiles {
		Current_Level.Level_In_Matrix[int(spike_tiles.Pos.Y/1024)][int(spike_tiles.Pos.X/1024)].Tiles[int(math.Mod((spike_tiles.Pos.Y/32), 32))][int(math.Mod(spike_tiles.Pos.X/32, 32))] = -8
		tile := &Current_Level.Level_In_Matrix[int(spike_tiles.Pos.Y/1024)][int(spike_tiles.Pos.X/1024)].Tiles[int(math.Mod(spike_tiles.Pos.Y/32, 32))][int(math.Mod(spike_tiles.Pos.X/32, 32))]
		Current_Level.SpikeTiles = append(Current_Level.SpikeTiles, spike_tiles.Deserialize(tile))
	}

	Current_Level.Level_In_Matrix[int(level.Player_Spawn.Y/1024)][int(level.Player_Spawn.X/1024)].Tiles[int(math.Mod((level.Player_Spawn.Y/32), 32))][int(math.Mod(level.Player_Spawn.X/32, 32))] = -2
	camera.Camera.Pos = utils.Vec2{X: level.Player_Spawn.X - 320, Y: level.Player_Spawn.Y - 240}

	Current_Level.Level_In_Matrix[int(level.End.Y/1024)][int(level.End.X/1024)].Tiles[int(math.Mod((level.End.Y/32), 32))][int(math.Mod(level.End.X/32, 32))] = -6

	R = float64(level.TileBorderColor.R) / 255
	G = float64(level.TileBorderColor.G) / 255
	B = float64(level.TileBorderColor.B) / 255

	RR = float64(level.TileColor.R) / 255
	GG = float64(level.TileColor.G) / 255
	BB = float64(level.TileColor.B) / 255

	Background_Red = float64(level.BackgroundColor.R) / 255
	Background_Green = float64(level.BackgroundColor.G) / 255
	Background_Blue = float64(level.BackgroundColor.B) / 255
}

var Current_Level = NewLevel(64, 64, "./art/tile_set.png")
