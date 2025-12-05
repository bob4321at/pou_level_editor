package grid

func GetSurroundingTiles(matrix *[][]Chunk, tile_x, tile_y, chunk_x, chunk_y int, level *Level) (int, int, int, int) {
	var tile_above int = 0
	var tile_left int = 0
	var tile_right int = 0
	var tile_down int = 0

	if tile_y == 0 {
		if chunk_y == 0 {
			tile_above = 0
		} else {
			tile_above = (*matrix)[chunk_y-1][chunk_x].Tiles[31][tile_x]
			(*matrix)[chunk_y-1][chunk_x].Changed = 1
		}
	} else {
		tile_above = (*matrix)[chunk_y][chunk_x].Tiles[tile_y-1][tile_x]
	}

	if tile_y == 31 {
		if chunk_y == len((*matrix))-1 {
			tile_down = 0
		} else {
			tile_down = (*matrix)[chunk_y+1][chunk_x].Tiles[0][tile_x]
			(*matrix)[chunk_y+1][chunk_x].Changed = 1
		}
	} else {
		tile_down = (*matrix)[chunk_y][chunk_x].Tiles[tile_y+1][tile_x]
	}

	if tile_x == 0 {
		if chunk_x == 0 {
			tile_left = 0
		} else {
			tile_left = (*matrix)[chunk_y][chunk_x-1].Tiles[tile_y][31]
			(*matrix)[chunk_y][chunk_x-1].Changed = 1
		}
	} else {
		tile_left = (*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x-1]
	}

	if tile_x == 31 {
		if chunk_x == len((*matrix)[chunk_y])-1 {
			tile_right = 0
		} else {
			tile_right = (*matrix)[chunk_y][chunk_x+1].Tiles[tile_y][0]
			(*matrix)[chunk_y][chunk_x+1].Changed = 1
		}
	} else {
		tile_right = (*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x+1]
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

	return tile_above, tile_left, tile_right, tile_down
}

func GetAboveTile(matrix *[][]Chunk, tile_x, tile_y, chunk_x, chunk_y int, level *Level) int {
	var tile_above int = 0

	if tile_y == 0 {
		if chunk_y == 0 {
			tile_above = 0
		} else {
			tile_above = (*matrix)[chunk_y-1][chunk_x].Tiles[31][tile_x]
		}
	} else {
		tile_above = (*matrix)[chunk_y][chunk_x].Tiles[tile_y-1][tile_x]
	}

	return tile_above
}

func (level *Level) NeighbourCheck(matrix *[][]Chunk) {
	for chunk_y, chunk_row := range *matrix {
		for chunk_x := range chunk_row {
			chunk := &(*matrix)[chunk_y][chunk_x]
			for tile_y, tile_row := range chunk.Tiles {
				for tile_x, tile := range tile_row {
					if tile < 0 {
						if tile == -8 || tile == -5 || tile == -9 || tile == -12 {
							if tile != -12 {
								tile_above, tile_left, tile_right, tile_down := GetSurroundingTiles(matrix, tile_x, tile_y, chunk_x, chunk_y, level)
								var Dir_To_Change *int

								for i := range level.SpikeTiles {
									spike := &level.SpikeTiles[i]
									tile_check_if_this_spike := &(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x]
									if spike.Tile == tile_check_if_this_spike {
										Dir_To_Change = &spike.Direction
									}
								}

								for i := range level.SpringTiles {
									spring := &level.SpringTiles[i]
									tile_check_if_this_spike := &(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x]
									if spring.Tile == tile_check_if_this_spike {
										Dir_To_Change = &spring.Direction
									}
								}

								for i := range level.TriggerTile {
									trigger := &level.TriggerTile[i]
									tile_check_if_this_spike := &(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x]
									if trigger.Tile == tile_check_if_this_spike {
										Dir_To_Change = &trigger.Direction
									}
								}

								if Dir_To_Change != nil {
									if tile_left > 0 {
										*Dir_To_Change = 1
									}
									if tile_right > 0 {
										*Dir_To_Change = 3
									}
									if tile_above > 0 {
										*Dir_To_Change = 2
									}
									if tile_down > 0 {
										*Dir_To_Change = 0
									}
								}
							} else {
								tile_above := GetAboveTile(matrix, tile_x, tile_y, chunk_x, chunk_y, level)
								if tile_above > 0 || tile_above == -12 {
									for i := range level.WaterTiles {
										water := &level.WaterTiles[i]
										tile := &(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x]
										if water.Tile == tile {
											water.Top_Bottom = true
											if tile_y != 0 {
												other_tile := &(*matrix)[chunk_y][chunk_x].Tiles[tile_y-1][tile_x]
												if *other_tile == -12 {
													for j := range level.WaterTiles {
														other_water_tile := &level.WaterTiles[j]
														if other_water_tile.Tile == other_tile {
															if other_water_tile.Dissapear_Or_Appear != water.Dissapear_Or_Appear {
																water.Dissapear_Or_Appear_Top_Or_Bottom = true
															} else {
																water.Dissapear_Or_Appear_Top_Or_Bottom = false
															}
														}
													}
												} else {
													water.Dissapear_Or_Appear_Top_Or_Bottom = false
												}
											} else {
												if chunk_y != 0 {
													other_tile := &(*matrix)[chunk_y-1][chunk_x].Tiles[31][tile_x]
													if *other_tile == -12 {
														for j := range level.WaterTiles {
															other_water_tile := &level.WaterTiles[j]
															if other_water_tile.Tile == other_tile {
																if other_water_tile.Dissapear_Or_Appear != water.Dissapear_Or_Appear {
																	water.Dissapear_Or_Appear_Top_Or_Bottom = true
																} else {
																	water.Dissapear_Or_Appear_Top_Or_Bottom = false
																}
															}
														}
													} else {
														water.Dissapear_Or_Appear_Top_Or_Bottom = false
													}
												} else {
													water.Dissapear_Or_Appear_Top_Or_Bottom = false
												}
											}
										}
									}
								} else {
									for i := range level.WaterTiles {
										water := &level.WaterTiles[i]
										tile := &(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x]
										if water.Tile == tile {
											water.Top_Bottom = false
										}
									}
								}
							}
						} else {
							tile = 0
						}
					}
					if tile != 0 && tile != -8 && tile != -5 && tile != -9 && tile != -12 {
						tile_above, tile_left, tile_right, tile_down := GetSurroundingTiles(matrix, tile_x, tile_y, chunk_x, chunk_y, level)

						if tile_above == 0 && tile_down != 0 {
							if tile_left == 0 && tile_right != 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 1
							}
							if tile_left != 0 && tile_right != 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 2
							}
							if tile_left != 0 && tile_right == 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 3
							}
						}

						if tile_above != 0 && tile_down != 0 {
							if tile_left == 0 && tile_right != 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 4
							}
							if tile_left != 0 && tile_right != 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 5
							}
							if tile_left != 0 && tile_right == 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 6
							}
						}

						if tile_above != 0 && tile_down == 0 {
							if tile_left == 0 && tile_right != 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 7
							}
							if tile_left != 0 && tile_right != 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 8
							}
							if tile_left != 0 && tile_right == 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 9
							}
						}

						if tile_right == 0 && tile_left == 0 {
							if tile_above == 0 && tile_down != 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 10
							}
							if tile_above != 0 && tile_down != 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 11
							}
							if tile_above != 0 && tile_down == 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 12
							}
						}

						if tile_above == 0 && tile_down == 0 {
							if tile_left == 0 && tile_right != 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 13
							}
							if tile_left != 0 && tile_right != 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 14
							}
							if tile_left != 0 && tile_right == 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 15
							}
						}
						if tile_above == 0 && tile_down == 0 {
							if tile_left == 0 && tile_right == 0 {
								(*matrix)[chunk_y][chunk_x].Tiles[tile_y][tile_x] = 16
							}
						}
					}
				}
			}
		}
	}
}
