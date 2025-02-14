package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

func main() {
	// Laissez l'utilisateur sélectionner un numéro de fichier
	var fileNumber int
	fmt.Print("Veuillez choisir un fichier (1-4) : ")
	_, err := fmt.Scanln(&fileNumber)
	if err != nil || fileNumber < 1 || fileNumber > 4 {
		fmt.Println("Numéro invalide ! Veuillez entrer 1, 2 ou 3.")
		return
	}

	// Laissez l'utilisateur saisir le nombre de Goroutines
	var numGoroutines int
	fmt.Print("Veuillez entrer le nombre de Goroutines : ")
	_, err = fmt.Scanln(&numGoroutines)
	if err != nil || numGoroutines <= 0 {
		fmt.Println("Nombre de Goroutines invalide !")
		return
	}

	// Se connecter au serveur
	conn, err := net.Dial("tcp", "localhost:8882")
	if err != nil {
		fmt.Println("Erreur de connexion au serveur :", err)
		return
	}
	defer conn.Close()

	// Envoyer le nombre de Goroutines au serveur
	N_routine := int64(numGoroutines)
	err = binary.Write(conn, binary.BigEndian, N_routine)
	if err != nil {
		fmt.Println("Erreur d'envoi du nombre de go routines au serveur", err)
		return
	}

	// Envoyer le numéro de fichier au serveur
	writer := bufio.NewWriter(conn)
	fmt.Fprintln(writer, strconv.Itoa(fileNumber))
	writer.Flush()

	// Lire les résultats renvoyés par le serveur
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
