package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// fauxTitres renvoie des titres d'exemple, en attendant Spotify
func fauxTitres() []Titre {
	return []Titre{
		{ID: "1", Nom: "Titre 1", Artistes: []string{"Artiste A"}, DateSortie: "2024-05-10", DureeMs: 222000, Type: "single", Album: "Album 1", LienSpotify: "https://open.spotify.com"},
		{ID: "2", Nom: "Titre 2", Artistes: []string{"Artiste A", "Artiste B"}, DateSortie: "2023-11-03", DureeMs: 198000, Type: "album", Album: "Album 2", LienSpotify: "https://open.spotify.com"},
		{ID: "3", Nom: "Titre 3", Artistes: []string{"Artiste A"}, DateSortie: "2022-02-18", DureeMs: 245000, Type: "album", Album: "Album 2", LienSpotify: "https://open.spotify.com"},
	}
}

// PageErreur contient ce qu'on affiche sur la page d'erreur
type PageErreur struct {
	Code    int
	Titre   string
	Message string
}

// afficherErreur montre la belle page d'erreur avec le bon code (404, 503...)
func afficherErreur(reponse http.ResponseWriter, code int, titre string, message string) {
	page, erreur := template.ParseFiles("templates/erreur.gohtml")
	if erreur != nil {
		// si même la page d'erreur est cassée, on envoie un texte simple
		http.Error(reponse, titre, code)
		return
	}

	reponse.WriteHeader(code) // le code doit être envoyé avant la page
	page.Execute(reponse, PageErreur{Code: code, Titre: titre, Message: message})
}

// Handler de l'accueil : récupère la requête (demande) et écrit la vue dans (reponse)
func accueil(reponse http.ResponseWriter, demande *http.Request) {
	// "/" attrape toutes les URLs : si ce n'est pas exactement "/", la page n'existe pas
	if demande.URL.Path != "/" {
		afficherErreur(reponse, http.StatusNotFound, "Page introuvable", "Cette page n'existe pas.")
		return
	}

	// Chargement du fichier HTML
	page, erreur := template.ParseFiles("templates/index.gohtml")
	if erreur != nil { // si erreur n'est pas vide (nil), fichier introuvable -> erreur 500 et on stoppe
		http.Error(reponse, "Impossible de charger la page", http.StatusInternalServerError)
		return
	}

	// La liste des titres (faux pour l'instant)
	titres := fauxTitres()

	// Envoi du HTML au navigateur, avec les titres injectés dans le template
	erreur = page.Execute(reponse, titres)
	if erreur != nil {
		http.Error(reponse, "Impossible d'afficher la page", http.StatusInternalServerError)
		return
	}
}

// detail affiche la page d'un seul titre, retrouvé grâce à l'id dans l'URL
func detail(reponse http.ResponseWriter, demande *http.Request) {
	id := demande.PathValue("id") // la partie {id} de /tracks/{id}

	page, erreur := template.ParseFiles("templates/detail.gohtml")
	if erreur != nil {
		http.Error(reponse, "Impossible de charger la page", http.StatusInternalServerError)
		return
	}

	// on cherche le titre qui a cet id
	for _, titre := range fauxTitres() {
		if titre.ID == id {
			erreur = page.Execute(reponse, titre)
			if erreur != nil {
				http.Error(reponse, "Impossible d'afficher la page", http.StatusInternalServerError)
			}
			return
		}
	}

	// aucun titre trouvé : page 404
	afficherErreur(reponse, http.StatusNotFound, "Titre introuvable", "Ce titre n'existe pas ou n'est plus disponible.")
}

func main() {
	// Route principale "/" reliée à la fonction accueil
	http.HandleFunc("/", accueil)
	http.HandleFunc("/tracks/{id}", detail)

	// Les fichiers du dossier styles sont servis sur /styles/
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))

	fmt.Println("Serveur lancé sur http://localhost:8080")

	// Lancement du serveur sur le port 8080 (nil = routeur par défaut de Go)
	erreur := http.ListenAndServe(":8080", nil)
	if erreur != nil {
		fmt.Println("Erreur serveur :", erreur)
	}
}
