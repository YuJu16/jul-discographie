package main

import (
	"fmt"
	"strings"
)

// Titre représente un morceau affiché sur le site
type Titre struct {
	ID          string
	Nom         string
	Artistes    []string
	DateSortie  string // format "2024-05-10"
	DureeMs     int    // durée en millisecondes, comme Spotify
	Type        string // single, album...
	Pochette    string // lien de l'image
	Album       string // nom de l'album
	LienSpotify string // lien pour ouvrir le titre sur Spotify
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

// Duree convertit des millisecondes en "03:12"
func (t Titre) Duree() string {
	secondes := t.DureeMs / 1000
	return fmt.Sprintf("%02d:%02d", secondes/60, secondes%60)
}

// DateFR transforme "2026-03-06" en "06/03/2026"
func (t Titre) DateFR() string {
	if len(t.DateSortie) != 10 {
		return t.DateSortie
	}
	return t.DateSortie[8:10] + "/" + t.DateSortie[5:7] + "/" + t.DateSortie[0:4]
}

// AFeat vaut vrai s'il y a plus d'un artiste sur le titre
func (t Titre) AFeat() bool {
	return len(t.Artistes) > 1
}
