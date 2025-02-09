package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
)

func main() {
	// 连接到 localhost 的 8080 端口
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 提示用户输入文件路径
	fmt.Print("Enter the path to the .txt file: ")
	reader := bufio.NewReader(os.Stdin)
	filePath, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading file path:", err)
		return
	}
	filePath = strings.TrimSpace(filePath)

	// 根据操作系统自动转换路径
	if runtime.GOOS == "windows" {
		// 将 Windows 路径转换为 Linux 路径
		filePath = strings.Replace(filePath, `\`, `/`, -1) // 替换反斜杠为正斜杠
		filePath = "/mnt" + filePath[1:]                   // 去掉 "C:" 前缀并加上 "/mnt"
	}

	// 发送文件路径到服务端
	_, err = conn.Write([]byte(filePath))
	if err != nil {
		fmt.Println("Error sending file path to server:", err)
		return
	}

	// 接收服务端的确认消息
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}

	// 打印服务端返回的结果
	fmt.Printf("Server response: %s\n", string(buffer[:n]))
}
