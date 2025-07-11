package ui

import (
	"fmt"
	"image"
	"main/camera"
	"main/grid"
	"main/utils"
	"strconv"

	"github.com/ebitengine/debugui"
)

var X_Size string = "32"
var Y_Size string = "32"

var Save_Name string
var Load_Name string

func Menu(ctx *debugui.Context) {
	ctx.Window("test", image.Rect(0, 0, 100, 300), func(layout debugui.ContainerLayout) {
		ctx.Header("New Level", false, func() {
			ctx.TextField(&X_Size)
			ctx.TextField(&Y_Size)
			ctx.Button("Create").On(func() {
				width, err := strconv.Atoi(X_Size)
				if err != nil {
					panic(err)
				}
				height, err := strconv.Atoi(Y_Size)
				if err != nil {
					panic(err)
				}
				camera.Camera.Pos = utils.Vec2{X: float64(width*16) / 2, Y: float64(height*16) / 2}

				grid.Current_Level = grid.NewLevel(width, height, "./art/tile_set.png")
			})
		})
		ctx.TextField(&Save_Name)
		ctx.Button("Save").On(func() {
			grid.Current_Level.Save(Save_Name)
		})
		ctx.TextField(&Load_Name)
		ctx.Button("Load").On(func() {
			grid.LoadLevel(Load_Name)
		})
	})
}

func EditTriggerTileUi(ctx *debugui.Context) {
	ctx.Window("Modify Trigger Tile Spawner", image.Rect(400, 0, 640, 320), func(layout debugui.ContainerLayout) {
		test := &grid.SelectedTriggerTile.Signal
		ctx.Text("Send")
		ctx.TextField(test)

		ctx.Button("Close Window").On(func() {
			grid.SelectedTriggerTile = nil
		})
	})
}

func EditBreakableTileUi(ctx *debugui.Context) {
	ctx.Window("Modify Breakable Tile Spawner", image.Rect(400, 0, 640, 320), func(layout debugui.ContainerLayout) {
		test := &grid.SelectedBreakableTile.Signal
		ctx.TextField(test)

		ctx.Button("Close Window").On(func() {
			grid.SelectedBreakableTile = nil
		})
	})
}

func EditEnemySpawnerUi(ctx *debugui.Context) {
	ctx.Window("Modify Enemy Spawner", image.Rect(400, 0, 640, 320), func(layout debugui.ContainerLayout) {
		ctx.Text("Enemies")

		for i := range grid.SelectedEnemySpawner.Enemies {
			ctx.IDScope(fmt.Sprintf("%d", i), func() {
				ctx.TextField(&grid.SelectedEnemySpawner.Enemies[i])
			})
		}

		ctx.Button("Add Enemy").On(func() {
			grid.SelectedEnemySpawner.Enemies = append(grid.SelectedEnemySpawner.Enemies, "0")
		})

		ctx.Text("Signals")
		ctx.Text("Send")
		ctx.TextField(&grid.SelectedEnemySpawner.SendSignal)
		ctx.Text("Receive")
		ctx.TextField(&grid.SelectedEnemySpawner.ReceiveSignal)

		ctx.Button("Close Window").On(func() {
			grid.SelectedEnemySpawner = nil
		})
	})
}
