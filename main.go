package main

import (
	"goTcpProject/blur"
	"image"
	"image/jpeg"
	"log"
	"os"
)

func main() {

	file, err := os.Open("input.jpg") // Remplacez par le chemin de votre image
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer file.Close()

	// Décoder l'image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	// Appliquer le flou simple
	blurredImg := blur.GaussianBlur(img)
	// Créer un fichier de sortie pour l'image floutée
	outFile, err := os.Create("output.jpg") // Chemin pour l'image floutée
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Encoder et sauvegarder l'image floutée
	err = jpeg.Encode(outFile, blurredImg, nil)
	if err != nil {
		log.Fatalf("Failed to encode output image: %v", err)
	}

	log.Println("Image floutée sauvegardée sous 'output.jpg'")

}
