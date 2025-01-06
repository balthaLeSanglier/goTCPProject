package main

import (
	"fmt"
	"os"
)

func main() {
	// Lire tout le contenu du fichier
	content, err := os.ReadFile("example.txt")
	if err != nil {
		fmt.Println("Erreur en lisant le fichier:", err)
		return
	}

	// Afficher le contenu
	fmt.Println(string(content))
}
