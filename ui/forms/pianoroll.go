package forms

import (
	"github.com/inkyblackness/imgui-go/v4"
	"golang.org/x/image/colornames"
)

var currentParameters PianoRollParameters

func pianoDrawer(input func(i int)) {
	clipper := imgui.ListClipper{ItemsCount: int(currentParameters.PianoCount), ItemsHeight: float32(currentParameters.PianoY)}
	for clipper.Step() {
		for i := clipper.DisplayStart; i < clipper.DisplayEnd; i++ {
			input(i)
		}
	}
}
func mouseToRollCoords(mouseRelative imgui.Vec2) (x, y int) {
	y = int(mouseRelative.Y / float32(currentParameters.PianoY))
	x = int(mouseRelative.X / float32(currentParameters.Roll))
	return
}
func note(i int) (out string) {
	switch i {
	case 0:
		out = "B"
		return
	case 1:
		out = "AB"
		return
	case 2:
		out = "A"
		return
	case 3:
		out = "GA"
		return
	case 4:
		out = "G"
		return
	case 5:
		out = "FG"
		return
	case 6:
		out = "F"
		return
	case 7:
		out = "E"
		return
	case 8:
		out = "DE"
		return
	case 9:
		out = "D"
		return
	case 10:
		out = "CD"
		return
	case 11:
		out = "C"
		return
	}
	return
}
func pianoNames(i int) (out string) {
	octaves := int(currentParameters.PianoOctaves) - i/12
	out = ""
	notes := i % 12
	out += note(notes)
	out += string(rune('0' + octaves - 1))
	return
}
func childRoll(size imgui.Vec2) {
	scrollY := imgui.ScrollY()
	if !imgui.BeginChildV("PianoRoll2", size, true, imgui.WindowFlagsNoScrollWithMouse|imgui.WindowFlagsNoScrollbar|imgui.WindowFlagsAlwaysHorizontalScrollbar) {
		imgui.EndChild()
		return
	}
	cursor := imgui.CursorScreenPos()
	mouse := imgui.MousePos()
	drawList := imgui.WindowDrawList()
	color1 := imgui.Packed(colornames.Gray)
	color2 := imgui.Packed(colornames.Black)
	var color imgui.PackedColor
	drawLines := func(i int) {
		i2 := i
		i2 %= 12
		if i2 != 1 && i2 != 3 && i2 != 6 && i2 != 8 && i2 != 10 {
			color = color1
		} else {
			color = color2
		}
		drawList.AddRectFilled(imgui.Vec2{X: cursor.X, Y: cursor.Y + float32(i)*float32(currentParameters.PianoY)}, imgui.Vec2{cursor.X + float32(10*200), cursor.Y + float32(i+1)*float32(currentParameters.PianoY)}, color)
	}
	pianoDrawer(drawLines)

	xMouse, yMouse := mouseToRollCoords(mouse.Minus(cursor))
	//fontY := imgui.FontSize()
	rollPosition := imgui.Vec2{
		X: cursor.X + float32(xMouse)*float32(currentParameters.Roll),
		Y: cursor.Y + float32(yMouse)*float32(currentParameters.PianoY)}
	rect := imgui.Vec2{float32(currentParameters.Roll), float32(currentParameters.PianoY)}
	drawList.AddRectFilled(rollPosition, rollPosition.Plus(rect), imgui.Packed(colornames.Purple))
	//drawList.AddText(imgui.Vec2{
	//	X: cursor.X + float32(xMouse)*Roll,
	//	Y: cursor.Y + float32(yMouse)*PianoY + fontY/2}, imgui.Packed(colornames.Purple), "Test2")
	imgui.SetScrollY(scrollY)
	if imgui.IsMouseDoubleClicked(1) {

	}
	imgui.EndChild()
}
func PianoRoll() {
	//imgui.SetNextWindowSize(imgui.Vec2{1000,1000})
	imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{0, 0})
	defer imgui.PopStyleVar()
	if !imgui.BeginV("PianoRoll", nil, imgui.WindowFlagsNoScrollbar|imgui.WindowFlagsAlwaysVerticalScrollbar) {
		imgui.End()
		return
	}
	cursor := imgui.CursorScreenPos()
	cursorWin := imgui.CursorPos()
	drawList := imgui.WindowDrawList()
	drawList.AddRectFilled(
		imgui.Vec2{X: cursor.X, Y: cursor.Y},
		imgui.Vec2{X: cursor.X + float32(currentParameters.PianoX2),
			Y: cursor.Y + float32(currentParameters.PianoY)*float32(currentParameters.PianoCount+1)},
		imgui.Packed(colornames.White))
	//sizeY := sizeY2 + sizeY2/2
	mouse := imgui.MousePos()
	drawButtons := func(i int) {
		i2 := i
		i2 %= 12
		if i2 != 1 && i2 != 3 && i2 != 6 && i2 != 8 && i2 != 10 {
		} else {
			drawList.AddRectFilled(
				imgui.Vec2{
					X: cursor.X,
					Y: cursor.Y + float32(i)*float32(currentParameters.PianoY)},
				imgui.Vec2{
					X: cursor.X + float32(currentParameters.PianoX),
					Y: cursor.Y + float32(i+1)*float32(currentParameters.PianoY)},
				imgui.Packed(colornames.Black))
		}
	}
	pianoDrawer(drawButtons)
	_, yMouse := mouseToRollCoords(mouse.Minus(cursor))
	fontY := imgui.FontSize()
	drawList.AddText(
		imgui.Vec2{
			X: cursor.X,
			Y: cursor.Y + float32(yMouse)*float32(currentParameters.PianoY) + fontY/2},
		imgui.Packed(colornames.Purple), pianoNames(yMouse))

	childCursor := imgui.Vec2{
		X: cursorWin.X + float32(currentParameters.PianoX2),
		Y: cursorWin.Y + imgui.ScrollY(),
	}
	dummySize := imgui.Vec2{
		X: 0,
		Y: float32(currentParameters.PianoCount)*float32(currentParameters.PianoY) + cursorWin.Y,
	}
	maxRegions := imgui.WindowContentRegionMax()
	rollSize := imgui.Vec2{
		X: imgui.WindowWidth() - float32(currentParameters.PianoX2),
		Y: imgui.WindowHeight() - cursorWin.Y,
	}
	if childCursor.Y+rollSize.Y > dummySize.Y {
		rollSize.Y -= (childCursor.Y + rollSize.Y) - dummySize.Y
	}

	if childCursor.X+rollSize.X > maxRegions.X {
		rollSize.X -= (childCursor.X + rollSize.X) - (maxRegions.X)
	}
	imgui.SetCursorPos(childCursor)
	childRoll(rollSize)
	imgui.End()
}
