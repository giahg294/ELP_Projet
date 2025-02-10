package traitement

import (
	"sort"
	"sync"
)

// Modularity calcule la modularité
func (g *Graph) Modularity() float64 {
	m := float64(0)
	for _, neighbors := range g.AdjList {
		m += float64(len(neighbors))
	}
	m /= 2

	var Q float64
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

// Louvain exécute l'algorithme Louvain avec un nombre spécifié de Goroutines
func (g *Graph) Louvain(maxIterations int, numGoroutines int) {
	nodes := make([]int, 0, len(g.AdjList))
	for node := range g.AdjList {
		nodes = append(nodes, node)
		g.Communities[node] = node
	}
	sort.Ints(nodes)

	for iter := 0; iter < maxIterations; iter++ {
		improvement := false

		// Allouer les nœuds à différents Goroutines
		nodeChunks := chunk(nodes, numGoroutines)
		var wg sync.WaitGroup
		mu := sync.Mutex{} // Le mutex protège la structure de données partagée.
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(chunk []int) {
				defer wg.Done()
				for _, node := range chunk {
					currentCommunity := g.Communities[node]
					bestCommunity := currentCommunity
					bestModularity := g.Modularity()

					for _, neighbor := range g.AdjList[node] {
						mu.Lock()
						g.Communities[node] = g.Communities[neighbor]
						mu.Unlock()
						newModularity := g.Modularity()
						if newModularity > bestModularity {
							bestModularity = newModularity
							bestCommunity = g.Communities[neighbor]
							improvement = true
						}
					}

					mu.Lock()
					g.Communities[node] = bestCommunity
					mu.Unlock()
				}
			}(nodeChunks[i])
		}
		wg.Wait()

		// Si aucune amélioration n'est apportée, alors quitter.
		if !improvement {
			break
		}

		// Fusionner les communautés.
		g.MergeCommunities()
	}
}

// Chunk divise les nœuds en plusieurs blocs pour un traitement parallèle.
func chunk(nodes []int, numChunks int) [][]int {
	chunks := make([][]int, numChunks)
	for i := range chunks {
		chunks[i] = make([]int, 0)
	}
	for i, node := range nodes {
		chunks[i%numChunks] = append(chunks[i%numChunks], node)
	}
	return chunks
}
