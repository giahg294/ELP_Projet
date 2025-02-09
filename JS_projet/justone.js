const fs = require("fs"); // Importer le module du système de fichiers
const readline = require("readline"); // Module pour des interactions dans l'invite de commande

// Interface pour entrer du texte dans le terminal
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

// Liste des mots mystères
const cartes = [
    { carte: ["chien", "chat", "lapin"] },
    { carte: ["arbre", "plante", "fleur"] },
    { carte: ["voiture", "vélo", "moto"] },
    { carte: ["pomme", "banane", "orange"] },
    { carte: ["mer", "plage", "océan"] },
    { carte: ["maison", "appartement", "villa"] },
    { carte: ["montagne", "colline", "vallée"] },
    { carte: ["ordinateur", "téléphone", "tablette"] },
    { carte: ["football", "basketball", "tennis"] },
    { carte: ["école", "université", "lycée"] }
];
const nround = 5;

let motMystere = "";
let round = 0; // Suivi du round actuel
let indices = []; // Liste des indices pour un round
let joueurs = 5; // Nombre de joueurs
let indexJoueur = 0; // Suivi du joueur actuel
let score = 0; // Compteur de points
let logPropositions = []; // Stocke toutes les propositions pour l'écriture dans le fichier
let cartes_piochees = [];
let pioche = [];

// Piocher 5 cartes aléatoires
for (let i = 0; i < nround; i++) {
    pioche.push(Math.floor(Math.random() * 10));
}

for (let i = 0; i < nround; i++) {
    let num = pioche[i];
    cartes_piochees.push(cartes[num]);
}

console.log("🎲 Bienvenue dans le jeu JUST ONE !");
console.log(`Les ${nround} cartes piochées sont :\n\n${cartes_piochees.map(card => card.carte).join("\n")}\n`);
console.log(`📝 Il y aura ${nround} rounds.`);

// Fonction principale pour démarrer un round
async function startgame() {
    if (round < nround) {
        motMystere = await demanderChoix(); // Attendre la sélection du mot mystère
        while (motMystere == undefined) {
            console.log(`\nJoueur actif doit choisir entre 1 ou 2 ou 3 uniquement.`);
            motMystere = await demanderChoix();};
        indices = []; // Réinitialiser les indices
        indexJoueur = 0; // Réinitialiser le compteur de joueur
        console.log(`\nRound ${round + 1}: Le mot mystère est "${motMystere}" (gardez-le secret !)\n`);
        // Demander les indices à chaque joueur
        demanderIndices();
    } else {
        afficherScoreFinal();
        rl.close()
    }
}

// Fonction pour demander un choix au joueur pour le mot mystère
function demanderChoix() {
    return new Promise((resolve) => {
        rl.question(`Joueur actif, choisissez un chiffre entre 1 et 3 : `, (choix) => {
            let mot = cartes_piochees[round].carte[parseInt(choix) - 1];
            resolve(mot)});
        });
};


// Fonction pour demander un indice à chaque joueur
function demanderIndices() {
    if (indexJoueur < joueurs) {
        rl.question(`Joueur ${indexJoueur + 1}, entrez votre indice : `, (mot) => {
            // Vérifie si l'indice est le même que le mot mystère
            if (mot.toLowerCase().trim() === motMystere.toLowerCase()) {
                console.log("L'indice doit être différent du mot mystère !");
                // Si l'indice est invalide, redemander un indice
                demanderIndices();
            } else {
                // Si l'indice est valide, l'ajouter à la liste et passer au joueur suivant
                indices.push(mot.toLowerCase().trim());
                indexJoueur++;
                demanderIndices(); // Appel récursif pour demander l'indice du joueur suivant
            }
        });
    } else {
        // Si tous les joueurs ont donné un indice, on passe à la vérification des indices
        verifierIndices();
    }
}


// Fonction pour vérifier les indices et afficher les valides avant la devinette
function verifierIndices() {
    let occurrences = {};
    indices.forEach(mot => {
        occurrences[mot] = (occurrences[mot] || 0) + 1;
    });

    let indicesValides = indices.filter(mot => occurrences[mot] === 1);

    // Enregistrer les propositions dans un fichier
    logPropositions.push(`Round ${round + 1} : ${indices.join(", ")}`);
    fs.writeFileSync("propositions.txt", logPropositions.join("\n"), "utf-8");

    // Afficher les indices valides avant la devinette
    let msg = indicesValides.length > 0 ? indicesValides.join(", ") : "Aucun indice valide... 😢\n";
    console.log("\nIndices valides :");
    console.log(msg);

    // Demander à l'utilisateur de deviner le mot mystère
    rl.question(`Devinez le mot mystère (Round ${round + 1}) : `, (reponse) => {
        // Vérifier si la réponse est correcte
        if (reponse.toLowerCase().trim() === motMystere.toLowerCase()) {
            score += 1;
            console.log("\n🎉 Round validé ! (+1 point)\n");
        } else if (reponse.toLowerCase().trim() === "pass") {
            console.log("\n-- Round passé. (+0 point)\n");
        } else {
            score -= 1;
            console.log("\n🚫 Round non validé (mot mystère incorrect). (-1 point)\n");
        }

        // Passer au round suivant
        round++;
        setTimeout(startgame, 1000); // Attente de 2 secondes avant de démarrer le prochain round
    });
}

// Afficher le score final
function afficherScoreFinal() {
    console.log(`\n🏁 Jeu terminé !`);
    console.log(`Votre score final est : ${score} point(s)`);
};

// Lancer le premier round
startgame();
