package main

import (
	"bufio"
	"os"
	"strings"
)

// chargerEnv lit le fichier .env (indispensable pour cacher le Client Secret Spotify)
func chargerEnv() {
	// On essaie d'ouvrir le fichier .env
	fichier, erreur := os.Open(".env")
	if erreur != nil {
		// S'il n'y a pas de fichier, on quitte la fonction sans faire planter le serveur
		return
	}

	// defer : permet de s'assurer que le fichier se fermera tout seul à la fin de la fonction
	defer fichier.Close()

	// On prépare un "lecteur" pour parcourir le fichier ligne par ligne
	lecteur := bufio.NewScanner(fichier)

	for lecteur.Scan() {
		// On récupère la ligne en enlevant les espaces en trop autour
		ligne := strings.TrimSpace(lecteur.Text())

		// On skip les lignes vides ou les commentaires (qui commencent par #)
		if ligne == "" || strings.HasPrefix(ligne, "#") {
			continue
		}

		// On coupe la ligne au niveau du signe "=" (ex: SPOTIFY_SECRET=12345)
		nom, valeur, trouve := strings.Cut(ligne, "=")

		// S'il n'y a pas de "=", ce n'est pas une variable valide, on passe à la suite
		if !trouve {
			continue
		}

		// On injecte la clé et la valeur dans l'environnement du serveur
		os.Setenv(strings.TrimSpace(nom), strings.TrimSpace(valeur))
	}
}
