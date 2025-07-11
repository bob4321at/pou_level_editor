package grid

func (level *Level) NeighbourCheck() {
	for chunk_y, chunk_row := range level.Level_In_Matrix {
		for chunk_x, chunk := range chunk_row {
			for tile_y, tile_row := range chunk.Tiles {
				for tile_x, tile := range tile_row {
					if tile < 0 {
						tile = 0
					}
					if tile != 0 {
						var tile_above int = 0
						var tile_left int = 0
						var tile_right int = 0
						var tile_down int = 0

						if tile_y == 0 {
							if chunk_y == 0 {
								tile_above = 0
							} else {
								tile_above = level.Level_In_Matrix[chunk_y-1][chunk_x].Tiles[31][tile_x]
								level.Level_In_Matrix[chunk_y-1][chunk_x].Changed = true
							}
						} else {
							tile_above = level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y-1][tile_x]
						}

						if tile_y == 31 {
							if chunk_y == len(level.Level_In_Matrix)-1 {
								tile_down = 0
							} else {
								tile_down = level.Level_In_Matrix[chunk_y+1][chunk_x].Tiles[0][tile_x]
								level.Level_In_Matrix[chunk_y+1][chunk_x].Changed = true
							}
						} else {
							tile_down = level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y+1][tile_x]
						}

						if tile_x == 0 {
							if chunk_x == 0 {
								tile_left = 0
							} else {
								tile_left = level.Level_In_Matrix[chunk_y][chunk_x-1].Tiles[tile_y][31]
								level.Level_In_Matrix[chunk_y][chunk_x-1].Changed = true
							}
						} else {
							tile_left = level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x-1]
						}

						if tile_x == 31 {
							if chunk_x == len(level.Level_In_Matrix[chunk_y])-1 {
								tile_right = 0
							} else {
								tile_right = level.Level_In_Matrix[chunk_y][chunk_x+1].Tiles[tile_y][0]
								level.Level_In_Matrix[chunk_y][chunk_x+1].Changed = true
							}
						} else {
							tile_right = level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x+1]
						}

						if tile_above < 0 {
							tile_above = 0
						}
						if tile_down < 0 {
							tile_down = 0
						}
						if tile_left < 0 {
							tile_left = 0
						}
						if tile_right < 0 {
							tile_right = 0
						}

						if tile_above == 0 && tile_down != 0 {
							if tile_left == 0 && tile_right != 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 1
							}
							if tile_left != 0 && tile_right != 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 2
							}
							if tile_left != 0 && tile_right == 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 3
							}
						}

						if tile_above != 0 && tile_down != 0 {
							if tile_left == 0 && tile_right != 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 4
							}
							if tile_left != 0 && tile_right != 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 5
							}
							if tile_left != 0 && tile_right == 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 6
							}
						}

						if tile_above != 0 && tile_down == 0 {
							if tile_left == 0 && tile_right != 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 7
							}
							if tile_left != 0 && tile_right != 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 8
							}
							if tile_left != 0 && tile_right == 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 9
							}
						}

						if tile_right == 0 && tile_left == 0 {
							if tile_above == 0 && tile_down != 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 10
							}
							if tile_above != 0 && tile_down != 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 11
							}
							if tile_above != 0 && tile_down == 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 12
							}
						}

						if tile_above == 0 && tile_down == 0 {
							if tile_left == 0 && tile_right != 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 13
							}
							if tile_left != 0 && tile_right != 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 14
							}
							if tile_left != 0 && tile_right == 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 15
							}
						}
						if tile_above == 0 && tile_down == 0 {
							if tile_left == 0 && tile_right == 0 {
								level.Level_In_Matrix[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 16
							}
						}
					}
					chunk.Changed = true

					if chunk_x-1 >= 0 {
						level.Level_In_Matrix[chunk_y][chunk_x-1].Changed = true
					}

					if chunk_y-1 >= 0 {
						level.Level_In_Matrix[chunk_y-1][chunk_x].Changed = true
					}

					if chunk_y+1 < len(level.Level_In_Matrix) {
						level.Level_In_Matrix[chunk_y+1][chunk_x].Changed = true
					}

					if chunk_x+1 < len(level.Level_In_Matrix) {
						level.Level_In_Matrix[chunk_y][chunk_x+1].Changed = true
					}
				}
			}
		}
	}
}
