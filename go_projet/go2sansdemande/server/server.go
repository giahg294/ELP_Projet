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

// SessionPool gérer TCP pool de connexion
type SessionPool struct {
	mu             sync.Mutex
	pool           map[net.Conn]time.Time // 存储连接及其最后使用时间
	maxConnections int                    // 最大连接数
	idleTimeout    time.Duration          // 空闲超时时间
}

// NewSessionPool créer un nouveau SessionPool
func NewSessionPool(maxConnections int, idleTimeout time.Duration) *SessionPool {
	return &SessionPool{
		pool:           make(map[net.Conn]time.Time),
		maxConnections: maxConnections,
		idleTimeout:    idleTimeout,
	}
}

// AddSession Ajouter une nouvelle connexion TCP au pool
func (sp *SessionPool) AddSession(conn net.Conn) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	// Si le nombre maximum de connexions est atteint, fermez la connexion la plus ancienne
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

// RemoveSession Supprimer une connexion du pool
func (sp *SessionPool) RemoveSession(conn net.Conn) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	delete(sp.pool, conn)
	fmt.Printf("Connection removed from pool: %v\n", conn.RemoteAddr())
}

// getOldestSession Obtenez la première connexion
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

// handleConnection Gérer les connexions clients
func handleConnection(sp *SessionPool, conn net.Conn) {

	defer func() {
		sp.RemoveSession(conn) // Supprimer la connexion lorsque la connexion est fermée
		conn.Close()
	}()
	// sp.AddSession(conn)

	scanner := bufio.NewScanner(conn)

	// Lire le numéro de fichier
	if !scanner.Scan() {
		fmt.Println("Erreur de lecture du numéro de fichier")
		return
	}
	fileNumber, err := strconv.Atoi(scanner.Text())
	if err != nil || fileNumber < 1 || fileNumber > 3 {
		fmt.Println("Numéro de fichier invalide")
		return
	}

	// Générer le nom du fichier
	fileName := fmt.Sprintf("data/minigraph%d.txt", fileNumber)

	// Lire le fichier correspondant
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

	// Calculer la simultanéité de Louvain (sans bloquer le thread principal)
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

	// Envoyer les résultats du calcul : appelez la fonction DisplayCommunities et envoyez la chaîne renvoyée au client
	communityOutput := graph.DisplayCommunities()
	fmt.Fprintf(conn, "%s", communityOutput)

	// // Mettre à jour la dernière fois que la connexion a été utilisée
	sp.AddSession(conn)
}

func main() {
	// Créez un pool de connexions avec un nombre maximum de connexions de 2 et un délai d'inactivité de 1 minute.
	sp := NewSessionPool(2, 20*time.Second)

	// Créer un écouteur TCP
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
		go handleConnection(sp, conn) // Gérer plusieurs clients simultanément
	}
}
