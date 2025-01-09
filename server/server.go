package server

import (
	"bufio"
	"fmt" // Importer le package pour le flou
	"net"
	"os"
)

func StartServer() {
	// Créer un écouteur sur le port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erreur lors de l'écoute:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Serveur TCP en écoute sur le port 8080...")

	// Boucle infinie pour accepter les connexions entrantes
	for {
		// Accepter une nouvelle connexion
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur lors de l'acceptation de la connexion:", err)
			continue
		}

		// Afficher un message lorsqu'une nouvelle connexion est reçue
		fmt.Println("Nouvelle connexion reçue !")

		// Gérer la connexion dans une goroutine pour ne pas bloquer l'acceptation de nouvelles connexions
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Créer un scanner pour lire les messages envoyés par le client
	scanner := bufio.NewScanner(conn)

	// Lire les messages envoyés par le client
	for scanner.Scan() {
		// Récupérer le message envoyé par le client
		message := scanner.Text()

		// Afficher le message reçu dans la console du serveur
		fmt.Println("Message reçu du client :", message)

		// Optionnel : si tu veux renvoyer un message au client
		_, err := conn.Write([]byte("Message reçu : " + message + "\n"))
		if err != nil {
			fmt.Println("Erreur lors de l'envoi au client :", err)
		}
	}

	// Vérifier s'il y a eu une erreur pendant la lecture du scanner
	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lors de la lecture de la connexion :", err)
	}
}
