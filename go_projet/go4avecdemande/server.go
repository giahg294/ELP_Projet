package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"go4avecdemande/traitement"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SessionPool gère le pool de connexions TCP
type SessionPool struct {
	mu             sync.Mutex
	pool           map[net.Conn]time.Time // 存储连接及其最后使用时间
	maxConnections int                    // 最大连接数
	idleTimeout    time.Duration          // 空闲超时时间
}

// NewSessionPool crée un nouveau SessionPool
func NewSessionPool(maxConnections int, idleTimeout time.Duration) *SessionPool {
	return &SessionPool{
		pool:           make(map[net.Conn]time.Time),
		maxConnections: maxConnections,
		idleTimeout:    idleTimeout,
	}
}

// AddSession ajoute une nouvelle connexion TCP au pool
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
		sp.RemoveSession(conn) // Supprimer la connexion lorsqu'elle est fermée
		conn.Close()
	}()
	sp.AddSession(conn)

	// Lire le nombre de Goroutines
	var N_routine int64
	err := binary.Read(conn, binary.BigEndian, &N_routine)
	if err != nil {
		fmt.Println("Erreur de lecture du nombre de go routines du client", err)
		return
	}
	fmt.Printf("Received number of go routines from client: %d\n", N_routine)

	// Lire le numéro de fichier
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		fmt.Println("Erreur de lecture du numéro de fichier")
		return
	}
	fileNumber, err := strconv.Atoi(scanner.Text())
	if err != nil || fileNumber < 1 || fileNumber > 4 {
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
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		nodes := strings.Fields(line)
		if len(nodes) == 2 {
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
		graph.Louvain(10, int(N_routine)) // Utiliser le nombre de Goroutines transmis par le client
	}()
	wg.Wait()

	elapsed := time.Since(start)
	if _, ok := sp.pool[conn]; ok {
		fmt.Printf("Louvain algorithm %v took %v to complete\n", conn.RemoteAddr(), elapsed)
	}

	// Envoyer les résultats du calcul : appelez la fonction DisplayCommunities et envoyez la chaîne renvoyée au client
	communityOutput := graph.DisplayCommunities()
	fmt.Fprintf(conn, "%s", communityOutput)

	// Envoyer le drapeau de fin
	fmt.Fprintln(conn, "FIN")
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
