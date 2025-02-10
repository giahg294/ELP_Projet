import random

# Fonction pour générer un graphe
def generate_graph(filename, num_nodes, density):
    # Calcul du nombre maximal d'arêtes dans un graphe complet
    max_edges = num_nodes * (num_nodes - 1) // 2

    # Calcul du nombre d'arêtes en fonction de la densité
    edges_count = int(max_edges * density)

    # Initialisation de l'ensemble pour éviter les doublons
    edges = set()

    # Génération des arêtes
    while len(edges) < edges_count:
        # Choisir deux noeuds aléatoires différents
        node1, node2 = random.sample(range(1, num_nodes + 1), 2)
        # Ajouter l'arête (dans un ordre spécifique pour éviter les doublons)
        edge = tuple(sorted((node1, node2)))
        edges.add(edge)

    # Écriture des arêtes dans un fichier
    with open(filename, 'w') as f:
        for edge in edges:
            f.write(f"{edge[0]} {edge[1]}\n")
    print(f"Graph written to {filename}")

# Paramètres du graphe
num_nodes = 20
density = 0.2
filename = "data/minigraph3.txt"

# Générer le graphe et écrire dans le fichier
generate_graph(filename, num_nodes, density)
