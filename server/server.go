package server

import (
	"bytes"
	"encoding/binary"
	"fmt" // Importer le package pour le flou
	"goTcpProjectServer/blur"
	"image"
	"image/png"
	"log"
	"net"
	"os"
	"time"
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

	fmt.Println("Connexion reçue, en attente du nombre de goRoutine...")

	var goRoutineNumber int32
	goRoutineNumber, err := receiveGoRoutineNumber(conn)
	if err != nil {
		fmt.Println("Erreur lors de la réception du nombre de GoRoutine :", err)
		return
	}

	var Radius int32
	Radius, err_radius := receiveRadius(conn)
	if err_radius != nil {
		fmt.Println("Erreur lors de la réception du rayon :", err)
		return
	}

	fmt.Printf("Nombre de GoRoutine reçu : %d\n", goRoutineNumber)

	fmt.Println("En attente de l'image...")

	img, err := receiveImage(conn)
	if err != nil {
		log.Println("Erreur lors de la réception de l'image :", err)
		return
	}

	fmt.Println("Image reçue avec succès.")
	start := time.Now()
	blurred := blur.StartGaussianBlur(img, int(goRoutineNumber), int(Radius))
	elapsed := time.Since(start)
	fmt.Println("Temps écoulé pendant le flou : ", elapsed)

	err = sendImage(conn, blurred)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de l'image:", err)
	}
}

// Décode les bytes envoyés par le client afin de retrouvé le nombre de goRoutine demandé par le client
//
// conn : La connexion TCP à partir de laquelle l'image est reçue (conn contient les octets envoyé par le client)
//
// Retourne :
//   - le nombre de goRoutine.
//   - Une erreur détaillant le problème si le décodage échoue.
func receiveGoRoutineNumber(conn net.Conn) (int32, any) {
	var goRoutineNumber int32
	err := binary.Read(conn, binary.BigEndian, &goRoutineNumber)
	return goRoutineNumber, err
}

// Similairement aux nombre de go routine, le rayon est aussi réccupéré

func receiveRadius(conn net.Conn) (int32, any) {
	var Radius int32
	err := binary.Read(conn, binary.BigEndian, &Radius)
	return Radius, err
}

// Décode les bytes envoyés par le client afin de reconstitué l'image
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

// sendImage encode une image au format PNG et l'envoie via la connexion TCP.
//
// Paramètres :
//
//	conn (net.Conn) : La connexion à traiter
//	img *image.RGBA : l'image à envoyer
func sendImage(conn net.Conn, img *image.RGBA) error {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return err
	}
	_, err = conn.Write(buf.Bytes())
	fmt.Println("image envoyé")
	fmt.Printf("Taille de l'image envoyée : %d octets\n", buf.Len())
	return err
}
