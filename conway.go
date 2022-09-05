package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Cell uint8

func (cell *Cell) IsLive() bool {
	return *cell == 255
}

func (cell *Cell) Die() {
	*cell = 0
}

func (cell *Cell) Decay() {
	if *cell > 0 {
		*cell -= 1
	}
}

func (cell *Cell) Spawn() {
	*cell = 255
}

type Arena struct {
	Cells []Cell
	Next  []Cell
	Size  int
}

func (arena *Arena) Pos(index int) (x, y int) {
	return index % arena.Size, index / arena.Size
}

func (arena *Arena) Index(x, y int) int {
	return y*arena.Size + x
}

func (arena *Arena) Neighbors(x, y int) (result int) {
	for j := -1; j <= 1; j += 1 {
		for i := -1; i <= 1; i += 1 {
			if i == 0 && j == 0 {
				continue
			}

			x2 := x + i
			y2 := y + j

			if x2 <= 0 || x2 >= arena.Size {
				continue
			}

			if y2 <= 0 || y2 >= arena.Size {
				continue
			}

			if arena.Cells[arena.Index(x2, y2)].IsLive() {
				result += 1
			}
		}
	}

	return
}

func (arena *Arena) Draw(pixels []byte) {
	for i, cell := range arena.Cells {
		if cell.IsLive() {
			pixels[4*i+0] = 0
			pixels[4*i+1] = 255
			pixels[4*i+2] = 0
			pixels[4*i+3] = 255
		} else {
			pixels[4*i+0] = 0
			pixels[4*i+1] = 255 / (uint8(cell) + 1)
			pixels[4*i+2] = 255 - uint8(cell)
			pixels[4*i+3] = uint8(cell)
		}
	}
}

func (arena *Arena) Tick() {
	for i, cell := range arena.Cells {
		count := arena.Neighbors(arena.Pos(i))
		next := &arena.Next[i]

		switch {
		case count < 2:
			next.Decay()

		case (count == 2 || count == 3) && cell.IsLive():
			next.Spawn()

		case count > 3:
			next.Decay()

		case count == 3:
			next.Spawn()
		}
	}

	for i, cell := range arena.Next {
		arena.Cells[i] = cell

		if cell.IsLive() {
			arena.Next[i].Decay()
		}
	}
}

type App struct {
	Arena  Arena
	Pixels []byte
}

func (app *App) Update() error {
	app.Arena.Tick()
	return nil
}

func (app *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	app.Arena.Draw(app.Pixels)
	screen.WritePixels(app.Pixels)
}

func (app *App) Layout(width, height int) (int, int) {
	return app.Arena.Size, app.Arena.Size
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ebiten.SetWindowSize(500, 500)
	ebiten.SetWindowTitle("Conway")

	size := 250

	app := App{
		Arena: Arena{
			Size:  size,
			Cells: make([]Cell, size*size),
			Next:  make([]Cell, size*size),
		},
		Pixels: make([]byte, size*size*4),
	}

	for i := 0; i < 10000; i += 1 {
		app.Arena.Cells[rand.Intn(len(app.Arena.Cells))].Spawn()
	}

	if err := ebiten.RunGame(&app); err != nil {
		log.Fatal(err)
	}
}
