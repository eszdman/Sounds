package forms

import (
	"github.com/eszdman/Sounds/env"
	"github.com/inkyblackness/imgui-go/v4"
	"golang.org/x/image/colornames"
	"image/color"
	"strconv"
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

var clickedRoll = 0
var clickedPianoY = 0

func callNote(voiceNote *env.VoiceNote, drawList imgui.DrawList, cursor imgui.Vec2, i int) {
	end := imgui.WindowSize().Plus(imgui.Vec2{X: imgui.ScrollX(), Y: imgui.ScrollY()})
	startTile := imgui.Vec2{
		X: cursor.X + float32(currentParameters.Roll)*float32(voiceNote.RollStart),
		Y: cursor.Y + float32(currentParameters.PianoY)*float32(voiceNote.RollPitch),
	}
	endTile := imgui.Vec2{
		X: cursor.X + float32(currentParameters.Roll)*float32(voiceNote.RollEnd),
		Y: cursor.Y + float32(currentParameters.PianoY)*float32(voiceNote.RollPitch+1),
	}
	start := cursor.Plus(imgui.Vec2{X: imgui.ScrollX(), Y: imgui.ScrollY()})
	if endTile.X > start.X && startTile.X < start.X+end.X && endTile.Y > start.Y && startTile.Y < start.Y+end.Y {
		cursorN := startTile.Minus(cursor)

		name := strconv.Itoa(i)
		end := endTile.Minus(imgui.Vec2{X: float32(currentParameters.Roll)})
		drawList.AddRectFilledV(startTile, end, imgui.Packed(colornames.White), 6, imgui.DrawFlagsRoundCornersLeft)
		drawList.AddRectFilledV(end.Minus(imgui.Vec2{Y: float32(currentParameters.PianoY)}), endTile, imgui.Packed(color.RGBA{81, 81, 163, 0xff}), 0, imgui.DrawFlagsNone)
		noteSize := imgui.Vec2{
			X: float32(currentParameters.Roll * int32(voiceNote.RollEnd-voiceNote.RollStart-1)),
			Y: float32(currentParameters.PianoY),
		}
		imgui.PushID("note" + name)
		defer imgui.PopID()
		imgui.SetCursorPos(cursorN.Plus(imgui.Vec2{Y: imgui.TextLineHeight()/2.0 + float32(currentParameters.PianoY)}))
		imgui.PushItemWidth(float32(currentParameters.Roll * int32(voiceNote.RollEnd-voiceNote.RollStart-1)))
		imgui.InputTextV("", &voiceNote.Lyrics, imgui.InputTextFlagsCharsNoBlank, nil)
		imgui.IsItemActive()
		imgui.PopItemWidth()
		if imgui.IsMouseDoubleClicked(0) && imgui.IsItemActive() {

		}
		imgui.SetCursorPos(cursorN)
		imgui.InvisibleButton("noteMain", imgui.Vec2{
			X: float32(currentParameters.Roll * int32(voiceNote.RollEnd-voiceNote.RollStart-1)),
			Y: float32(currentParameters.PianoY),
		})

		xMouse, yMouse := mouseToRollCoords(imgui.MousePos().Minus(cursor))
		if imgui.IsItemClicked() {
			clickedRoll = xMouse - voiceNote.RollStart
			clickedPianoY = yMouse - voiceNote.RollPitch
		}
		if imgui.IsItemActive() && imgui.IsMouseDragging(0, 0) {
			xMouse -= clickedRoll + voiceNote.RollStart
			yMouse -= clickedPianoY + voiceNote.RollPitch
			voiceNote.RollPitch += yMouse
			voiceNote.RollStart += xMouse
			voiceNote.RollEnd += xMouse
		}

		imgui.SetCursorPos(cursorN.Plus(imgui.Vec2{X: noteSize.X}))
		imgui.InvisibleButton("nodeEnd", imgui.Vec2{
			X: float32(currentParameters.Roll),
			Y: float32(currentParameters.PianoY),
		})
		if imgui.IsItemClicked() {
			clickedRoll = xMouse - voiceNote.RollEnd
		}
		if imgui.IsItemActive() && imgui.IsMouseDragging(0, 0) {
			xMouse -= clickedRoll + voiceNote.RollEnd
			voiceNote.RollEnd += xMouse
			if voiceNote.RollEnd-voiceNote.RollStart < 2 {
				voiceNote.RollEnd = voiceNote.RollStart + 2
			}
		}
	}

}
func childRoll(size imgui.Vec2) {
	scrollY := imgui.ScrollY()
	defer imgui.EndChild()
	if !imgui.BeginChildV("PianoRoll2", size, true,
		imgui.WindowFlagsNoScrollWithMouse|
			imgui.WindowFlagsNoScrollbar|
			imgui.WindowFlagsAlwaysHorizontalScrollbar) {
		return
	}
	cursor := imgui.CursorScreenPos()
	mouse := imgui.MousePos()
	drawList := imgui.WindowDrawList()
	color1 := imgui.Packed(colornames.Darkgray)
	color2 := imgui.Packed(colornames.Gray)
	var color imgui.PackedColor
	drawLines := func(i int) {
		i2 := i
		i2 %= 12
		if i2 != 1 && i2 != 3 && i2 != 6 && i2 != 8 && i2 != 10 {
			color = color1
		} else {
			color = color2
		}
		drawList.AddRectFilled(imgui.Vec2{X: cursor.X, Y: cursor.Y + float32(i)*float32(currentParameters.PianoY)},
			imgui.Vec2{X: cursor.X + float32(10*200), Y: cursor.Y + float32(i+1)*float32(currentParameters.PianoY)}, color)
	}
	pianoDrawer(drawLines)

	xMouse, yMouse := mouseToRollCoords(mouse.Minus(cursor))
	//fontY := imgui.FontSize()
	rollPosition := imgui.Vec2{
		X: cursor.X + float32(xMouse)*float32(currentParameters.Roll),
		Y: cursor.Y + float32(yMouse)*float32(currentParameters.PianoY)}
	rect := imgui.Vec2{float32(currentParameters.Roll), float32(currentParameters.PianoY)}
	for i := 0; i < len(env.VoiceNotes); i++ {
		callNote(&env.VoiceNotes[i], drawList, cursor, i)
	}
	imgui.SetScrollY(scrollY)
	drawList.AddRectFilled(rollPosition, rollPosition.Plus(rect), imgui.Packed(colornames.Purple))
	imgui.Dummy(imgui.Vec2{float32(10 * 200), 0})

	if imgui.IsMouseDoubleClicked(0) && imgui.IsWindowFocused() {
		env.VoiceNotes = append(env.VoiceNotes, env.VoiceNote{RollStart: xMouse, RollEnd: xMouse + 3, RollPitch: yMouse})
	}

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
