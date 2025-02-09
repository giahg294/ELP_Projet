# JS_projet

## Présentation générale du jeu

Le jeu **Just One** se base sur 10 cartes chaque contient 3 propositions de mot mystère à faire deviner. Des indices sont fournis par les joueurs pour deviner ce mot. À la fin de chaque round, le joueur doit deviner le mot mystère en fonction des indices donnés par les autres joueurs. Le score est mis à jour en fonction de la validité des indices (doublons ou non) et si la devinette est correcte.

## Structure du code

Le code est principalement composé de 3 grandes sections :
1. **Initialisation du jeu** (listes, variables)
2. **Interaction avec les joueurs** (saisie des indices, devinette)
3. **Gestion de la logique du jeu** (validation des indices, calcul du score)

## Résumé global du flux du jeu
1. **Initialisation** : Le jeu commence par une pioche de 5 cartes.
2. **Demande du choix de mot** : Le joueur actif choisit un numéro entre 1 et 3.
2. **Demande des indices** : Les joueurs donnent chacun un indice.
3. **Vérification des indices** : Les indices sont analysés pour détecter les doublons et ceux valides sont affichés.
4. **Devinette du mot mystère** : Le joueur tente de deviner le mot mystère en fonction des indices.
5. **Calcul du score** : Le score est ajusté en fonction de la réussite ou de l'échec de la devinette.
6. **Fichier de sauvegarde** : Les propositions de chaque round sont sauvegardées dans un fichier pour garder une trace de la partie.

## Lancement du jeu
Executer la commande npm start dans le terminal.

## Conclusion
Ce code représente une version simple mais complète du jeu **Just One**, où la logique principale est de donner des indices, deviner le mot mystère, et mettre à jour un score. La fonctionnalité de sauvegarde des propositions dans un fichier permet de conserver un historique du jeu. Ce jeu peut être amélioré avec des fonctionnalités supplémentaires comme un système de chronomètre ou des niveaux de difficulté.
