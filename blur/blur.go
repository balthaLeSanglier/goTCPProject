package blur

import (
	"image"
	"image/color"
	"math"
)

// Générer dynamiquement un noyau gaussien en fonction du radius
func generateGaussianKernel(radius int, dim int) [][]float64 {
	kernel := make([][]float64, dim)
	sigma := float64(radius) / 2.0 // Ajuste l'écart-type pour le flou
	twoSigmaSq := 2 * sigma * sigma
	center := radius

	for i := 0; i < dim; i++ {
		kernel[i] = make([]float64, dim)
		for j := 0; j < dim; j++ {
			x := float64(i - center)
			y := float64(j - center)
			kernel[i][j] = math.Exp(-(x*x+y*y)/twoSigmaSq) / (math.Pi * twoSigmaSq)
		}
	}

	return kernel
}

// Calculer la somme des poids dans le noyau
func sumKernel(kernel [][]float64) float64 {
	sum := 0.0
	for i := 0; i < len(kernel); i++ {
		for j := 0; j < len(kernel[i]); j++ {
			sum += kernel[i][j]
		}
	}
	return sum
}

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

	// Dimension du noyau
	dim := 2*params.radius + 1
	kernel := generateGaussianKernel(params.radius, dim)
	kernelSum := sumKernel(kernel)

	for x := params.bounds.Min.X; x < params.bounds.Max.X; x++ {
		for y := params.yStart; y < params.yEnd; y++ {
			var r, g, b float64

			// Appliquer le flou sur chaque pixel
			for i := -params.radius; i <= params.radius; i++ {
				for j := -params.radius; j <= params.radius; j++ {
					px := clamp(x+i, params.bounds.Min.X, params.bounds.Max.X-1)
					py := clamp(y+j, params.bounds.Min.Y, params.bounds.Max.Y-1)
					srcColor := params.img.At(px, py)
					rSrc, gSrc, bSrc, _ := srcColor.RGBA()
					weight := kernel[i+params.radius][j+params.radius]
					r += float64(rSrc) * weight
					g += float64(gSrc) * weight
					b += float64(bSrc) * weight
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
