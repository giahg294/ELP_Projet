package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
)

// Graph 定义一个图结构
type Graph struct {
	AdjList map[string][]string
	mu      sync.Mutex // 用于并发保护
}

// AddEdge 向图中添加一条边
func (g *Graph) AddEdge(node1, node2 string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.AdjList[node1] = append(g.AdjList[node1], node2)
	g.AdjList[node2] = append(g.AdjList[node2], node1) // 无向图
}

// Print 打印图的邻接表
func (g *Graph) Print() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for node, neighbors := range g.AdjList {
		fmt.Printf("%s: %v\n", node, neighbors)
	}
}

// LoadGraphFromFile 从指定路径加载图
func LoadGraphFromFile(filePath string) (*Graph, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	graph := &Graph{AdjList: make(map[string][]string)}
	reader := bufio.NewReader(file)

	// 使用 bufio.Reader 逐行读取文件，处理长行
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				return nil, err
			}
			break
		}
		line = strings.TrimSpace(line)
		if line != "" {
			// 修改这里：使用制表符 `\t` 分割行
			nodes := strings.Split(line, "\t")
			if len(nodes) == 2 {
				graph.AddEdge(nodes[0], nodes[1]) // 添加边
			}
		}
	}

	return graph, nil
}

// handleConnection 处理每个客户端连接
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected.")

	// 接收客户端传递的文件路径
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return
	}

	// 解析客户端发送的文件路径
	filePath := strings.TrimSpace(string(buffer[:n]))
	fmt.Printf("Received file path from client: %s\n", filePath)

	// 如果客户端传递的是 Windows 风格路径，需要转换为 Linux 路径
	if runtime.GOOS == "linux" && strings.HasPrefix(filePath, "C:") {
		// 转换 Windows 路径为 Linux 路径，确保不多出冒号
		filePath = "/mnt/c" + filePath[2:]
		// 替换反斜杠 `\` 为正斜杠 `/`
		filePath = strings.ReplaceAll(filePath, "\\", "/")
	}

	// 加载图结构
	graph, err := LoadGraphFromFile(filePath)
	if err != nil {
		fmt.Fprintf(conn, "Error loading graph from file: %v\n", err)
		fmt.Println("Error loading graph:", err)
		return
	}

	// 打印图结构到服务器控制台
	fmt.Println("Graph loaded successfully. Adjacency list:")
	graph.Print()

	// 向客户端确认加载成功
	_, err = conn.Write([]byte("Graph loaded successfully. Check server logs for adjacency list.\n"))
	if err != nil {
		fmt.Println("Error sending confirmation to client:", err)
		return
	}

	fmt.Println("Graph loading confirmation sent to client.")
}

// main 启动 TCP 服务器
func main() {
	// 监听本地的 8080 端口
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080")

	// 等待客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// 启动一个 Goroutine 处理每个客户端连接
		go handleConnection(conn)
	}
}
