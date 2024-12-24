package main

import (
	"main/wind"

	"github.com/go-vgo/robotgo"
)

//#include <stdlib.h>
import "C"

func main() {
	wind.InitApplication()

	title := "School 21"

	screenWidth, screenHeight := robotgo.GetScreenSize()

	windPtr := wind.CreateWindow(screenWidth/2, screenHeight/2, 300, 200, title)
	defer C.free(windPtr)

	wind.MakeKeyAndOrderFront(windPtr)
	wind.RunApplication()
}
