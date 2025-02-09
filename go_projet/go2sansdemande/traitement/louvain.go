package traitement

import (
	"sort"
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

// Louvain exécute l'algorithme Louvain
func (g *Graph) Louvain(maxIterations int) {
	nodes := make([]int, 0, len(g.AdjList))
	for node := range g.AdjList {
		nodes = append(nodes, node)
		g.Communities[node] = node
	}
	sort.Ints(nodes)

	for iter := 0; iter < maxIterations; iter++ {
		improvement := false

		// Optimisation locale de la modularité
		for _, node := range nodes {
			currentCommunity := g.Communities[node]
			bestCommunity := currentCommunity
			bestModularity := g.Modularity()

			for _, neighbor := range g.AdjList[node] {
				g.Communities[node] = g.Communities[neighbor]
				newModularity := g.Modularity()
				if newModularity > bestModularity {
					bestModularity = newModularity
					bestCommunity = g.Communities[neighbor]
					improvement = true
				}
			}

			g.Communities[node] = bestCommunity
		}

		// Si aucune amélioration n'a été faite, terminer l'algorithme
		if !improvement {
			break
		}

		// Fusionner les communautés pour créer un nouveau graphe réduit
		g.MergeCommunities()
	}
}
