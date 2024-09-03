package visual

import rl "github.com/gen2brain/raylib-go/raylib"

type UI struct {
	Width  int32
	Height int32
	Title  string
}

func (u *UI) Init() {
	rl.InitWindow(u.Width, u.Height, u.Title)
	rl.SetTargetFPS(60)
}

func (u *UI) Run() {
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
		rl.EndDrawing()
	}
}

func (u *UI) Destroy() {
	rl.CloseWindow()
}

func NewUI() *UI {
	return &UI{
		Width:  640,
		Height: 480,
		Title:  "KO : Resampler",
	}
}
