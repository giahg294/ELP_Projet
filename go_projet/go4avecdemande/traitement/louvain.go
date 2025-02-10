package traitement

import (
	"sort"
	"sync"
	"strings"
	"fmt"
)

// Graph est une structure représentant un graphe avec des communautés
type Graph struct {
	AdjList     map[int][]int
	Communities map[int]int
	mu          sync.Mutex // Mutex pour protéger l'accès aux communautés
}

// NewGraph crée un nouveau graphe vide
func NewGraph() *Graph {
	return &Graph{
		AdjList:     make(map[int][]int),
		Communities: make(map[int]int),
	}
}

func (g *Graph) DisplayCommunities() string {
	var result strings.Builder
	communityGroups := make(map[int][]int)
	for node, community := range g.Communities {
		communityGroups[community] = append(communityGroups[community], node)
	}

	// Sort communities and nodes
	for community, nodes := range communityGroups {
		sort.Ints(nodes)
		result.WriteString(fmt.Sprintf("Community %d: %v\n", community, nodes))
	}
	return result.String()
}

// Modularity calcule la modularité du graphe
func (g *Graph) Modularity() float64 {
	m := float64(0)
	// Calcul du nombre total d'arêtes
	for _, neighbors := range g.AdjList {
		m += float64(len(neighbors))
	}
	m /= 2

	var Q float64
	// Calcul de la modularité
	for node, neighbors := range g.AdjList {
		community := g.Communities[node]
		ki := float64(len(neighbors))
		for _, neighbor := range neighbors {
			kj := float64(len(g.AdjList[neighbor]))
			if g.Communities[neighbor] == community {
				Q += 1.0 - (ki*kj)/(2.0*m)
			}
		}
	}
	return Q / (2 * m)
}

// MergeCommunities fusionne les communautés en un nouveau graphe réduit
func (g *Graph) MergeCommunities() {
	newGraph := NewGraph()
	newCommunities := make(map[int]int)
	communityMap := make(map[int]int)

	// Fusionner les arêtes entre les communautés
	for node, community := range g.Communities {
		for _, neighbor := range g.AdjList[node] {
			// Ne pas fusionner les arêtes dans la même communauté
			if g.Communities[neighbor] != community {
				newGraph.AddEdge(community, g.Communities[neighbor])
			}
		}
		if _, exists := communityMap[community]; !exists {
			communityMap[community] = len(communityMap) + 1
		}
		newCommunities[node] = communityMap[community]
	}

	// Remplacer l'ancien graphe par le nouveau graphe réduit
	*g = *newGraph
	g.Communities = newCommunities
}

// AddEdge ajoute une arête entre les nœuds u et v
func (g *Graph) AddEdge(u, v int) {
	if !contains(g.AdjList[u], v) {
		g.AdjList[u] = append(g.AdjList[u], v)
		g.AdjList[v] = append(g.AdjList[v], u)
	}
}

// contains vérifie si une valeur existe dans une slice
func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Louvain exécute l'algorithme Louvain avec un nombre spécifié de goroutines (workers)
func (g *Graph) Louvain(maxIterations int, numGoroutines int) {
	nodes := make([]int, 0, len(g.AdjList))
	for node := range g.AdjList {
		nodes = append(nodes, node)
		g.Communities[node] = node
	}
	sort.Ints(nodes)

	for iter := 0; iter < maxIterations; iter++ {
		improvement := false
		var wg sync.WaitGroup
		mu := sync.Mutex{} // Protège g.Communities et improvement

		// Canal pour distribuer les nœuds aux goroutines
		nodeChan := make(chan int, len(nodes))
		for _, node := range nodes {
			nodeChan <- node
		}
		close(nodeChan)

		// Lancer les workers
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				for node := range nodeChan {
					// Verrouiller l'accès à g.Communities
					mu.Lock()
					currentCommunity := g.Communities[node]
					mu.Unlock()

					bestCommunity := currentCommunity
					bestModularity := g.localModularity(node, currentCommunity)

					// Essayer de trouver une meilleure communauté en fonction des voisins
					for _, neighbor := range g.AdjList[node] {
						mu.Lock()
						neighborCommunity := g.Communities[neighbor]
						g.Communities[node] = neighborCommunity
						mu.Unlock()

						newModularity := g.localModularity(node, neighborCommunity)
						if newModularity > bestModularity {
							bestModularity = newModularity
							bestCommunity = neighborCommunity
							mu.Lock()
							improvement = true
							mu.Unlock()
						}
					}

					// Mettre à jour la communauté du nœud
					mu.Lock()
					g.Communities[node] = bestCommunity
					mu.Unlock()
				}
			}(i)
		}

		// Attendre la fin de toutes les goroutines
		wg.Wait()

		// Si aucune amélioration n'a été faite, on arrête
		if !improvement {
			break
		}

		// Fusionner les communautés après l'amélioration
		g.MergeCommunities()
	}
}

// localModularity calcule la modularité locale pour un nœud et une communauté donnée
func (g *Graph) localModularity(node, community int) float64 {
	m := float64(0)
	// Calcul du nombre total d'arêtes
	for _, neighbors := range g.AdjList {
		m += float64(len(neighbors))
	}
	m /= 2

	var Q float64
	// Calcul de la modularité locale
	ki := float64(len(g.AdjList[node]))
	for _, neighbor := range g.AdjList[node] {
		kj := float64(len(g.AdjList[neighbor]))
		if g.Communities[neighbor] == community {
			Q += 1.0 - (ki*kj)/(2.0*m)
		}
	}
	return Q / (2 * m)
}

func main() {
	// Création du graphe
	g := NewGraph()
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 1)
	g.AddEdge(4, 1)

	// Exécution de l'algorithme Louvain avec 2 workers
	g.Louvain(10, 4)
	fmt.Println("Communautés après Louvain:", g.Communities)
}