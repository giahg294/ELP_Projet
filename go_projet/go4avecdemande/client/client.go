package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

func main() {
	// 让用户选择文件编号
	var fileNumber int
	fmt.Print("Veuillez choisir un fichier (1-3) : ")
	_, err := fmt.Scanln(&fileNumber)
	if err != nil || fileNumber < 1 || fileNumber > 3 {
		fmt.Println("Numéro invalide ! Veuillez entrer 1, 2 ou 3.")
		return
	}

	// 让用户输入 Goroutines 的数量
	var numGoroutines int
	fmt.Print("Veuillez entrer le nombre de Goroutines : ")
	_, err = fmt.Scanln(&numGoroutines)
	if err != nil || numGoroutines <= 0 {
		fmt.Println("Nombre de Goroutines invalide !")
		return
	}

	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8882")
	if err != nil {
		fmt.Println("Erreur de connexion au serveur :", err)
		return
	}
	defer conn.Close()

	// 发送 Goroutines 的数量到服务器
	N_routine := int64(numGoroutines)
	err = binary.Write(conn, binary.BigEndian, N_routine)
	if err != nil {
		fmt.Println("Erreur d'envoi du nombre de go routines au serveur", err)
		return
	}

	// 发送文件编号给服务器
	writer := bufio.NewWriter(conn)
	fmt.Fprintln(writer, strconv.Itoa(fileNumber))
	writer.Flush()

	// 读取服务器返回的结果
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
