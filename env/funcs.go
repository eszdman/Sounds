package env

import "github.com/inkyblackness/imgui-go/v4"

func IsFramed(startTile, endTile, cursor imgui.Vec2) bool {
	end := imgui.WindowSize().Plus(imgui.Vec2{X: imgui.ScrollX(), Y: imgui.ScrollY()})
	start := cursor.Plus(imgui.Vec2{X: imgui.ScrollX(), Y: imgui.ScrollY()})
	return IsFramedV(startTile, endTile, start, end, cursor)
}
func ClipHorizontal(size float32, drawing func(w int)) {
	end := imgui.WindowSize().Plus(imgui.Vec2{X: imgui.ScrollX(), Y: imgui.ScrollY()})
	start := imgui.Vec2{X: imgui.ScrollX(), Y: imgui.ScrollY()}
	for i := int(start.X / size); i < int(end.X/size); i++ {
		drawing(i)
	}
}
func IsFramedV(startTile, endTile, start, end, cursor imgui.Vec2) bool {
	return endTile.X > start.X && startTile.X < start.X+end.X && endTile.Y > start.Y && startTile.Y < start.Y+end.Y
}
