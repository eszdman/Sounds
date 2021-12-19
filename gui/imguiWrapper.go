package gui

import (
	"audioGeneration/platform"
	"github.com/inkyblackness/imgui-go"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v3.2-core/gl"
	"math"
)

// ImguiWrapping is the state holder for the imgui framework
type ImguiWrapping struct {
	io               *imgui.IO
	time             float64
	window *glfw.Window
	platform *platform.Platform
	context *imgui.Context
	renderer *OpenGL3
	runner func()
	mouseJustPressed [3]bool
}

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
func NewImgui(platform *platform.Platform,renderer *OpenGL3,context *imgui.Context) (*ImguiWrapping) {
	imguiIO = imgui.CurrentIO()
	inputState = ImguiWrapping{
		io: &imguiIO,
		window: platform.Window,
		platform: platform,
		renderer: renderer,
		context: context,
		time: 0}
	inputState.setKeyMapping()
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
func (input *ImguiWrapping) NewFrame() {
	cursorX, cursorY := input.platform.GetCursorPos()
	mouseState := ImguiMouseState{
		MousePosX:  float32(cursorX),
		MousePosY:  float32(cursorY),
		MousePress: input.platform.GetMousePresses123(),
	}
	// Setup display size (every frame to accommodate for window resizing)
	sizes := input.platform.DisplaySize()
	for i:=0; i<len(sizes);i++ {
		if(sizes[i] <= 500){
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
		down := input.mouseJustPressed[i] || mouseState.MousePress[0] == true
		input.io.SetMouseButtonDown(i, down)
		input.mouseJustPressed[i] = false
	}

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
func (input *ImguiWrapping) Run(runner func()){
	input.runner = runner
	for !input.platform.ShouldStop() {
		input.Render()
	}
}
func (input *ImguiWrapping) Render(){
	input.runner()
	p := input.platform
	r := input.renderer
	imgui.Render()
	input.Clear()
	r.Render(p.DisplaySize(),p.FramebufferSize(),imgui.RenderedDrawData())
	p.PostRender()
}
func (input *ImguiWrapping) Clear(){
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
