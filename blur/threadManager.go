// TODO
package blur

import (
	"image"
	"sync"
)

type BlurParams struct {
	img       image.Image
	bounds    image.Rectangle
	blurred   *image.RGBA
	yStart    int
	yEnd      int
	waitGroup *sync.WaitGroup
}

// Lance le traitement du flou gaussien sur toute l'image:
//   - Calcul la longueur de bandes horizontal
//   - Lance x thread pour calculer les flous sur les bandes
func StartGaussianBlur(img image.Image, threadAmount int) *image.RGBA {
	bounds := img.Bounds()
	blurred := image.NewRGBA(bounds)
	bandHeight := bounds.Max.Y / threadAmount

	var WaitGroup sync.WaitGroup

	for i := 0; i < threadAmount; i++ {
		yStart := bounds.Min.Y + i*bandHeight
		yEnd := yStart + bandHeight
		if i == threadAmount-1 {
			yEnd = bounds.Max.Y
		}
		WaitGroup.Add(1)
		params := BlurParams{
			bounds:    bounds,
			img:       img,
			blurred:   blurred,
			yStart:    yStart,
			yEnd:      yEnd,
			waitGroup: &WaitGroup,
		}
		go applyBlurSection(params)
	}
	WaitGroup.Wait()
	return blurred
}
