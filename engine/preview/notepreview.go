package preview

import (
	"fmt"
	"github.com/gen2brain/malgo"
	"math"
	"os"
)

type NotePreviewer struct {
	context      *malgo.AllocatedContext
	device       *malgo.Device
	deviceConfig *malgo.DeviceConfig
	sound        []int32
	dragging     bool
	reset        bool
	freq         float64
}

var notePreviewer NotePreviewer

const sampleRate = 44100

func OnMouseReleased() {
	notePreviewer.dragging = false
}
func onChangedPitch() {
	notePreviewer.reset = true
}
func OnMouseDrag(frequency float64) {
	if notePreviewer.freq != frequency {
		onChangedPitch()
	}
	notePreviewer.freq = frequency
	notePreviewer.dragging = true
}

// PutUint32 encodes a uint32 into buf and returns the number of bytes written.
// If the buffer is too small, PutUint32 will panic.
func PutUint32(buf []byte, x uint32) int {
	i := 0
	for x > 0xFF {
		buf[i] = byte(x) | 0xFF
		x >>= 8
		i++
	}
	buf[i] = byte(x)
	return i + 1
}
func Init() {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(-4)
	}
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS32
	deviceConfig.Playback.Channels = 1
	deviceConfig.SampleRate = sampleRate
	deviceConfig.Alsa.NoMMap = 1
	notePreviewer.deviceConfig = &deviceConfig
	notePreviewer.context = ctx
	// This is the function that's used for sending more data to the device for playback.
	onSamples := func(pOutputSample, pInputSamples []byte, frameCount uint32) {
		intArr := notePreviewer.sound[0:frameCount]
		soundGen(intArr)
		for i := 0; i < len(pOutputSample)/4; i++ {
			PutUint32(pOutputSample[i*4:i*4+4], uint32(intArr[i]))
		}
	}
	deviceCallbacks := malgo.DeviceCallbacks{
		Data: onSamples,
	}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		fmt.Println(err)
		os.Exit(-4)
	}
	notePreviewer.sound = make([]int32, deviceConfig.SampleRate/2)
	notePreviewer.device = device
	err = device.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(-4)
	}
}

var cnt = 0
var step = 0
var prevFreq = float64(0)

func soundGen(arr []int32) {
	for i := 0; i < len(arr); i++ {
		if notePreviewer.dragging && !notePreviewer.reset {
			step++
		} else {
			step--
		}
		if step > 128 {
			step = 128
		} else if step < 0 {
			step = 0
			prevFreq = notePreviewer.freq
			notePreviewer.reset = false
		}
		arr[i] = int32(math.Sin(float64(i+cnt)*(prevFreq*10)/float64(sampleRate)) * (math.MaxInt32 / 8) * float64(step) / float64(128))
	}
	cnt += len(arr)
}

func DeInit() {
	notePreviewer.device.Uninit()
	_ = notePreviewer.context.Uninit()
	notePreviewer.context.Free()
}
