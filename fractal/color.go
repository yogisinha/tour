// The fractal package provides a simple fractal browsing web page.
package fractal

import (
	"image/color"
)

var black = color.Black

// Ramp returns the color to use if the Mandelbrot calculation
// stopped after iter iterations out of a maximum possible maxIter.
// The color is black for points still in the set (iter == maxIter) and
// otherwise a blue whose brightness corresponds to the number
// of iterations.
func Ramp(iter, maxIter int) color.Color {
	i := iter
	n := maxIter
	if i >= n {
		return black
	}
	if i >= n/2 {
		v := byte((i - n/2) * 255 / (n - n/2))
		//return image.RGBAColor{v, v, 255, 255}
		return MyColor{uint32(v), uint32(v), 255, 255}
	}
	v := byte(i * 255 / (n / 2))
	//return image.RGBAColor{0, 0, v, 255}
	return MyColor{0, 0, uint32(v), 255}
}

type MyColor struct {
	r, g, b, a uint32
}

func (c MyColor) RGBA() (r, g, b, a uint32) {
	return c.r, c.g, c.b, c.a
}

// Cycle returns the color to use if the Mandelbrot calculation
// stopped after iter iterations out of a maximum possible maxIter.
// The color is black for points still in the set (iter == maxIter) and
// otherwise cycles between red and blue depending on the number
// of iterations.
// func Cycle(iter, maxIter int) image.RGBAColor {
// 	i := iter
// 	n := maxIter
// 	if i >= n {
// 		return black
// 	}
// 	if i%2 == 1 {
// 		return image.RGBAColor{255, 0, 0, 255}
// 	}
// 	return image.RGBAColor{0, 0, 255, 255}
// }
