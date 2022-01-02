package forms

import (
	"github.com/eszdman/Sounds/engine/preview"
	"github.com/eszdman/Sounds/env"
	"github.com/eszdman/Sounds/setting"
	"github.com/eszdman/Sounds/ui/widgets"
	"github.com/inkyblackness/imgui-go/v4"
	"golang.org/x/image/colornames"
	"image/color"
	"strconv"
)

func callNote(voiceNote *env.VoiceNote, drawList imgui.DrawList, cursor imgui.Vec2, i int) {
	startTile := imgui.Vec2{
		X: cursor.X + float32(setting.CurrentParameters.Roll)*float32(voiceNote.RollStart),
		Y: cursor.Y + float32(setting.CurrentParameters.PianoY)*float32(voiceNote.RollPitch),
	}
	endTile := imgui.Vec2{
		X: cursor.X + float32(setting.CurrentParameters.Roll)*float32(voiceNote.RollEnd),
		Y: cursor.Y + float32(setting.CurrentParameters.PianoY)*float32(voiceNote.RollPitch+1),
	}
	start := cursor.Plus(imgui.Vec2{X: imgui.ScrollX(), Y: imgui.ScrollY()})
	end := cursor.Plus(imgui.WindowSize().Plus(imgui.Vec2{X: imgui.ScrollX(), Y: imgui.ScrollY()}))
	if env.IsFramedV(startTile, endTile, start, end, cursor) {
		cursorN := startTile.Minus(cursor)

		name := strconv.Itoa(i)
		end := endTile.Minus(imgui.Vec2{X: float32(setting.CurrentParameters.Roll)})
		drawList.AddRectFilledV(startTile, end, imgui.Packed(colornames.White), 6, imgui.DrawFlagsRoundCornersLeft)
		widgets.PlotX(drawList, imgui.Packed(color.RGBA{R: 81, G: 81, B: 163, A: 0xff}), startTile, cursor.X+float32(voiceNote.RollEnd*int(setting.CurrentParameters.Roll)),
			8, voiceNote.VibratoParams.GetVibrato())
		//color := imgui.Packed(color.RGBA{R: 81, G: 81, B: 163, A: 0xff})
		color := imgui.Packed(colornames.Black)
		drawList.AddRectFilledV(end.Minus(imgui.Vec2{Y: float32(setting.CurrentParameters.PianoY)}), endTile, color, 0, imgui.DrawFlagsNone)
		drawList.AddText(endTile, imgui.Packed(colornames.White), name)
		noteSize := imgui.Vec2{
			X: float32(setting.CurrentParameters.Roll * int32(voiceNote.RollEnd-voiceNote.RollStart-1)),
			Y: float32(setting.CurrentParameters.PianoY),
		}
		imgui.PushID("note" + name)
		defer imgui.PopID()
		imgui.SetCursorPos(cursorN)
		imgui.InvisibleButton("noteMain", imgui.Vec2{
			X: float32(setting.CurrentParameters.Roll * int32(voiceNote.RollEnd-voiceNote.RollStart-1)),
			Y: float32(setting.CurrentParameters.PianoY),
		})

		xMouse, yMouse := mouseToRollCoords(imgui.MousePos().Minus(cursor))
		if imgui.IsItemClicked() {
			clickedRoll = xMouse - voiceNote.RollStart
			clickedPianoY = yMouse - voiceNote.RollPitch
		}
		focusText := false
		if imgui.IsMouseDoubleClicked(0) && imgui.IsItemActive() {
			noteDC = true
			focusText = true
		}
		if imgui.IsItemActive() && imgui.IsMouseDragging(0, 0) {

			xMouse -= clickedRoll + voiceNote.RollStart
			yMouse -= clickedPianoY + voiceNote.RollPitch
			voiceNote.RollPitch += yMouse
			if voiceNote.RollPitch < 0 {
				voiceNote.RollPitch = 0
			}
			if voiceNote.RollStart+xMouse >= 0 {
				voiceNote.RollStart += xMouse
				voiceNote.RollEnd += xMouse
			}
			preview.OnMouseDrag(voiceNote.GetFrequency())
		}
		if imgui.IsMouseReleased(0) {
			preview.OnMouseReleased()
		}
		imgui.SetCursorPos(cursorN.Plus(imgui.Vec2{Y: imgui.TextLineHeight()/2.0 + float32(setting.CurrentParameters.PianoY)}))
		imgui.PushItemWidth(float32(setting.CurrentParameters.Roll * int32(voiceNote.RollEnd-voiceNote.RollStart-1)))
		if focusText {
			imgui.SetKeyboardFocusHere()
		}
		imgui.InputTextV("", &voiceNote.Lyrics, imgui.InputTextFlagsCharsNoBlank, nil)
		imgui.IsItemActive()
		imgui.PopItemWidth()
		if imgui.IsMouseDoubleClicked(0) && imgui.IsItemActive() {
			noteDC = true
		}

		imgui.SetCursorPos(cursorN.Plus(imgui.Vec2{X: noteSize.X}))
		imgui.InvisibleButton("nodeEnd", imgui.Vec2{
			X: float32(setting.CurrentParameters.Roll),
			Y: float32(setting.CurrentParameters.PianoY),
		})
		if imgui.IsItemClicked() {
			clickedRoll = xMouse - voiceNote.RollEnd
		}
		if imgui.IsItemHovered() {
			imgui.SetMouseCursor(imgui.MouseCursorResizeEW)
		}
		if imgui.IsMouseDoubleClicked(0) && imgui.IsItemActive() {
			noteDC = true
		}
		if imgui.IsItemActive() && imgui.IsMouseDragging(0, 0) {
			imgui.SetMouseCursor(imgui.MouseCursorResizeEW)

			xMouse -= clickedRoll + voiceNote.RollEnd
			voiceNote.RollEnd += xMouse
			if voiceNote.RollEnd-voiceNote.RollStart < 2 {
				voiceNote.RollEnd = voiceNote.RollStart + 2
			}
		}
	}

}
