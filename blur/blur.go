package blur

import (
	"image"
	"image/color"
)

// Flou gaussien avec pondération appliqué à une section de l'image.
//
// utilise le multithreading, décrémente le WaitGroup de 1 une fois la fonction terminé
//
// Paramètres :
//
//	img *image.Image : l'image à envoyer
//
// Retourne :
//
//	*image.RGBA : image flouté
func applyBlurSection(params BlurParams) {
	defer params.waitGroup.Done()

	// kernel := [3][3]float64{
	// 	{1, 2, 1},
	// 	{2, 4, 2},
	// 	{1, 2, 1},
	// }
	kernelSum := 121.0

	for x := params.bounds.Min.X; x < params.bounds.Max.X; x++ {
		for y := params.yStart; y < params.yEnd; y++ {
			var r, g, b float64

			// Appliquer le flou sur chaque pixel
			for i := -5; i <= 5; i++ {
				for j := -5; j <= 5; j++ {
					px := clamp(x+i, params.bounds.Min.X, params.bounds.Max.X-1)
					py := clamp(y+j, params.bounds.Min.Y, params.bounds.Max.Y-1)
					srcColor := params.img.At(px, py)
					rSrc, gSrc, bSrc, _ := srcColor.RGBA()
					//weight := kernel[i+1][j+1]
					r += float64(rSrc) // * weight
					g += float64(gSrc) // * weight
					b += float64(bSrc) // * weight
				}
			}

			params.blurred.Set(x, y, color.RGBA{
				R: uint8(r / kernelSum / 256),
				G: uint8(g / kernelSum / 256),
				B: uint8(b / kernelSum / 256),
				A: uint8(255), // Garder l'opacité maximale
			})
		}
	}
}

// Flou gaussien de avec pondération
//
// Paramètres :
//
//	img *image.Image : l'image à envoyer
//
// Retourne :
//
//	*image.RGBA : image flouté
func GaussianBlur(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	blurred := image.NewRGBA(bounds) //nouvelle image (la version flouté de l'input)

	kernel := [3][3]float64{
		{1, 2, 1},
		{2, 4, 2},
		{1, 2, 1},
	}
	kernelSum := 16.0

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
					weight := kernel[i+1][j+1]
					r += float64(rSrc) * weight
					g += float64(gSrc) * weight
					b += float64(bSrc) * weight
				}
			}

			blurred.Set(x, y, color.RGBA{
				R: uint8(r / kernelSum / 256),
				G: uint8(g / kernelSum / 256),
				B: uint8(b / kernelSum / 256),
				A: uint8(255), // Garder l'opacité maximale
			})
		}
	}
	return blurred
}

// Flou gaussien de base, sans pondération
//
// Paramètres :
//
//	img *image.Image : l'image à envoyer
//
// Retourne :
//
//	*image.RGBA : image flouté
func GaussianBlurSimple(img image.Image) *image.RGBA {
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
			})
		}
	}
	return blurred
}

// clamp limite une valeur à une plage donnée.
// Si la valeur est inférieure au minimum, elle est fixée au minimum.
// Si la valeur est supérieure au maximum, elle est fixée au maximum.
// Sinon, la valeur reste inchangée.
func clamp(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
