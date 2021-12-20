package wrapper

import (
	"github.com/eszdman/Sounds/env"
	"github.com/eszdman/Sounds/ui/platform"
	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/inkyblackness/imgui-go/v4"
	"math"
	"time"
)

// ImguiWrapping is the state holder for the imgui framework
type ImguiWrapping struct {
	io               *imgui.IO
	time             float64
	window           *glfw.Window
	platform         *platform.Platform
	context          *imgui.Context
	renderer         *OpenGL3
	runner           func()
	mouseJustPressed [3]bool
}

var cursors []*glfw.Cursor

const (
	mouseButtonPrimary   = 0
	mouseButtonSecondary = 1
	mouseButtonTertiary  = 2
	mouseButtonCount     = 3
)

var glfwButtonIndexByID = map[glfw.MouseButton]int{
	glfw.MouseButton1: mouseButtonPrimary,
	glfw.MouseButton2: mouseButtonSecondary,
	glfw.MouseButton3: mouseButtonTertiary,
}

var glfwButtonIDByIndex = map[int]glfw.MouseButton{
	mouseButtonPrimary:   glfw.MouseButton1,
	mouseButtonSecondary: glfw.MouseButton2,
	mouseButtonTertiary:  glfw.MouseButton3,
}

// ImguiMouseState is provided to NewFrame(...), containing the mouse state
type ImguiMouseState struct {
	MousePosX  float32
	MousePosY  float32
	MousePress [3]bool
}

var imguiIO imgui.IO
var inputState ImguiWrapping

// NewImgui initializes a new imgui context and a input object
func NewImgui(platform *platform.Platform, renderer *OpenGL3, context *imgui.Context) *ImguiWrapping {
	imguiIO = imgui.CurrentIO()
	inputState = ImguiWrapping{
		io:       &imguiIO,
		window:   platform.Window,
		platform: platform,
		renderer: renderer,
		context:  context,
		time:     0}
	inputState.setKeyMapping()
	inputState.io.SetMouseDrawCursor(true)
	inputState.io.SetBackendFlags(imgui.BackendFlagsHasMouseCursors)
	cursors = make([]*glfw.Cursor, imgui.MouseCursorCount)
	cursors[imgui.MouseCursorArrow] = glfw.CreateStandardCursor(glfw.ArrowCursor)
	cursors[imgui.MouseCursorTextInput] = glfw.CreateStandardCursor(glfw.IBeamCursor)
	cursors[imgui.MouseCursorResizeAll] = glfw.CreateStandardCursor(glfw.ArrowCursor)
	cursors[imgui.MouseCursorResizeNS] = glfw.CreateStandardCursor(glfw.VResizeCursor)
	cursors[imgui.MouseCursorResizeEW] = glfw.CreateStandardCursor(glfw.HResizeCursor)
	cursors[imgui.MouseCursorResizeNESW] = glfw.CreateStandardCursor(glfw.ArrowCursor)
	cursors[imgui.MouseCursorResizeNWSE] = glfw.CreateStandardCursor(glfw.ArrowCursor)
	cursors[imgui.MouseCursorHand] = glfw.CreateStandardCursor(glfw.HandCursor)

	inputState.installCallbacks()
	return &inputState
}
func (input *ImguiWrapping) installCallbacks() {
	input.window.SetMouseButtonCallback(input.mouseButtonChange)
	input.window.SetScrollCallback(input.mouseScrollChange)
	input.window.SetKeyCallback(input.keyChange)
	input.window.SetCharCallback(input.charChange)
	input.window.SetSizeCallback(input.sizeChange)
}
func (input *ImguiWrapping) updateMouseCursor() {
	//imgui_cursor := imgui.MouseCursor()
	if input.window.GetInputMode(glfw.CursorMode) == glfw.CursorDisabled {
		return
	}
	//input.window.SetCursor(cursors[imgui_cursor])
	input.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
}
func (input *ImguiWrapping) NewFrame() {
	cursorX, cursorY := input.platform.GetCursorPos()
	mouseState := ImguiMouseState{
		MousePosX:  float32(cursorX),
		MousePosY:  float32(cursorY),
		MousePress: input.platform.GetMousePresses123(),
	}
	// Setup display size (every frame to accommodate for window resizing)
	sizes := input.platform.DisplaySize()
	for i := 0; i < len(sizes); i++ {
		if sizes[i] <= 500 {
			sizes[i] = 500
		}
	}
	input.io.SetDisplaySize(imgui.Vec2{X: input.platform.DisplaySize()[0], Y: input.platform.DisplaySize()[1]})

	// Setup time step
	currentTime := glfw.GetTime()
	if input.time > 0 {
		input.io.SetDeltaTime(float32(currentTime - input.time))
	}
	input.time = currentTime

	// Setup inputs
	if input.platform.IsFocused() {
		input.io.SetMousePosition(imgui.Vec2{X: mouseState.MousePosX, Y: mouseState.MousePosY})
	} else {
		input.io.SetMousePosition(imgui.Vec2{X: -math.MaxFloat32, Y: -math.MaxFloat32})
	}
	for i := 0; i < len(input.mouseJustPressed); i++ {
		down := input.mouseJustPressed[i] || (input.window.GetMouseButton(glfwButtonIDByIndex[i]) == glfw.Press)
		input.io.SetMouseButtonDown(i, down)
		input.mouseJustPressed[i] = false
	}
	input.updateMouseCursor()

	imgui.NewFrame()
}
func (input *ImguiWrapping) setKeyMapping() {
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	input.io.KeyMap(imgui.KeyTab, int(glfw.KeyTab))
	input.io.KeyMap(imgui.KeyLeftArrow, int(glfw.KeyLeft))
	input.io.KeyMap(imgui.KeyRightArrow, int(glfw.KeyRight))
	input.io.KeyMap(imgui.KeyUpArrow, int(glfw.KeyUp))
	input.io.KeyMap(imgui.KeyDownArrow, int(glfw.KeyDown))
	input.io.KeyMap(imgui.KeyPageUp, int(glfw.KeyPageUp))
	input.io.KeyMap(imgui.KeyPageDown, int(glfw.KeyPageDown))
	input.io.KeyMap(imgui.KeyHome, int(glfw.KeyHome))
	input.io.KeyMap(imgui.KeyEnd, int(glfw.KeyEnd))
	input.io.KeyMap(imgui.KeyInsert, int(glfw.KeyInsert))
	input.io.KeyMap(imgui.KeyDelete, int(glfw.KeyDelete))
	input.io.KeyMap(imgui.KeyBackspace, int(glfw.KeyBackspace))
	input.io.KeyMap(imgui.KeySpace, int(glfw.KeySpace))
	input.io.KeyMap(imgui.KeyEnter, int(glfw.KeyEnter))
	input.io.KeyMap(imgui.KeyEscape, int(glfw.KeyEscape))
	input.io.KeyMap(imgui.KeyA, int(glfw.KeyA))
	input.io.KeyMap(imgui.KeyC, int(glfw.KeyC))
	input.io.KeyMap(imgui.KeyV, int(glfw.KeyV))
	input.io.KeyMap(imgui.KeyX, int(glfw.KeyX))
	input.io.KeyMap(imgui.KeyY, int(glfw.KeyY))
	input.io.KeyMap(imgui.KeyZ, int(glfw.KeyZ))

}
func (input *ImguiWrapping) keyChange(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		input.io.KeyPress(int(key))
	}
	if action == glfw.Release {
		input.io.KeyRelease(int(key))
	}

	// Modifiers are not reliable across systems
	input.io.KeyCtrl(int(glfw.KeyLeftControl), int(glfw.KeyRightControl))
	input.io.KeyShift(int(glfw.KeyLeftShift), int(glfw.KeyRightShift))
	input.io.KeyAlt(int(glfw.KeyLeftAlt), int(glfw.KeyRightAlt))
	input.io.KeySuper(int(glfw.KeyLeftSuper), int(glfw.KeyRightSuper))
}
func (input *ImguiWrapping) charChange(window *glfw.Window, char rune) {
	input.io.AddInputCharacters(string(char))
}
func (input *ImguiWrapping) PID(ticker *time.Ticker) {
	controlTicker := time.NewTicker(time.Second / 2)
	pidtime := int64(10000)
	kp := float32(1000000000)
	ki := float32(900000000)
	kd := float32(500000000)
	I := float32(0)
	prErr := float32(0.0)
	for true {
		frameTime := input.io.Framerate()
		if frameTime < 0.00001 {
			<-controlTicker.C
			continue
		}
		Time := 1.0 / (input.io.Framerate())
		RequiredTime := 1.0 / (float32(env.FPS) - 0.001)
		errorVal := RequiredTime - Time
		P := errorVal
		I += errorVal
		D := errorVal - prErr
		pidtime = int64(P*kp + I*ki + D*kd)
		println("PID:", pidtime, "Time:", Time, "RequiredTime:", RequiredTime)
		if pidtime < 0 {
			pidtime = 0
		} else if time.Duration(pidtime) > time.Second/2 {
			pidtime = int64(time.Second / 2)
		}
		ticker.Reset(time.Nanosecond*time.Duration(pidtime) + 2)
		prErr = P
		<-controlTicker.C
	}
}
func (input *ImguiWrapping) Run(runner func(), ticker *time.Ticker) {
	input.runner = runner
	//FrameTime Controller
	//go input.PID(ticker)

	mainthread.Run(func() {
		for !input.platform.ShouldStop() {
			mainthread.Call(input.Render)
			<-ticker.C
		}
	})

}
func (input *ImguiWrapping) Render() {
	input.platform.ProcessEvents()
	input.NewFrame()
	input.runner()
	p := input.platform
	r := input.renderer
	imgui.Render()
	input.Clear()
	r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
	p.PostRender()
}
func (input *ImguiWrapping) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(
		0,
		0,
		0,
		255)
}
func (input *ImguiWrapping) sizeChange(w *glfw.Window, width int, height int) {
	input.Render()
}
func (input *ImguiWrapping) mouseButtonChange(window *glfw.Window, rawButton glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	buttonIndex, known := glfwButtonIndexByID[rawButton]

	if known && (action == glfw.Press) {
		input.mouseJustPressed[buttonIndex] = true
	}
}
func (input *ImguiWrapping) mouseScrollChange(window *glfw.Window, x, y float64) {
	input.io.AddMouseWheelDelta(float32(x), float32(y))
}
