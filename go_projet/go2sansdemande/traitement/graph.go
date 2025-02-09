package traitement

import (
	"fmt"
	"sort"
	"strings"
)

type Graph struct {
	AdjList     map[int][]int
	Communities map[int]int
}

// NewGraph crée un graphe vide
func NewGraph() *Graph {
	return &Graph{
		AdjList:     make(map[int][]int),
		Communities: make(map[int]int),
	}
}

// AddEdge ajoute une arête au graphe sans doublon
func (g *Graph) AddEdge(u, v int) {
	// Vérifier si l'arête existe déjà
	if contains(g.AdjList[u], v) || u == v {
		return
	}
	g.AdjList[u] = append(g.AdjList[u], v)
	g.AdjList[v] = append(g.AdjList[v], u)
}

// contains vérifie si une slice contient une valeur
func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// DisplayCommunities affiche les communautés sous forme de groupes de nœuds
// func (g *Graph) DisplayCommunities() {
// 	communityGroups := make(map[int][]int)
// 	for node, community := range g.Communities {
// 		communityGroups[community] = append(communityGroups[community], node)
// 	}

//		// Tri des communautés et des nœuds
//		for community, nodes := range communityGroups {
//			sort.Ints(nodes)
//			fmt.Printf("Community %d: %v\n", community, nodes)
//		}
//	}
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

// PrintGraph affiche le graphe sous forme de listes d'adjacence
func (g *Graph) PrintGraph() {
	fmt.Println("Graphe (listes d'adjacence) :")
	for node, neighbors := range g.AdjList {
		fmt.Printf("%d: %v\n", node, neighbors)
	}
}
