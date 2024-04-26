package utils

import (
	"image/color"
)

func GetRedColor() color.Color {
	return color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}
}

func GetGreenColor() color.Color {
	return color.RGBA {
		R: 70,
		G: 255,
		B: 0,
		A: 255,
	}
}

func GetWhiteColor() color.Color {
	return color.RGBA {
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
}

func GetBlackColor() color.Color {
	return color.RGBA{
		R: 10, 
		G: 30, 
		B: 10,
		A: 255,
	}
}
