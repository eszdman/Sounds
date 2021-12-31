package forms

import (
	"encoding/gob"
	"github.com/eszdman/Sounds/engine/generator"
	"github.com/eszdman/Sounds/env"
	"github.com/eszdman/Sounds/setting"
	"github.com/inkyblackness/imgui-go/v4"
	"golang.org/x/image/colornames"
	"os"
)

func pianoDrawer(input func(i int)) {
	clipper := imgui.ListClipper{ItemsCount: int(setting.CurrentParameters.PianoCount), ItemsHeight: float32(setting.CurrentParameters.PianoY)}
	for clipper.Step() {
		for i := clipper.DisplayStart; i < clipper.DisplayEnd; i++ {
			input(i)
		}
	}
}
func mouseToRollCoords(mouseRelative imgui.Vec2) (x, y int) {
	y = int(mouseRelative.Y / float32(setting.CurrentParameters.PianoY))
	x = int(mouseRelative.X / float32(setting.CurrentParameters.Roll))
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
	octaves := int(setting.CurrentParameters.PianoOctaves) - i/12
	out = ""
	notes := (i) % 12
	out += note(notes)
	out += string(rune('0' + octaves - 1))
	return
}

var clickedRoll = 0
var clickedPianoY = 0
var noteDC = false

func childRoll(size imgui.Vec2) {
	//imgui.PushStyleVarFloat(imgui.StyleVarScrollbarSize, float32(setting.CurrentParameters.PianoY))
	//defer imgui.PopStyleVar()
	scrollY := imgui.ScrollY()
	if !env.ImguiWrapping.IO.KeyShiftPressed() {
		_, w := env.ImguiWrapping.IO.MouseWheel()
		imgui.SetScrollY(scrollY + w*-float32(setting.CurrentParameters.PianoY*2))
	}
	defer imgui.EndChild()

	if !imgui.BeginChildV("PianoRoll2", size, true,
		imgui.WindowFlagsNoScrollWithMouse|
			imgui.WindowFlagsNoScrollbar|
			imgui.WindowFlagsAlwaysHorizontalScrollbar) {
		return
	}
	scrollX := imgui.ScrollX()
	if env.ImguiWrapping.IO.KeyShiftPressed() {
		_, w := env.ImguiWrapping.IO.MouseWheel()
		imgui.SetScrollX(scrollX + w*-float32(setting.CurrentParameters.RollMpy*setting.CurrentParameters.Roll))
	}
	if imgui.IsKeyDown(341) {
		_, w := env.ImguiWrapping.IO.MouseWheel()
		setting.CurrentParameters.Roll += int32(w)
		if setting.CurrentParameters.Roll < 4 {
			setting.CurrentParameters.Roll = 4
		}
	}

	cursor := imgui.CursorScreenPos()
	mouse := imgui.MousePos()
	drawList := imgui.WindowDrawList()
	color1 := imgui.Packed(colornames.Darkgray)
	color2 := imgui.Packed(colornames.Gray)
	lines := imgui.Packed(colornames.Dimgray)
	var selectedLineColor imgui.PackedColor
	start := cursor.Plus(imgui.Vec2{X: imgui.ScrollX(), Y: scrollY})
	end := cursor.Plus(imgui.WindowSize().Plus(imgui.Vec2{X: imgui.ScrollX(), Y: scrollY}))
	drawLines := func(i int) {
		i2 := (i + 5) % 12
		if i2 != 1 && i2 != 3 && i2 != 6 && i2 != 8 && i2 != 10 {
			selectedLineColor = color1
		} else {
			selectedLineColor = color2
		}
		drawList.AddRectFilledV(imgui.Vec2{X: start.X, Y: cursor.Y + float32(i)*float32(setting.CurrentParameters.PianoY)},
			imgui.Vec2{X: end.X, Y: cursor.Y + float32(i+1)*float32(setting.CurrentParameters.PianoY)}, selectedLineColor, 0, imgui.DrawFlagsNone)
		if i2 == 5 || i2 == 0 {
			drawList.AddLineV(imgui.Vec2{X: start.X, Y: cursor.Y + float32(i)*float32(setting.CurrentParameters.PianoY)},
				imgui.Vec2{X: end.X, Y: cursor.Y + float32(i)*float32(setting.CurrentParameters.PianoY)}, lines, 1.5)
		}
	}

	drawVerticals := func(w int) {
		thickness := float32(1.0)
		if w%int(setting.CurrentParameters.RollSubMpy) == 0 {
			thickness = 2
		}
		drawList.AddLineV(
			imgui.Vec2{X: cursor.X + float32(w*int(setting.CurrentParameters.Roll*setting.CurrentParameters.RollMpy)), Y: start.Y},
			imgui.Vec2{X: cursor.X + float32(w*int(setting.CurrentParameters.Roll*setting.CurrentParameters.RollMpy)), Y: end.Y}, lines, thickness)
	}
	pianoDrawer(drawLines)
	env.ClipHorizontal(float32(setting.CurrentParameters.Roll*4), drawVerticals)

	xMouse, yMouse := mouseToRollCoords(mouse.Minus(cursor))
	for i := 0; i < len(env.VoiceNotes); i++ {
		callNote(&env.VoiceNotes[i], drawList, cursor, i)
	}
	imgui.SetScrollY(scrollY)
	imgui.Dummy(imgui.Vec2{X: float32(10 * 200)})

	if imgui.IsMouseDoubleClicked(0) && imgui.IsWindowFocused() && !noteDC {
		note := generator.FillNote()
		note.RollStart = xMouse
		note.RollEnd = xMouse + 3
		note.RollPitch = yMouse
		note.VibratoParams.Frequency = 5
		env.VoiceNotes = append(env.VoiceNotes, note)
	}
	noteDC = false

}

var isPinned = false

func PianoRoll() {
	//imgui.SetNextWindowSize(imgui.Vec2{1000,1000})
	imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{})
	defer imgui.PopStyleVar()
	flags := imgui.WindowFlagsNoScrollbar | imgui.WindowFlagsAlwaysVerticalScrollbar | imgui.WindowFlagsMenuBar | imgui.WindowFlagsNoScrollWithMouse
	if isPinned {
		flags |= imgui.WindowFlagsNoMove | imgui.WindowFlagsNoResize
	}
	if !imgui.BeginV("PianoRoll", nil, flags) {
		imgui.End()
		return
	}
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.MenuItem("Open") {
				go func() {
					file, _ := os.Open("pianoRoll.roll")
					encoder := gob.NewDecoder(file)
					_ = encoder.Decode(&env.VoiceNotes)
					_ = file.Close()
				}()
			}
			if imgui.MenuItem("Save") {
				go func() {
					file, _ := os.Create("pianoRoll.roll")
					encoder := gob.NewEncoder(file)
					_ = encoder.Encode(env.VoiceNotes)
					_ = file.Close()
				}()
			}
			imgui.EndMenu()
		}
		imgui.Checkbox("Pinned", &isPinned)

		imgui.EndMenuBar()
	}

	cursor := imgui.CursorScreenPos()
	cursorWin := imgui.CursorPos()
	drawList := imgui.WindowDrawList()
	drawList.AddRectFilledV(
		imgui.Vec2{X: cursor.X, Y: cursor.Y},
		imgui.Vec2{X: cursor.X + float32(setting.CurrentParameters.PianoX2),
			Y: cursor.Y + float32(setting.CurrentParameters.PianoY)*float32(setting.CurrentParameters.PianoCount+1)},
		imgui.Packed(colornames.White), 0, imgui.DrawFlagsNone)
	//sizeY := sizeY2 + sizeY2/2
	mouse := imgui.MousePos()
	lines := imgui.Packed(colornames.Dimgray)
	fontY := imgui.FontSize()
	drawButtons := func(i int) {
		i2 := (i + 5) % 12
		if i2 != 1 && i2 != 3 && i2 != 6 && i2 != 8 && i2 != 10 {
		} else {
			drawList.AddRectFilledV(
				imgui.Vec2{
					X: cursor.X,
					Y: cursor.Y + float32(i)*float32(setting.CurrentParameters.PianoY)},
				imgui.Vec2{
					X: cursor.X + float32(setting.CurrentParameters.PianoX),
					Y: cursor.Y + float32(i+1)*float32(setting.CurrentParameters.PianoY)},
				imgui.Packed(colornames.Black), 6, imgui.DrawFlagsRoundCornersRight)
		}
		if i2 == 5 || i2 == 0 {
			drawList.AddLineV(imgui.Vec2{X: cursor.X, Y: cursor.Y + float32(i)*float32(setting.CurrentParameters.PianoY)},
				imgui.Vec2{X: cursor.X + float32(setting.CurrentParameters.PianoX2), Y: cursor.Y + float32(i)*float32(setting.CurrentParameters.PianoY)}, lines, 1.5)
		}
		if i2 == 4 {
			drawList.AddText(
				imgui.Vec2{
					X: cursor.X + float32(setting.CurrentParameters.PianoX),
					Y: cursor.Y + float32(i)*float32(setting.CurrentParameters.PianoY) + fontY},
				imgui.Packed(colornames.Black), pianoNames(i))
		}
	}
	pianoDrawer(drawButtons)
	_, yMouse := mouseToRollCoords(mouse.Minus(cursor))

	drawList.AddText(
		imgui.Vec2{
			X: cursor.X + 10,
			Y: cursor.Y + float32(yMouse)*float32(setting.CurrentParameters.PianoY) + fontY},
		imgui.Packed(colornames.Gray), pianoNames(yMouse))

	childCursor := imgui.Vec2{
		X: cursorWin.X + float32(setting.CurrentParameters.PianoX2),
		Y: cursorWin.Y + imgui.ScrollY(),
	}
	dummySize := imgui.Vec2{
		X: 0,
		Y: float32(setting.CurrentParameters.PianoCount)*float32(setting.CurrentParameters.PianoY) + cursorWin.Y,
	}
	maxRegions := imgui.WindowContentRegionMax()
	rollSize := imgui.Vec2{
		X: imgui.WindowWidth() - float32(setting.CurrentParameters.PianoX2),
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
