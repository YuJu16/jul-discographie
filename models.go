package main

import (
	"fmt"
	"strings"
)

// Titre représente un morceau affiché sur le site
type Titre struct {
	ID         string
	Nom        string
	Artistes   []string
	DateSortie string // format "2024-05-10"
	DureeMs    int    // durée en millisecondes, comme Spotify
	Type       string // single, album...
	Pochette   string // lien de l'image
}

// ArtistesTexte colle les artistes dans un seul texte : "Artiste A, Artiste B"
func (t Titre) ArtistesTexte() string {
	return strings.Join(t.Artistes, ", ")
}

// Annee garde seulement l'année de "2024-05-10"
func (t Titre) Annee() string {
	if len(t.DateSortie) < 4 {
		return ""
	}
	return t.DateSortie[:4]
}

// Duree convertit des millisecondes en "3:42"
func (t Titre) Duree() string {
	secondes := t.DureeMs / 1000
	return fmt.Sprintf("%d:%02d", secondes/60, secondes%60)
}

// AFeat vaut vrai s'il y a plus d'un artiste sur le titre
func (t Titre) AFeat() bool {
	return len(t.Artistes) > 1
}
