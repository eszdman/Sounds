package renderer

import "time"

var FPS = int32(180)
var Ticker = time.NewTicker(time.Second / time.Duration(FPS))
