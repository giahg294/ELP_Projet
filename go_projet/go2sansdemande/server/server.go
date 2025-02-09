package main

import (
	"bufio"
	"fmt"
	"go2sansdemande/traitement"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SessionPool 管理 TCP 连接池
type SessionPool struct {
	mu             sync.Mutex
	pool           map[net.Conn]time.Time // 存储连接及其最后使用时间
	maxConnections int                    // 最大连接数
	idleTimeout    time.Duration          // 空闲超时时间
}

// NewSessionPool 创建一个新的 SessionPool
func NewSessionPool(maxConnections int, idleTimeout time.Duration) *SessionPool {
	return &SessionPool{
		pool:           make(map[net.Conn]time.Time),
		maxConnections: maxConnections,
		idleTimeout:    idleTimeout,
	}
}

// AddSession 添加一个新的 TCP 连接到池中
func (sp *SessionPool) AddSession(conn net.Conn) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	// 如果达到最大连接数，关闭最早的连接
	if len(sp.pool) >= sp.maxConnections {
		oldestConn := sp.getOldestSession()
		if oldestConn != nil {
			fmt.Printf("Connection pool is full. Closing oldest connection: %v\n", oldestConn.RemoteAddr())
			oldestConn.Close()
			delete(sp.pool, oldestConn)
		}
	}

	sp.pool[conn] = time.Now()
	fmt.Printf("New connection added to pool: %v\n", conn.RemoteAddr())
}

// RemoveSession 从池中移除一个连接
func (sp *SessionPool) RemoveSession(conn net.Conn) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	delete(sp.pool, conn)
	fmt.Printf("Connection removed from pool: %v\n", conn.RemoteAddr())
}

// getOldestSession 获取最早的连接
func (sp *SessionPool) getOldestSession() net.Conn {
	var oldestConn net.Conn
	var oldestTime time.Time
	for conn, t := range sp.pool {
		if oldestConn == nil || t.Before(oldestTime) {
			oldestConn = conn
			oldestTime = t
		}
	}
	return oldestConn
}

// handleConnection 处理客户端连接
func handleConnection(sp *SessionPool, conn net.Conn) {

	defer func() {
		sp.RemoveSession(conn) // 在连接关闭时移除连接
		conn.Close()
	}()
	// sp.AddSession(conn)

	scanner := bufio.NewScanner(conn)

	// 读取文件编号
	if !scanner.Scan() {
		fmt.Println("Erreur de lecture du numéro de fichier")
		return
	}
	fileNumber, err := strconv.Atoi(scanner.Text())
	if err != nil || fileNumber < 1 || fileNumber > 3 {
		fmt.Println("Numéro de fichier invalide")
		return
	}

	// 生成文件名
	fileName := fmt.Sprintf("minigraph%d.txt", fileNumber)

	// 读取对应的文件
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(conn, "Erreur: impossible d'ouvrir %s\n", fileName)
		return
	}
	defer file.Close()

	// Créer un nouveau graphe
	graph := traitement.NewGraph()

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		// Ignorer les lignes vides ou les commentaires (lignes qui commencent par "#")
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		nodes := strings.Fields(line) // Découper la ligne par espace
		if len(nodes) == 2 {
			// Convertir les chaînes en entiers
			u, err := strconv.Atoi(nodes[0])
			if err != nil {
				fmt.Println("Erreur lors de la conversion de", nodes[0], "en entier:", err)
				continue
			}
			v, err := strconv.Atoi(nodes[1])
			if err != nil {
				fmt.Println("Erreur lors de la conversion de", nodes[1], "en entier:", err)
				continue
			}

			// Ajouter l'arête au graphe
			graph.AddEdge(u, v)
		}
	}

	if err := fileScanner.Err(); err != nil {
		fmt.Fprintf(conn, "Erreur lors de la lecture du fichier : %s\n", err)
		return
	}

	// 计算 Louvain 并发（不阻塞主线程）
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		graph.Louvain(10)
	}()
	wg.Wait()

	elapsed := time.Since(start)
	if _, ok := sp.pool[conn]; ok {
		fmt.Printf("Louvain algorithm %v took %v to complete\n", conn.RemoteAddr(), elapsed)
	}

	// 发送计算结果：调用 DisplayCommunities 函数，并将返回的字符串发送到客户端
	communityOutput := graph.DisplayCommunities()
	fmt.Fprintf(conn, "%s", communityOutput)

	// // 更新连接的最后使用时间
	sp.AddSession(conn)
}

func main() {
	// 创建一个最大连接数为 2 的连接池，空闲超时时间为 1 分钟
	sp := NewSessionPool(2, 20*time.Second)

	// 创建 TCP 监听
	listener, err := net.Listen("tcp", ":8882")
	if err != nil {
		fmt.Println("Erreur de démarrage du serveur :", err)
		return
	}
	defer listener.Close()

	fmt.Println("Serveur en attente de connexions sur le port 8882...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur d'acceptation :", err)
			continue
		}
		go handleConnection(sp, conn) // 并发处理多个客户端
	}
}
