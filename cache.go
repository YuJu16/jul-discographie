package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Notre fichier de cache local pour éviter de cramer notre quota Spotify (anti-erreur 429)
const fichierCache = "cache_titres.json"

// Variable globale qui garde les sons en mémoire tant que le serveur est allumé
var titresJul []Titre

// Tente de récupérer les sons depuis notre fichier JSON local
func lireCache() ([]Titre, bool) {
	// os.ReadFile fait tout d'un coup (ouvre, lit, et ferme le fichier)
	contenu, erreur := os.ReadFile(fichierCache)
	if erreur != nil {
		return nil, false // Fichier introuvable, on renvoie false
	}

	var titres []Titre
	// On transforme le texte du fichier en vrai tableau Go
	erreur = json.Unmarshal(contenu, &titres)

	// Si le fichier est mal formaté ou vide, on fait comme s'il n'y avait pas de cache
	if erreur != nil || len(titres) == 0 {
		return nil, false
	}

	return titres, true
}

// Sauvegarde notre tableau de sons dans le fichier
func ecrireCache(titres []Titre) error {
	// MarshalIndent = comme JSON.stringify en JS, mais formaté proprement avec des espaces pour être lisible
	contenu, erreur := json.MarshalIndent(titres, "", "  ")
	if erreur != nil {
		return erreur
	}

	// 0644 : permissions standards de fichier (lecture/écriture pour toi, lecture pour le reste)
	return os.WriteFile(fichierCache, contenu, 0644)
}

// Le chef d'orchestre : remplit titresJul depuis le cache, ou via l'API si le cache est vide
func chargerTitres() error {
	// 1. On check d'abord notre fichier local
	if titres, trouve := lireCache(); trouve {
		titresJul = titres
		fmt.Println("📦 Sons chargés depuis le cache local :", len(titres))
		return nil // On arrête la fonction ici, pas besoin d'appeler Spotify !
	}

	// 2. Si on arrive là, c'est qu'on a pas de cache, donc on appelle l'API
	fmt.Println("🌐 Aucun cache trouvé, on va taper l'API Spotify...")

	token, erreur := demanderToken()
	if erreur != nil {
		return erreur
	}

	albums, erreur := recupererAlbums(token)
	if erreur != nil {
		return erreur
	}

	titres, erreur := recupererTitres(token, albums)
	if erreur != nil {
		return erreur
	}

	// 3. On met à jour notre variable en mémoire...
	titresJul = titres
	fmt.Println("✅ Sons récupérés via Spotify :", len(titres))

	// ...et on crée le fichier cache pour la prochaine fois !
	return ecrireCache(titres)
}
