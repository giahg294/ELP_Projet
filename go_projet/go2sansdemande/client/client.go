package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

func main() {
	// Laissez l'utilisateur sélectionner un numéro de fichier
	var fileNumber int
	fmt.Print("Veuillez choisir un fichier (1-3) : ")
	_, err := fmt.Scanln(&fileNumber)
	if err != nil || fileNumber < 1 || fileNumber > 3 {
		fmt.Println("Numéro invalide ! Veuillez entrer 1, 2 ou 3.")
		return
	}

	// Se connecter au serveur
	conn, err := net.Dial("tcp", "localhost:8882")
	if err != nil {
		fmt.Println("Erreur de connexion au serveur :", err)
		return
	}
	defer conn.Close()

<<<<<<< HEAD
	// Envoyer le numéro au serveur
=======
	// Envoyer le numero au serveur
>>>>>>> 47ccd1611529c2176587f72461728384faaad031
	writer := bufio.NewWriter(conn)
	fmt.Fprintln(writer, strconv.Itoa(fileNumber))
	writer.Flush()

<<<<<<< HEAD
	// Lire les résultats renvoyés par le serveur
=======
	// Lire le résultat renvoyé par le serveur
>>>>>>> 47ccd1611529c2176587f72461728384faaad031
	serverReader := bufio.NewScanner(conn)
	fmt.Println("Résultat de la détection de communauté :")
	for serverReader.Scan() {
		line := serverReader.Text()
		if line == "FIN" {
			break
		}
		fmt.Println(line)
	}
}
