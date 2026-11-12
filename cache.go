package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Fichier de cache local pour éviter de re-taper Spotify à chaque lancement (anti-429)
const fichierCache = "cache_titres.json"

// Garde les sons en mémoire tant que le serveur tourne
var titresJul []Titre

// Tente de récupérer les sons depuis le fichier JSON local
func lireCache() ([]Titre, bool) {
	// os.ReadFile ouvre, lit et ferme le fichier d'un coup
	contenu, erreur := os.ReadFile(fichierCache)
	if erreur != nil {
		return nil, false // pas de fichier = pas de cache
	}

	var titres []Titre
	// On transforme le texte du fichier en tableau Go
	erreur = json.Unmarshal(contenu, &titres)

	// Fichier mal formé ou vide : on fait comme s'il n'y avait pas de cache
	if erreur != nil || len(titres) == 0 {
		return nil, false
	}

	return titres, true
}

// Sauvegarde le tableau de sons dans le fichier
func ecrireCache(titres []Titre) error {
	// MarshalIndent = comme JSON.stringify en JS, mais bien indenté pour être lisible
	contenu, erreur := json.MarshalIndent(titres, "", "  ")
	if erreur != nil {
		return erreur
	}

	// 0644 : permissions classiques (écriture pour moi, lecture pour les autres)
	return os.WriteFile(fichierCache, contenu, 0644)
}

// Remplit titresJul : d'abord depuis le cache, sinon depuis Spotify
func chargerTitres() error {

	// tokenActuel = "faux-token" // TEST à supprimer

	// 1. On tente le raccourci : le fichier local
	titres, trouve := lireCache()
	if trouve {
		titresJul = titres
		fmt.Println("📦 Sons chargés depuis le cache local :", len(titres))
		return nil // on s'arrête là, pas d'appel API pour rien
	}

	// 2. Pas de cache : on interroge Spotify (le token est géré dans appelSpotify)
	fmt.Println("🌐 Aucun cache trouvé, on va chercher sur Spotify...")

	albums, erreur := recupererAlbums()
	if erreur != nil {
		return erreur
	}

	titres, erreur = recupererTitres(albums)
	if erreur != nil {
		return erreur
	}

	// 3. On garde les titres en mémoire...
	titresJul = titres
	fmt.Println("✅ Sons récupérés via l'API :", len(titres))

	// ... et dans le fichier pour la prochaine fois
	return ecrireCache(titres)
}
