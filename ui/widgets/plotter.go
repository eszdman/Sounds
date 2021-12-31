package widgets

import "github.com/inkyblackness/imgui-go/v4"

func PlotX(drawList imgui.DrawList, drawColor imgui.PackedColor, start imgui.Vec2, end, kx float32, drawer func(x float32) float32) {
	prev := imgui.Vec2{}
	for i := 0; i < int((end-start.X)/kx); i++ {
		x := float32(i) * kx
		val := drawer(x)
		output := imgui.Vec2{
			X: x,
			Y: val,
		}
		drawList.AddLineV(start.Plus(prev), start.Plus(output), drawColor, 2)
		prev = output.Minus(imgui.Vec2{X: 1})
	}
}
