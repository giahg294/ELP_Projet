package main

import (
	"bufio"
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

	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8882")
	if err != nil {
		fmt.Println("Erreur de connexion au serveur :", err)
		return
	}
	defer conn.Close()

	// 发送编号给服务器
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
