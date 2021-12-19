package ui

import (
	"fmt"
	"github.com/eszdman/Sounds/env"
	"github.com/eszdman/Sounds/ui/forms"
	"github.com/eszdman/Sounds/ui/gui"
	"github.com/eszdman/Sounds/ui/platform"
	"github.com/inkyblackness/imgui-go/v4"
	"os"
)

const (
	ScrollBarSize = 10
	PianoOctaves  = 7
	PianoCount    = PianoOctaves * 12
)
const windowWidth = 1600
const windowHeight = 900

func RunUI() {
	var err error
	env.NewPlatform, err = platform.NewPlatform(windowWidth, windowHeight, "SoundEngine")
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
	defer env.NewPlatform.Dispose()
	defer context.Destroy()
	defer imguiRenderer.Dispose()
	imguiWrapping := gui.NewImgui(env.NewPlatform, imguiRenderer, context)
	imgui.StyleColorsClassic()
	forms.UseDefaultPianoRoll()
	mainRunner := func() {
		env.NewPlatform.ProcessEvents()
		imguiWrapping.NewFrame()
		{
			forms.MenuBar(&keep, env.NewPlatform.DisplaySize())
			if keep {
				imgui.ShowDemoWindow(&keep)
			}
			forms.PianoRoll()
			if env.PianoSettings {
				forms.PianoRollSettings(&env.PianoSettings)
			}
		}
	}
	imguiWrapping.Run(mainRunner)
}
