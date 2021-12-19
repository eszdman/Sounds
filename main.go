package main

import (
	"fmt"
	"github.com/eszdman/Sounds/gui"
	"github.com/eszdman/Sounds/platform"
	"github.com/inkyblackness/imgui-go/v4"
	"golang.org/x/image/colornames"
	"strconv"
	//"google.golang.org/grpc/balancer/grpclb/state"
	"os"
)

const windowWidth = 1600
const windowHeight = 900

var newPlatform *platform.Platform

func Show(keep *bool) {
	// Use fixed width for labels (by passing a negative value), the rest goes to widgets.
	// We choose a width proportional to our font size.
	size := newPlatform.DisplaySize()
	imgui.SetNextWindowPosV(imgui.Vec2{X: 0, Y: 0}, imgui.ConditionFirstUseEver, imgui.Vec2{})
	imgui.SetNextWindowSize(imgui.Vec2{X: size[0], Y: size[1]})
	//imgui.StyleVarWindowBorderSize

	if !imgui.BeginV("Window", keep, imgui.WindowFlagsNoMove|
		imgui.WindowFlagsMenuBar|imgui.WindowFlagsNoTitleBar|
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoBringToFrontOnFocus|
		imgui.WindowFlagsNoBackground) {
		// Early out if the window is collapsed, as an optimization.
		imgui.End()
		return
	}
	imgui.PushItemWidth(imgui.FontSize() * -12)
	// MenuBar
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.MenuItem("New") {
			}
			if imgui.MenuItemV("Open", "Ctrl+O", false, true) {
				println("Opened!")
			}
			if imgui.MenuItemV("Save", "Ctrl+S", false, true) {
				println("Saved!")
			}
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Examples") {
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Tools") {
			imgui.EndMenu()
		}
		imgui.EndMenuBar()
	}
	// End of ShowDemoWindow()
	imgui.End()
}
func abs(float1 float32) float32 {
	if float1 < 0 {
		float1 = -float1
	}
	return float1
}

var stringArr = make([]string, 100000)

func clipperTest() {
	if !imgui.BeginV("ClipperTest", nil, imgui.WindowFlagsAlwaysHorizontalScrollbar|imgui.WindowFlagsAlwaysVerticalScrollbar) {
		imgui.End()
		return
	}
	clipper := imgui.ListClipper{ItemsCount: len(stringArr)}
	for clipper.Step() {
		for i := clipper.DisplayStart; i < clipper.DisplayEnd; i++ {
			imgui.Text(stringArr[i])
		}
	}
	imgui.End()
}
func pianoRoll() {
	//imgui.SetNextWindowSize(imgui.Vec2{1000,1000})
	if !imgui.BeginV("PianoRoll", nil, imgui.WindowFlagsAlwaysHorizontalScrollbar|imgui.WindowFlagsAlwaysVerticalScrollbar) {
		imgui.End()
		return
	}

	cursor := imgui.CursorScreenPos()
	drawList := imgui.WindowDrawList()
	drawList.AddRectFilled(imgui.Vec2{cursor.X, cursor.Y}, imgui.Vec2{cursor.X + 1000, cursor.X + 1000}, imgui.Packed(colornames.Gray))
	cnt := 3
	sizeX := float32(150)
	sizeX2 := float32(100)
	sizeY2 := float32(30)
	sizeY := sizeY2 + sizeY2/2
	mouseY := imgui.MousePos().Y
	//dist2 := float32(0.0)
	//cond := false
	yCoord := float32(cursor.Y)
	for i := 0; i < 20; i++ {
		addStep := float32(0)
		//dist1 := abs(yCoord-mouseY)
		drawList.AddRectFilled(imgui.Vec2{cursor.X, yCoord + 1}, imgui.Vec2{cursor.X + sizeX, yCoord + sizeY - 1}, imgui.Packed(colornames.White))
		cnt %= 7
		if cnt != 2 && cnt != 3 && cnt != 5 && cnt != 6 {
			drawList.AddRectFilled(imgui.Vec2{cursor.X, yCoord + sizeY - 1}, imgui.Vec2{cursor.X + sizeX, yCoord + sizeY + sizeY2/2.0 - 1}, imgui.Packed(colornames.White))
			addStep = sizeY2 / 2.0
		}
		if cnt == 0 || cnt == 1 || cnt == 2 || cnt == 4 || cnt == 5 {
			drawList.AddRectFilled(imgui.Vec2{cursor.X, yCoord - sizeY2/2.0}, imgui.Vec2{cursor.X + sizeX2, yCoord + sizeY2/2.0}, imgui.Packed(colornames.Black))
			//dist2 = abs(yCoord-mouseY-sizeY2/2.0)
		}
		yCoord += sizeY + addStep
		cnt++
	}
	relative := mouseY - cursor.Y
	drawList.AddText(imgui.Vec2{cursor.X, cursor.Y + relative - float32(int(relative)%int(sizeY2))}, imgui.Packed(colornames.Purple), "Test")
	imgui.Dummy(imgui.Vec2{360, 360})
	imgui.End()
}
func tableTest() {
	imgui.BeginTable("test", 10)
	for row := 0; row < 4; row++ {
		imgui.TableNextRow()
		imgui.TableNextColumn()
		imgui.Text("Row " + strconv.Itoa(row))
		imgui.TableNextColumn()
		imgui.Text("Some contents")
		imgui.TableNextColumn()
		imgui.Text("123.456")
	}
	imgui.EndTable()
}
func main() {
	for i := 0; i < len(stringArr); i++ {
		stringArr[i] = strconv.Itoa(i)
	}
	var err error
	newPlatform, err = platform.NewPlatform(windowWidth, windowHeight, "SoundEngine")
	context := imgui.CreateContext(nil)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	keep := true
	imguiRenderer, err := gui.NewOpenGL3("#version 430")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer newPlatform.Dispose()
	defer context.Destroy()
	defer imguiRenderer.Dispose()
	imguiWrapping := gui.NewImgui(newPlatform, imguiRenderer, context)
	imgui.StyleColorsClassic()
	mainRunner := func() {
		newPlatform.ProcessEvents()
		imguiWrapping.NewFrame()
		{
			imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 0)
			Show(&keep)
			imgui.PopStyleVar()
			imgui.ShowDemoWindow(&keep)
			pianoRoll()
			clipperTest()
		}
	}
	imguiWrapping.Run(mainRunner)
}
