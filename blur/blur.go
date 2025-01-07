// TODO
package blur

import (
	"image"
	"image/color"
)

func GaussianBlur(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	blurred := image.NewRGBA(bounds) //nouvelle image (la version flouté de l'input)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			var r, g, b float64

			// Appliquer le flou sur chaque pixel
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					px := clamp(x+i, bounds.Min.X, bounds.Max.X-1)
					py := clamp(y+j, bounds.Min.Y, bounds.Max.Y-1)
					srcColor := img.At(px, py)
					rSrc, gSrc, bSrc, _ := srcColor.RGBA()

					r += float64(rSrc)
					g += float64(gSrc)
					b += float64(bSrc)
				}
			}
			blurred.Set(x, y, color.RGBA{
				R: uint8(r / 9 / 256),
				G: uint8(g / 9 / 256),
				B: uint8(b / 9 / 256),
				A: uint8(255), // Garder l'opacité maximale
			})
		}
	}
	return blurred
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
