package gologo

import "github.com/go-gl/glfw/v3.2/glfw"

var TimeState = TickState{}

/////////////////////////////////////////////////////////////
// Tick
//

type TickState struct {
	Zero     float64
	Start    float64
	End      float64
	Interval float64
}

func InitTick() error {
	TimeState.Zero = glfw.GetTime()
	TimeState.End = TimeState.Zero

	return nil
}

func GetTime() int {
	return int(1000 * (glfw.GetTime() - TimeState.Zero))
}

func GetTickTime() int {
	return int(1000 * (TimeState.End - TimeState.Zero))
}

func Tick() {
	time := glfw.GetTime()
	TimeState.Interval = time - TimeState.End
	TimeState.End = time
}
