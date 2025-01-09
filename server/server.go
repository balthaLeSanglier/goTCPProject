package server

import (
	"bytes"
	"fmt" // Importer le package pour le flou
	"goTcpProjectServer/blur"
	"image"
	"image/jpeg"
	"log"
	"net"
	"os"
)

// Initialise le serveur sur le port 8080. Accepte toutes les connexions entrantes et
// lance le traitement de l'image pour chaque connexion reçu.
func StartServer() {
	// Créer un écouteur sur le port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erreur lors de l'écoute:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Serveur TCP en écoute sur le port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur lors de l'acceptation de la connexion:", err)
			continue
		}
		fmt.Println("Nouvelle connexion reçue !")
		go handleConnection(conn) //Les connexions sont traitées dans des go routine pour ne pas bloqué l'éxecution
	}
}

// Traitement de la connexion :
//   - reception de l'image
//   - application du flou
//   - envoi de l'image flouté
//
// Paramètres :
//
//	conn (net.Conn) : La connexion à traiter
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connexion reçue, en attente de l'image...")

	img, err := receiveImage(conn)
	if err != nil {
		log.Println("Erreur lors de la réception de l'image :", err)
		return
	}

	fmt.Println("Image reçue avec succès.")
	blurred := blur.StartGaussianBlur(img, 4)

	err = sendImage(conn, blurred)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de l'image:", err)
	}
}

// sendImage encode une image au format JPEG et l'envoie via la connexion TCP.
//
// Paramètres :
//
//	conn (net.Conn) : La connexion à traiter
//	img *image.RGBA : l'image à envoyer
func sendImage(conn net.Conn, img *image.RGBA) error {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		return err
	}
	_, err = conn.Write(buf.Bytes())
	fmt.Println("image envoyé")
	fmt.Printf("Taille de l'image envoyée : %d octets\n", buf.Len())
	return err
}

// Décode les bytes envoyés par le client afin
//
// conn : La connexion TCP à partir de laquelle l'image est reçue (conn contient les octets envoyé par le client)
//
// Retourne :
//   - Une image décodée de type image.Image si le décodage est réussi.
//   - Une erreur détaillant le problème si le décodage échoue.
func receiveImage(conn net.Conn) (image.Image, error) {
	img, _, err := image.Decode(conn)
	if err != nil {
		return nil, fmt.Errorf("échec du décodage de l'image : %v", err)
	}
	return img, nil
}
