package main

import "github.com/fogleman/gg"

func main() {
	const width = 300
	const height = 300
	dc := gg.NewContext(width, height)
	dc.SetRGBA(100, 255, 100, 0.2)
	for i := 0; i < 360; i += 30 {
		dc.Push()
		dc.RotateAbout(gg.Radians(float64(i)), width/2, height/2)
		dc.DrawEllipse(width/2, height/2, height*7/16, height/8)
		dc.Fill()
		dc.Pop()
	}
	dc.SetRGBA(144, 25, 96, 0.2)
	for i := 15; i < 360; i += 30 {
		dc.Push()
		dc.RotateAbout(gg.Radians(float64(i)), width/2, height/2)
		dc.DrawEllipse(width/2, height/2, height*7/16, height/8)
		dc.Fill()
		dc.Pop()
	}
	// Настраиваем цвет фона для текста
	dc.SetRGB(0, 0, 0)

	// Устанавливаем шрифт и размер
	dc.LoadFontFace("sans-serif", 12)

	dc.DrawStringWrapped("Samson Lucius", width/2, height/2, 0.5, 0.5, 5, 1, 1)

	dc.SavePNG("amazing_logo.png")
}
